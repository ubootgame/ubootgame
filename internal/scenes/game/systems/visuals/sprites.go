package visuals

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility/draw"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/injector"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/systems"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/colornames"
	"gonum.org/v1/gonum/spatial/r2"
	"math"
)

type SpriteSystem struct {
	systems.BaseSystem

	camera  *game_system.CameraData
	debug   *game_system.DebugData
	display *game_system.DisplayData

	query *donburi.Query

	spriteDrawImageOptions *ebiten.DrawImageOptions

	debugText             string
	debugTextOptions      *text.DrawOptions
	debugTextPositionOpts *ebiten.DrawImageOptions
}

func NewSpriteSystem() *SpriteSystem {
	system := &SpriteSystem{
		query:                  donburi.NewQuery(filter.Contains(visuals.Sprite, transform.Transform)),
		spriteDrawImageOptions: &ebiten.DrawImageOptions{},
		debugTextOptions: &text.DrawOptions{
			DrawImageOptions: ebiten.DrawImageOptions{
				Filter: ebiten.FilterLinear,
			},
		},
		debugTextPositionOpts: &ebiten.DrawImageOptions{},
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

func (system *SpriteSystem) Layers() []lo.Tuple2[ecs.LayerID, systems.Renderer] {
	return []lo.Tuple2[ecs.LayerID, systems.Renderer]{
		{A: layers.Game, B: system.Draw},
		{A: layers.Debug, B: system.DrawDebug},
	}
}

func (system *SpriteSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	system.query.Each(e.World, func(entry *donburi.Entry) {
		sprite := visuals.Sprite.Get(entry)

		worldRotation := transform.WorldRotation(entry)
		worldScale := transform.WorldScale(entry)
		worldPosition := transform.WorldPosition(entry)

		system.spriteDrawImageOptions.GeoM.Reset()

		if sprite.FlipX {
			system.spriteDrawImageOptions.GeoM.Scale(1, -1)
			//system.spriteDrawImageOptions.GeoM.Translate(0, float64(sprite.Image.Bounds().Size().Y))
		}
		if sprite.FlipY {
			system.spriteDrawImageOptions.GeoM.Scale(-1, 1)
			//system.spriteDrawImageOptions.GeoM.Translate(float64(sprite.Image.Bounds().Size().X), 0)
		}

		system.spriteDrawImageOptions.GeoM.Translate(-float64(sprite.Image.Bounds().Size().X/2), -float64(sprite.Image.Bounds().Size().Y/2))
		system.spriteDrawImageOptions.GeoM.Rotate(float64(worldRotation) * 2 * math.Pi / 360)
		system.spriteDrawImageOptions.GeoM.Scale(worldScale.X, worldScale.Y)
		system.spriteDrawImageOptions.GeoM.Translate(worldPosition.X, worldPosition.Y)
		system.spriteDrawImageOptions.GeoM.Concat(*system.camera.Matrix)

		system.spriteDrawImageOptions.Filter = ebiten.FilterLinear

		screen.DrawImage(sprite.Image, system.spriteDrawImageOptions)
	})
}

func (system *SpriteSystem) DrawDebug(e *ecs.ECS, screen *ebiten.Image) {
	if !system.debug.DrawPositions {
		return
	}

	system.query.Each(e.World, func(entry *donburi.Entry) {
		worldPosition := transform.WorldPosition(entry)

		// Center dot
		spriteCenter := system.camera.WorldToScreenPosition(r2.Vec(worldPosition))
		draw.BigDot(screen, spriteCenter, colornames.Yellow)

		// Debug text
		debugText := fmt.Sprintf("Transform: %.3f, %.3f",
			worldPosition.X, worldPosition.Y)

		if entry.HasComponent(geometry.Velocity) {
			velocity := geometry.Velocity.Get(entry)
			debugText += fmt.Sprintf("\nVelocity: %.3f, %.3f", velocity.X, velocity.Y)
		}

		spriteBottomRight := system.camera.WorldToScreenPosition(r2.Vec{
			X: worldPosition.X,
			Y: worldPosition.Y,
		})

		system.debugTextOptions.GeoM.Reset()
		system.debugTextOptions.GeoM.Translate(spriteBottomRight.X, spriteBottomRight.Y)

		metrics := system.debug.FontFace.Metrics()
		system.debugTextOptions.LineSpacing = metrics.HAscent + metrics.HDescent + metrics.HLineGap

		text.Draw(screen, debugText, system.debug.FontFace, system.debugTextOptions)
	})
}
