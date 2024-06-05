package weapons

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/actors"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/weapons"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/injector"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/systems"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/colornames"
)

type BulletSystem struct {
	systems.BaseSystem

	camera *game_system.CameraData

	query *donburi.Query
	tick  uint64
	image *ebiten.Image
}

func NewBulletSystem() *BulletSystem {
	system := &BulletSystem{
		query: donburi.NewQuery(filter.Contains(weapons.BulletTag)),
	}
	system.Injector = injector.NewInjector([]injector.Injection{
		injector.Once([]injector.Injection{
			injector.Component(&system.camera, game_system.Camera),
		}),
	})
	return system
}

func (system *BulletSystem) Update(e *ecs.ECS) {
	system.BaseSystem.Update(e)

	weapons.BulletTag.Each(e.World, func(bulletEntry *donburi.Entry) {
		transform := geometry.Transform.Get(bulletEntry)

		bulletScreen := system.camera.WorldToScreenPosition(transform.Center)

		actors.EnemyTag.Each(e.World, func(enemyEntry *donburi.Entry) {
			shape := geometry.Shape.Get(enemyEntry)

			if shape.PointInside(resolv.Vector{X: bulletScreen.X, Y: bulletScreen.Y}) {
				e.World.Remove(enemyEntry.Entity())
				e.World.Remove(bulletEntry.Entity())
			}
		})
	})
}

func (system *BulletSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	if system.image == nil {
		system.image = ebiten.NewImage(2, 2)
		system.image.Fill(colornames.White)
	}

	system.query.Each(e.World, func(entry *donburi.Entry) {
		transform := geometry.Transform.Get(entry)

		op := &ebiten.DrawImageOptions{}

		op.GeoM.Translate(-1, -1)
		op.GeoM.Scale(0.001, 0.001)
		op.GeoM.Translate(transform.Center.X, transform.Center.Y)
		op.GeoM.Concat(*system.camera.Matrix)

		screen.DrawImage(system.image, op)
	})
}
