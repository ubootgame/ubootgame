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

	query      *donburi.Query
	debugQuery *donburi.Query

	spriteDrawImageOptions *ebiten.DrawImageOptions

	debugText             string
	debugTextOptions      *text.DrawOptions
	debugTextPositionOpts *ebiten.DrawImageOptions
}

func NewSpriteSystem() *SpriteSystem {
	system := &SpriteSystem{
		query:                  donburi.NewQuery(filter.Contains(visuals.Sprite, transform.Transform)),
		debugQuery:             donburi.NewQuery(filter.Contains(visuals.Sprite, transform.Transform, geometry.Velocity)),
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
		scale := geometry.Scale.Get(entry)

		worldRotation := transform.WorldRotation(entry)
		worldScale := transform.WorldScale(entry)
		worldPosition := transform.WorldPosition(entry)

		system.spriteDrawImageOptions.GeoM.Reset()

		// Convert to world coordinates
		system.spriteDrawImageOptions.GeoM.Scale(scale.NormalizedScale, scale.NormalizedScale)
		system.spriteDrawImageOptions.GeoM.Translate(-scale.NormalizedSize.X/2, -scale.NormalizedSize.Y/2)

		// Set direction
		if entry.HasComponent(geometry.Direction) {
			direction := geometry.Direction.Get(entry)

			if direction.Horizontal != direction.HorizontalBase {
				system.spriteDrawImageOptions.GeoM.Scale(-1, 1)
			}
			if direction.Vertical != direction.VerticalBase {
				system.spriteDrawImageOptions.GeoM.Scale(1, -1)
			}
		}

		// Position in world coordinates
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

	system.debugQuery.Each(e.World, func(entry *donburi.Entry) {
		velocity := geometry.Velocity.Get(entry)
		scale := geometry.Scale.Get(entry)

		worldPosition := transform.WorldPosition(entry)
		worldScale := transform.WorldScale(entry)

		// Center dot
		spriteCenter := system.camera.WorldToScreenPosition(r2.Vec(worldPosition))
		draw.BigDot(screen, spriteCenter, colornames.Yellow)

		// Debug text
		debugText := fmt.Sprintf("Transform: %.3f, %.3f\nVelocity: %.3f, %.3f",
			worldPosition.X, worldPosition.Y,
			velocity.X, velocity.Y)

		spriteBottomRight := system.camera.WorldToScreenPosition(r2.Vec{
			X: worldPosition.X + (scale.NormalizedSize.X*worldScale.X)/2,
			Y: worldPosition.Y + (scale.NormalizedSize.Y*worldScale.Y)/2,
		})

		metrics := system.debug.FontFace.Metrics()
		system.debugTextOptions.LineSpacing = metrics.HAscent + metrics.HDescent + metrics.HLineGap
		system.debugTextOptions.GeoM.Reset()
		system.debugTextOptions.GeoM.Translate(spriteBottomRight.X, spriteBottomRight.Y)

		text.Draw(screen, debugText, system.debug.FontFace, system.debugTextOptions)
	})
}
