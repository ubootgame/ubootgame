package systems

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type sprites struct {
	query     *donburi.Query
	debugText string
}

var Sprites = &sprites{
	query: ecs.NewQuery(layers.Foreground, filter.Contains(components.Sprite, components.Position)),
}

func (system *sprites) Update(e *ecs.ECS) {
	if config.C.Debug {
		system.query.Each(e.World, func(entry *donburi.Entry) {
			spriteData := components.Sprite.Get(entry)
			positionData := components.Position.Get(entry)

			debugText := fmt.Sprintf("Position: %.3f, %.3f\nSize: %.3f, %.3f",
				positionData.Center.X, positionData.Center.Y,
				positionData.Size.X, positionData.Size.Y)

			if entry.HasComponent(components.Velocity) {
				velocityData := components.Velocity.Get(entry)
				debugText += fmt.Sprintf("\nVelocity: %.3f, %.3f", velocityData.X, velocityData.Y)
			}

			spriteData.DebugText = debugText
		})
	}
}

func (system *sprites) Draw(e *ecs.ECS, screen *ebiten.Image) {
	system.query.Each(e.World, func(entry *donburi.Entry) {
		spriteData := components.Sprite.Get(entry)
		positionData := components.Position.Get(entry)

		cameraEntry, _ := components.Camera.First(e.World)
		camera := components.Camera.Get(cameraEntry)

		op := &ebiten.DrawImageOptions{}

		op.GeoM.Translate(-float64(spriteData.Image.Bounds().Size().X/2), -float64(spriteData.Image.Bounds().Size().Y/2))
		op.GeoM.Scale(spriteData.Scale*positionData.Size.X, spriteData.Scale*positionData.Size.X)
		op.GeoM.Translate(positionData.Center.X, positionData.Center.Y)
		op.GeoM.Concat(utility.CameraMatrix(camera))

		op.Filter = ebiten.FilterLinear

		screen.DrawImage(spriteData.Image, op)

		if config.C.Debug {
			debugOpts := &ebiten.DrawImageOptions{}
			debugOpts.GeoM.Scale(1/config.C.VirtualResolution.X, 1/config.C.VirtualResolution.X)
			debugOpts.GeoM.Translate(positionData.Center.X+positionData.Size.X/2, positionData.Center.Y+positionData.Size.Y/2)
			debugOpts.GeoM.Concat(utility.CameraMatrix(camera))

			Debug.printDebugTextAt(screen, spriteData.DebugText, debugOpts)
		}
	})
}
