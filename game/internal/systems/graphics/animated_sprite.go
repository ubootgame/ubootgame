package graphics

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/framework"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/components"
	"github.com/ubootgame/ubootgame/internal/components/graphics"
	"github.com/ubootgame/ubootgame/internal/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"image"
)

type AnimatedSpriteSystem struct {
	settings framework.SettingsService[internal.Settings]

	camera *components.CameraData

	updateQuery *donburi.Query
	drawQuery   *donburi.Query

	spriteDrawImageOptions *ebiten.DrawImageOptions
}

func NewAnimatedSpriteSystem(settings framework.SettingsService[internal.Settings]) *AnimatedSpriteSystem {
	return &AnimatedSpriteSystem{
		settings:               settings,
		updateQuery:            donburi.NewQuery(filter.Contains(graphics.AnimatedSprite)),
		drawQuery:              donburi.NewQuery(filter.Contains(graphics.AnimatedSprite, transform.Transform)),
		spriteDrawImageOptions: &ebiten.DrawImageOptions{},
	}
}

func (system *AnimatedSpriteSystem) Layers() []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer]{
		{A: layers.Game, B: system.Draw},
		{A: layers.Debug, B: system.DrawDebug},
	}
}

func (system *AnimatedSpriteSystem) Update(e *ecs.ECS) {
	if entry, found := components.Camera.First(e.World); found {
		system.camera = components.Camera.Get(entry)
	}

	system.updateQuery.Each(e.World, func(entry *donburi.Entry) {
		animatedSprite := graphics.AnimatedSprite.Get(entry)
		animatedSprite.Aseprite.Player.Update(1.0 / float32(system.settings.Settings().Internals.TPS) * animatedSprite.Speed)
	})
}

func (system *AnimatedSpriteSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	system.drawQuery.Each(e.World, func(entry *donburi.Entry) {
		animatedSprite := graphics.AnimatedSprite.Get(entry)

		worldPosition := transform.WorldPosition(entry)
		worldScale := transform.WorldScale(entry)

		system.spriteDrawImageOptions.GeoM.Reset()

		system.spriteDrawImageOptions.GeoM.Translate(-float64(animatedSprite.Aseprite.Player.File.FrameWidth/2), -float64(animatedSprite.Aseprite.Player.File.FrameHeight/2))
		system.spriteDrawImageOptions.GeoM.Scale(worldScale.X, worldScale.Y)
		system.spriteDrawImageOptions.GeoM.Translate(worldPosition.X, worldPosition.Y)

		system.spriteDrawImageOptions.Filter = ebiten.FilterLinear

		sub := animatedSprite.Aseprite.Image.SubImage(image.Rect(animatedSprite.Aseprite.Player.CurrentFrameCoords()))

		system.camera.Camera.Draw(sub.(*ebiten.Image), system.spriteDrawImageOptions, screen)
	})
}

func (system *AnimatedSpriteSystem) DrawDebug(e *ecs.ECS, screen *ebiten.Image) {
	//system.drawQuery.Each(e.World, func(entry *donburi.Entry) {
	//	worldPosition := transform.WorldPosition(entry)
	//	spriteCenter := system.camera.WorldToScreenPosition(r2.Vec(worldPosition))
	//
	//	d2d.Dot(screen, spriteCenter, colornames.Green)
	//})
}
