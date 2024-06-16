package weapons

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/framework"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/internal/components"
	"github.com/ubootgame/ubootgame/internal/components/physics"
	"github.com/ubootgame/ubootgame/internal/entities/weapons"
	"github.com/ubootgame/ubootgame/internal/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/colornames"
	"gonum.org/v1/gonum/spatial/r2"
)

type BulletSystem struct {
	display framework.DisplayService

	camera *components.CameraData

	query *donburi.Query
	tick  uint64
	image *ebiten.Image

	drawImageOptions *ebiten.DrawImageOptions
}

func NewBulletSystem(display framework.DisplayService) *BulletSystem {
	return &BulletSystem{
		display:          display,
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
	if entry, found := components.Camera.First(e.World); found {
		system.camera = components.Camera.Get(entry)
	}

	weapons.BulletTag.Each(e.World, func(bulletEntry *donburi.Entry) {
		//worldPosition := transform.WorldPosition(bulletEntry)
		//
		//actors.EnemyTag.Each(e.World, func(enemyEntry *donburi.Entry) {
		//	bounds := static.Bounds.Get(enemyEntry)
		//
		//	if bounds.PointInside(resolv.Vector{X: worldPosition.X, Y: worldPosition.Y}) {
		//		e.World.Remove(enemyEntry.Entity())
		//		e.World.Remove(bulletEntry.Entity())
		//	}
		//})
	})
}

func (system *BulletSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	if system.image == nil {
		system.image = ebiten.NewImage(2, 2)
		system.image.Fill(colornames.White)
	}

	system.query.Each(e.World, func(entry *donburi.Entry) {
		body := physics.Body.Get(entry)

		position := body.Position()

		system.drawImageOptions.GeoM.Reset()
		system.drawImageOptions.GeoM.Translate(-1, -1)
		system.drawImageOptions.GeoM.Translate(system.display.WorldToScreen(r2.Vec(position)))

		system.camera.Camera.Draw(system.image, system.drawImageOptions, screen)
	})
}
