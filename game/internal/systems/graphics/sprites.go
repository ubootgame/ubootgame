package graphics

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/framework"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/framework/graphics/d2d"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/components"
	"github.com/ubootgame/ubootgame/internal/components/graphics"
	"github.com/ubootgame/ubootgame/internal/components/physics"
	"github.com/ubootgame/ubootgame/internal/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/colornames"
	"gonum.org/v1/gonum/spatial/r2"
	"math"
)

type SpriteSystem struct {
	settings framework.SettingsService[internal.Settings]
	display  framework.DisplayService

	camera *components.CameraData

	query *donburi.Query
}

func NewSpriteSystem(settings framework.SettingsService[internal.Settings], display framework.DisplayService) *SpriteSystem {
	return &SpriteSystem{
		settings: settings,
		display:  display,
		query: donburi.NewQuery(
			filter.And(
				filter.Contains(graphics.Sprite),
				filter.Or(
					filter.Contains(transform.Transform),
					filter.Contains(physics.Body)))),
	}
}

func (system *SpriteSystem) Layers() []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer]{
		{A: layers.Game, B: system.Draw},
		{A: layers.Debug, B: system.DrawDebug},
	}
}

func (system *SpriteSystem) Update(e *ecs.ECS) {
	if entry, found := components.Camera.First(e.World); found {
		system.camera = components.Camera.Get(entry)
	}
}

func (system *SpriteSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	spriteDrawImageOptions := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}

	var (
		position, screenSize r2.Vec
		screenScale, angle   float64
	)

	system.query.Each(e.World, func(entry *donburi.Entry) {
		sprite := graphics.Sprite.Get(entry)

		if entry.HasComponent(transform.Transform) {
			position = r2.Vec(transform.WorldPosition(entry))
			angle = transform.WorldRotation(entry) * 2 * math.Pi / 360
		} else if entry.HasComponent(physics.Body) {
			body := physics.Body.Get(entry)

			position = r2.Vec(body.Position())
			angle = body.Angle()
		}

		spriteDrawImageOptions.GeoM.Reset()

		// Convert to world coordinates
		screenSize = sprite.ScreenSize(system.display.VirtualResolution())
		screenScale = sprite.ScreenScale(system.display.VirtualResolution())

		spriteDrawImageOptions.GeoM.Scale(screenScale, screenScale)
		spriteDrawImageOptions.GeoM.Translate(-screenSize.X/2, -screenSize.Y/2)

		// Set direction
		if sprite.FlipY {
			spriteDrawImageOptions.GeoM.Scale(-1, 1)
		}
		if sprite.FlipX {
			spriteDrawImageOptions.GeoM.Scale(1, -1)
		}

		// Position in world coordinates
		spriteDrawImageOptions.GeoM.Rotate(angle)
		spriteDrawImageOptions.GeoM.Translate(system.display.WorldToScreen(position))

		system.camera.Camera.Draw(sprite.Image, spriteDrawImageOptions, screen)
	})
}

func (system *SpriteSystem) DrawDebug(e *ecs.ECS, screen *ebiten.Image) {
	if !system.settings.Settings().Debug.DrawPositions {
		return
	}

	var (
		position, velocity, worldSize r2.Vec
	)

	system.query.Each(e.World, func(entry *donburi.Entry) {
		sprite := graphics.Sprite.Get(entry)

		worldSize = sprite.WorldSize()

		var debugText string

		if entry.HasComponent(transform.Transform) {
			position = r2.Vec(transform.WorldPosition(entry))

			debugText = fmt.Sprintf("Transform: %.3f, %.3f", position.X, position.Y)
		} else if entry.HasComponent(physics.Body) {
			body := physics.Body.Get(entry)

			position = r2.Vec(body.Position())
			velocity = r2.Vec(body.Velocity())

			debugText = fmt.Sprintf("Transform: %.3f, %.3f\nVelocity: %.3f, %.3f",
				position.X, position.Y,
				velocity.X, velocity.Y)
		}

		metrics := system.settings.Settings().Debug.FontFace.Metrics()

		debugTextOptions := text.DrawOptions{
			LayoutOptions: text.LayoutOptions{
				LineSpacing: metrics.HAscent + metrics.HDescent + metrics.HLineGap,
			},
			DrawImageOptions: ebiten.DrawImageOptions{
				Filter: ebiten.FilterLinear,
			},
		}

		textPosition := r2.Vec{X: position.X + worldSize.X/2, Y: position.Y + worldSize.Y/2}

		debugTextOptions.GeoM.Reset()
		debugTextOptions.GeoM.Rotate(-system.camera.Camera.Rotation() * 2 * math.Pi / 360)
		debugTextOptions.GeoM.Translate(system.display.WorldToScreen(textPosition))

		system.camera.Camera.ApplyCameraTransform(&debugTextOptions.GeoM)

		text.Draw(screen, debugText, system.settings.Settings().Debug.FontFace, &debugTextOptions)

		// Center dot
		debugDotOptions := ebiten.DrawImageOptions{
			Filter: ebiten.FilterLinear,
		}

		debugDotOptions.GeoM.Translate(system.display.WorldToScreen(position))

		system.camera.Camera.ApplyCameraTransform(&debugDotOptions.GeoM)

		d2d.Dot(screen, &debugDotOptions, colornames.Yellow)
	})

	dotOptions := ebiten.DrawImageOptions{}
	system.camera.Camera.ApplyCameraTransform(&dotOptions.GeoM)

	d2d.Dot(screen, &dotOptions, colornames.Red)
}
