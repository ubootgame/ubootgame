package systems

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type spriteSystem struct {
	cameraEntry, debugEntry, displayEntry *donburi.Entry
	query                                 *donburi.Query
	debugText                             string
}

var Sprite = &spriteSystem{
	query: ecs.NewQuery(layers.Foreground, filter.Contains(components.Sprite, components.Transform)),
}

func (system *spriteSystem) Update(e *ecs.ECS) {
	var ok bool
	if system.cameraEntry == nil {
		if system.cameraEntry, ok = components.Camera.First(e.World); !ok {
			panic("no camera found")
		}
	}
	if system.debugEntry == nil {
		if system.debugEntry, ok = components.Debug.First(e.World); !ok {
			panic("no debug found")
		}
	}
	if system.displayEntry == nil {
		if system.displayEntry, ok = components.Display.First(e.World); !ok {
			panic("no display found")
		}
	}

	debug := components.Debug.Get(system.debugEntry)

	if debug.Enabled && debug.DrawPositions {
		system.query.Each(e.World, func(entry *donburi.Entry) {
			sprite := components.Sprite.Get(entry)
			transform := components.Transform.Get(entry)

			debugText := fmt.Sprintf("Transform: %.3f, %.3f\nSize: %.3f, %.3f",
				transform.Center.X, transform.Center.Y,
				transform.Size.X, transform.Size.Y)

			if entry.HasComponent(components.Velocity) {
				velocity := components.Velocity.Get(entry)
				debugText += fmt.Sprintf("\nVelocity: %.3f, %.3f", velocity.X, velocity.Y)
			}

			sprite.DebugText = debugText
		})
	}
}

func (system *spriteSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	debug := components.Debug.Get(system.debugEntry)
	camera := components.Camera.Get(system.cameraEntry)
	display := components.Display.Get(system.displayEntry)

	system.query.Each(e.World, func(entry *donburi.Entry) {
		sprite := components.Sprite.Get(entry)
		transform := components.Transform.Get(entry)

		op := &ebiten.DrawImageOptions{}

		if transform.FlipX {
			op.GeoM.Scale(1, -11)
			op.GeoM.Translate(0, float64(sprite.Image.Bounds().Size().Y))
		}
		if transform.FlipY {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(float64(sprite.Image.Bounds().Size().X), 0)
		}

		op.GeoM.Translate(-float64(sprite.Image.Bounds().Size().X/2), -float64(sprite.Image.Bounds().Size().Y/2))
		op.GeoM.Scale(sprite.Scale*transform.Size.X, sprite.Scale*transform.Size.X)
		op.GeoM.Translate(transform.Center.X, transform.Center.Y)
		op.GeoM.Concat(*camera.Matrix)

		op.Filter = ebiten.FilterLinear

		screen.DrawImage(sprite.Image, op)

		if debug.Enabled && debug.DrawPositions {
			system.drawDebug(display, camera, transform, screen, sprite)
		}
	})
}

func (system *spriteSystem) drawDebug(display *components.DisplayData, camera *components.CameraData, transform *components.TransformData, screen *ebiten.Image, sprite *components.SpriteData) {
	debugOpts := &ebiten.DrawImageOptions{}
	debugOpts.GeoM.Scale(1/display.VirtualResolution.X, 1/display.VirtualResolution.X)
	debugOpts.GeoM.Scale(1.0/camera.ZoomFactor, 1.0/camera.ZoomFactor)
	debugOpts.GeoM.Translate(transform.Center.X+transform.Size.X/2, transform.Center.Y+transform.Size.Y/2)
	debugOpts.GeoM.Concat(*camera.Matrix)

	utility.Debug.PrintDebugTextAt(screen, sprite.DebugText, debugOpts)
}
