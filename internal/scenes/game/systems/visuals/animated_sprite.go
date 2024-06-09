package visuals

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/framework"
	"github.com/ubootgame/ubootgame/internal/framework/draw"
	ecs2 "github.com/ubootgame/ubootgame/internal/framework/ecs"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/colornames"
	"gonum.org/v1/gonum/spatial/r2"
	"image"
)

type AnimatedSpriteSystem struct {
	ecs2.System

	settings *internal.Settings

	camera *framework.Camera

	updateQuery *donburi.Query
	drawQuery   *donburi.Query

	spriteDrawImageOptions *ebiten.DrawImageOptions
}

func NewAnimatedSpriteSystem(settings *internal.Settings) *AnimatedSpriteSystem {
	return &AnimatedSpriteSystem{
		settings:               settings,
		updateQuery:            donburi.NewQuery(filter.Contains(visuals.AnimatedSprite)),
		drawQuery:              donburi.NewQuery(filter.Contains(visuals.AnimatedSprite, transform.Transform)),
		spriteDrawImageOptions: &ebiten.DrawImageOptions{},
	}
}

func (system *AnimatedSpriteSystem) Layers() []lo.Tuple2[ecs.LayerID, ecs2.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecs2.Renderer]{
		{A: layers.Game, B: system.Draw},
		{A: layers.Debug, B: system.DrawDebug},
	}
}

func (system *AnimatedSpriteSystem) Update(e *ecs.ECS) {
	system.System.Update(e)

	system.updateQuery.Each(e.World, func(entry *donburi.Entry) {
		animatedSprite := visuals.AnimatedSprite.Get(entry)
		animatedSprite.Aseprite.Player.Update(1.0 / float32(system.settings.TargetTPS) * animatedSprite.Speed)
	})
}

func (system *AnimatedSpriteSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	system.drawQuery.Each(e.World, func(entry *donburi.Entry) {
		animatedSprite := visuals.AnimatedSprite.Get(entry)

		worldPosition := transform.WorldPosition(entry)
		worldScale := transform.WorldScale(entry)

		system.spriteDrawImageOptions.GeoM.Reset()

		system.spriteDrawImageOptions.GeoM.Translate(-float64(animatedSprite.Aseprite.Player.File.FrameWidth/2), -float64(animatedSprite.Aseprite.Player.File.FrameHeight/2))
		system.spriteDrawImageOptions.GeoM.Scale(worldScale.X, worldScale.Y)
		system.spriteDrawImageOptions.GeoM.Translate(worldPosition.X, worldPosition.Y)
		system.camera.Apply(&system.spriteDrawImageOptions.GeoM)

		system.spriteDrawImageOptions.Filter = ebiten.FilterLinear

		sub := animatedSprite.Aseprite.Image.SubImage(image.Rect(animatedSprite.Aseprite.Player.CurrentFrameCoords()))

		screen.DrawImage(sub.(*ebiten.Image), system.spriteDrawImageOptions)
	})
}

func (system *AnimatedSpriteSystem) DrawDebug(e *ecs.ECS, screen *ebiten.Image) {
	system.drawQuery.Each(e.World, func(entry *donburi.Entry) {
		worldPosition := transform.WorldPosition(entry)
		spriteCenter := system.camera.WorldToScreenPosition(r2.Vec(worldPosition))

		draw.Dot(screen, spriteCenter, colornames.Green)
	})
}
