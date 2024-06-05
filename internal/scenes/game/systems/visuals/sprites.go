package visuals

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/injector"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/systems"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type SpriteSystem struct {
	systems.BaseSystem

	camera  *game_system.CameraData
	debug   *game_system.DebugData
	display *game_system.DisplayData

	query     *donburi.Query
	debugText string
}

func NewSpriteSystem() *SpriteSystem {
	system := &SpriteSystem{
		query: ecs.NewQuery(layers.Foreground, filter.Contains(visuals.Sprite, geometry.Transform)),
	}
	system.Injector = injector.NewInjector([]injector.Injection{
		injector.Once([]injector.Injection{
			injector.Component(&system.camera, game_system.Camera),
			injector.Component(&system.debug, game_system.Debug),
			injector.Component(&system.display, game_system.Display),
		}),
	})
	return system
}

func (system *SpriteSystem) Update(e *ecs.ECS) {
	system.BaseSystem.Update(e)

	if system.debug.Enabled && system.debug.DrawPositions {
		system.query.Each(e.World, func(entry *donburi.Entry) {
			sprite := visuals.Sprite.Get(entry)
			transform := geometry.Transform.Get(entry)

			debugText := fmt.Sprintf("Transform: %.3f, %.3f\nSize: %.3f, %.3f",
				transform.Center.X, transform.Center.Y,
				transform.Size.X, transform.Size.Y)

			if entry.HasComponent(geometry.Velocity) {
				velocity := geometry.Velocity.Get(entry)
				debugText += fmt.Sprintf("\nVelocity: %.3f, %.3f", velocity.X, velocity.Y)
			}

			sprite.DebugText = debugText
		})
	}
}

func (system *SpriteSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	system.query.Each(e.World, func(entry *donburi.Entry) {
		sprite := visuals.Sprite.Get(entry)
		transform := geometry.Transform.Get(entry)

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
		op.GeoM.Concat(*system.camera.Matrix)

		op.Filter = ebiten.FilterLinear

		screen.DrawImage(sprite.Image, op)

		if system.debug.Enabled && system.debug.DrawPositions {
			system.drawDebug(transform, screen, sprite)
		}
	})
}

func (system *SpriteSystem) drawDebug(transform *geometry.TransformData, screen *ebiten.Image, sprite *visuals.SpriteData) {
	debugOpts := &ebiten.DrawImageOptions{}
	debugOpts.GeoM.Scale(1/system.display.VirtualResolution.X, 1/system.display.VirtualResolution.X)
	debugOpts.GeoM.Scale(1.0/system.camera.ZoomFactor, 1.0/system.camera.ZoomFactor)
	debugOpts.GeoM.Translate(transform.Center.X+transform.Size.X/2, transform.Center.Y+transform.Size.Y/2)
	debugOpts.GeoM.Concat(*system.camera.Matrix)

	utility.Debug.PrintDebugTextAt(screen, sprite.DebugText, debugOpts)
}
