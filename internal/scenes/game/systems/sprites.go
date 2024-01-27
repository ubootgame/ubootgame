package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"gonum.org/v1/gonum/spatial/r2"
)

func DrawSprites(e *ecs.ECS, screen *ebiten.Image) {
	ecs.NewQuery(layers.Foreground, filter.Contains(components.Sprite, components.Position)).Each(e.World, func(entry *donburi.Entry) {
		spriteData := components.Sprite.Get(entry)
		positionData := components.Position.Get(entry)

		cameraEntry, _ := components.Camera.First(e.World)
		camera := components.Camera.Get(cameraEntry)

		w, h := float64(spriteData.Image.Bounds().Dx()), float64(spriteData.Image.Bounds().Dy())

		op := &ebiten.DrawImageOptions{}

		op.GeoM = utility.PositionMatrix(positionData.Center, r2.Vec{X: w, Y: h}, positionData.Scale, positionData.ScaleDirection)
		op.GeoM.Concat(utility.CameraMatrix(camera))

		op.Filter = ebiten.FilterLinear

		screen.DrawImage(spriteData.Image, op)
	})
}
