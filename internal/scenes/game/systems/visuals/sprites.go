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
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/colornames"
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
		query:                  donburi.NewQuery(filter.Contains(visuals.Sprite, geometry.Transform)),
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
		transform := geometry.Transform.Get(entry)

		system.spriteDrawImageOptions.GeoM.Reset()

		if transform.FlipX {
			system.spriteDrawImageOptions.GeoM.Scale(1, -1)
			system.spriteDrawImageOptions.GeoM.Translate(0, float64(sprite.Image.Bounds().Size().Y))
		}
		if transform.FlipY {
			system.spriteDrawImageOptions.GeoM.Scale(-1, 1)
			system.spriteDrawImageOptions.GeoM.Translate(float64(sprite.Image.Bounds().Size().X), 0)
		}

		system.spriteDrawImageOptions.GeoM.Translate(-float64(sprite.Image.Bounds().Size().X/2), -float64(sprite.Image.Bounds().Size().Y/2))
		system.spriteDrawImageOptions.GeoM.Scale(sprite.Scale*transform.Size.X, sprite.Scale*transform.Size.X)
		system.spriteDrawImageOptions.GeoM.Translate(transform.Center.X, transform.Center.Y)
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
		transform := geometry.Transform.Get(entry)

		spriteCenter := system.camera.WorldToScreenPosition(transform.Center)
		draw.BigDot(screen, spriteCenter, colornames.Yellow)

		debugText := fmt.Sprintf("Transform: %.3f, %.3f\nSize: %.3f, %.3f",
			transform.Center.X, transform.Center.Y,
			transform.Size.X, transform.Size.Y)

		if entry.HasComponent(geometry.Velocity) {
			velocity := geometry.Velocity.Get(entry)
			debugText += fmt.Sprintf("\nVelocity: %.3f, %.3f", velocity.X, velocity.Y)
		}

		metrics := system.debug.FontFace.Metrics()
		system.debugTextOptions.LineSpacing = metrics.HAscent + metrics.HDescent + metrics.HLineGap

		textImage := draw.TextWithOptions(debugText, system.debug.FontFace, system.debugTextOptions)

		system.debugTextPositionOpts.GeoM.Reset()
		system.debugTextPositionOpts.GeoM.Scale(1/(system.display.VirtualResolution.X*system.camera.ZoomFactor), 1/(system.display.VirtualResolution.X*system.camera.ZoomFactor))
		system.debugTextPositionOpts.GeoM.Translate(transform.Center.X+transform.Size.X/2, transform.Center.Y+transform.Size.Y/2)
		system.debugTextPositionOpts.GeoM.Concat(*system.camera.Matrix)

		screen.DrawImage(textImage, system.debugTextPositionOpts)
	})
}
