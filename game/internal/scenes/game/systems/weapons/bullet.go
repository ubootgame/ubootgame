package weapons

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"github.com/solarlune/resolv"
	ecsFramework "github.com/ubootgame/ubootgame/framework/game/ecs"
	"github.com/ubootgame/ubootgame/framework/graphics/camera"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/actors"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/weapons"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/colornames"
)

type BulletSystem struct {
	camera *camera.Camera

	query *donburi.Query
	tick  uint64
	image *ebiten.Image

	drawImageOptions *ebiten.DrawImageOptions
}

func NewBulletSystem(camera *camera.Camera) *BulletSystem {
	return &BulletSystem{
		camera:           camera,
		query:            donburi.NewQuery(filter.Contains(weapons.BulletTag)),
		drawImageOptions: &ebiten.DrawImageOptions{},
	}
}

func (system *BulletSystem) Layers() []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer]{
		{A: layers.Game, B: system.Draw},
	}
}

func (system *BulletSystem) Update(e *ecs.ECS) {
	weapons.BulletTag.Each(e.World, func(bulletEntry *donburi.Entry) {
		worldPosition := transform.WorldPosition(bulletEntry)

		actors.EnemyTag.Each(e.World, func(enemyEntry *donburi.Entry) {
			bounds := geometry.Bounds.Get(enemyEntry)

			if bounds.PointInside(resolv.Vector{X: worldPosition.X, Y: worldPosition.Y}) {
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
		worldPosition := transform.WorldPosition(entry)
		worldScale := transform.WorldScale(entry)

		system.drawImageOptions.GeoM.Reset()
		system.drawImageOptions.GeoM.Scale(worldScale.X, worldScale.Y)
		system.drawImageOptions.GeoM.Translate(worldPosition.X, worldPosition.Y)
		system.camera.Apply(&system.drawImageOptions.GeoM)

		screen.DrawImage(system.image, system.drawImageOptions)
	})
}
