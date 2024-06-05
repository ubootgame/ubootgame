package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quartercastle/vector"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/colornames"
)

type bulletSystem struct {
	cameraEntry *donburi.Entry
	query       *donburi.Query
	tick        uint64
	image       *ebiten.Image
}

var Bullet = &bulletSystem{
	query: donburi.NewQuery(filter.Contains(entities.BulletTag)),
}

func (system *bulletSystem) Update(e *ecs.ECS) {
	if system.cameraEntry == nil {
		var ok bool
		if system.cameraEntry, ok = game_system.Camera.First(e.World); !ok {
			panic("no camera found")
		}
	}

	camera := game_system.Camera.Get(system.cameraEntry)

	entities.BulletTag.Each(e.World, func(bulletEntry *donburi.Entry) {
		transform := components.Transform.Get(bulletEntry)

		bulletScreen := camera.WorldToScreenPosition(transform.Center)

		entities.EnemyTag.Each(e.World, func(enemyEntry *donburi.Entry) {
			shape := components.Shape.Get(enemyEntry)

			if shape.PointInside(vector.Vector{bulletScreen.X, bulletScreen.Y}) {
				e.World.Remove(enemyEntry.Entity())
			}
		})
	})
}

func (system *bulletSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	if system.image == nil {
		system.image = ebiten.NewImage(2, 2)
		system.image.Fill(colornames.White)
	}

	camera := game_system.Camera.Get(system.cameraEntry)

	system.query.Each(e.World, func(entry *donburi.Entry) {
		transform := components.Transform.Get(entry)

		op := &ebiten.DrawImageOptions{}

		op.GeoM.Translate(-1, -1)
		op.GeoM.Scale(0.001, 0.001)
		op.GeoM.Translate(transform.Center.X, transform.Center.Y)
		op.GeoM.Concat(*camera.Matrix)

		screen.DrawImage(system.image, op)
	})
}
