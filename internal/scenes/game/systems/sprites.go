package systems

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"math"
)

type spriteSystem struct {
	cameraEntry *donburi.Entry
	query       *donburi.Query
	debugText   string
}

var Sprite = &spriteSystem{
	query: ecs.NewQuery(layers.Foreground, filter.Contains(components.Sprite, components.Position)),
}

func (system *spriteSystem) Update(e *ecs.ECS) {
	if system.cameraEntry == nil {
		var ok bool
		if system.cameraEntry, ok = components.Camera.First(e.World); !ok {
			panic("no camera found")
		}
	}

	if config.C.Debug {
		system.query.Each(e.World, func(entry *donburi.Entry) {
			sprite := components.Sprite.Get(entry)
			position := components.Position.Get(entry)

			debugText := fmt.Sprintf("Position: %.3f, %.3f\nSize: %.3f, %.3f",
				position.Center.X, position.Center.Y,
				position.Size.X, position.Size.Y)

			if entry.HasComponent(components.Velocity) {
				velocity := components.Velocity.Get(entry)
				debugText += fmt.Sprintf("\nVelocity: %.3f, %.3f", velocity.X, velocity.Y)
			}

			sprite.DebugText = debugText
		})
	}
}

func (system *spriteSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	system.query.Each(e.World, func(entry *donburi.Entry) {
		sprite := components.Sprite.Get(entry)
		position := components.Position.Get(entry)
		camera := components.Camera.Get(system.cameraEntry)

		op := &ebiten.DrawImageOptions{}

		op.GeoM.Translate(-float64(sprite.Image.Bounds().Size().X/2), -float64(sprite.Image.Bounds().Size().Y/2))
		op.GeoM.Scale(sprite.Scale*position.Size.X, sprite.Scale*position.Size.X)
		op.GeoM.Translate(position.Center.X, position.Center.Y)
		op.GeoM.Concat(*camera.Matrix)

		op.Filter = ebiten.FilterLinear

		screen.DrawImage(sprite.Image, op)

		if config.C.Debug {
			debugOpts := &ebiten.DrawImageOptions{}
			debugOpts.GeoM.Scale(1/config.C.VirtualResolution.X, 1/config.C.VirtualResolution.X)
			debugOpts.GeoM.Scale(1.0/camera.ZoomFactor, 1.0/camera.ZoomFactor)
			debugOpts.GeoM.Translate(position.Center.X+math.Abs(position.Size.X)/2, position.Center.Y+math.Abs(position.Size.Y)/2)
			debugOpts.GeoM.Concat(*camera.Matrix)

			Debug.printDebugTextAt(screen, sprite.DebugText, debugOpts)
		}
	})
}
