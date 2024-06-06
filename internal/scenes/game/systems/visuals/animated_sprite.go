package visuals

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/internal/config"
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
	"image"
)

type AnimatedSpriteSystem struct {
	systems.BaseSystem

	camera *game_system.CameraData

	updateQuery *donburi.Query
	drawQuery   *donburi.Query

	spriteDrawImageOptions *ebiten.DrawImageOptions
}

func NewAnimatedSpriteSystem() *AnimatedSpriteSystem {
	system := &AnimatedSpriteSystem{
		updateQuery:            donburi.NewQuery(filter.Contains(visuals.AnimatedSprite)),
		drawQuery:              donburi.NewQuery(filter.Contains(visuals.AnimatedSprite, geometry.Transform)),
		spriteDrawImageOptions: &ebiten.DrawImageOptions{},
	}
	system.Injector = injector.NewInjector([]injector.Injection{
		injector.Once([]injector.Injection{
			injector.Component(&system.camera, game_system.Camera),
		}),
	})
	return system
}

func (system *AnimatedSpriteSystem) Layers() []lo.Tuple2[ecs.LayerID, systems.Renderer] {
	return []lo.Tuple2[ecs.LayerID, systems.Renderer]{
		{A: layers.Game, B: system.Draw},
		{A: layers.Debug, B: system.DrawDebug},
	}
}

func (system *AnimatedSpriteSystem) Update(e *ecs.ECS) {
	system.BaseSystem.Update(e)

	system.updateQuery.Each(e.World, func(entry *donburi.Entry) {
		animatedSprite := visuals.AnimatedSprite.Get(entry)
		animatedSprite.Aseprite.Player.Update(1.0 / float32(config.C.TargetTPS) * animatedSprite.Speed)
	})
}

func (system *AnimatedSpriteSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	system.drawQuery.Each(e.World, func(entry *donburi.Entry) {
		animatedSprite := visuals.AnimatedSprite.Get(entry)
		transform := geometry.Transform.Get(entry)

		system.spriteDrawImageOptions.GeoM.Reset()

		system.spriteDrawImageOptions.GeoM.Translate(-float64(animatedSprite.Aseprite.Player.File.FrameWidth/2), -float64(animatedSprite.Aseprite.Player.File.FrameHeight/2))
		system.spriteDrawImageOptions.GeoM.Scale(animatedSprite.Scale*transform.Size.X, animatedSprite.Scale*transform.Size.X)
		system.spriteDrawImageOptions.GeoM.Translate(transform.Center.X, transform.Center.Y)
		system.spriteDrawImageOptions.GeoM.Concat(*system.camera.Matrix)

		system.spriteDrawImageOptions.Filter = ebiten.FilterLinear

		sub := animatedSprite.Aseprite.Image.SubImage(image.Rect(animatedSprite.Aseprite.Player.CurrentFrameCoords()))

		screen.DrawImage(sub.(*ebiten.Image), system.spriteDrawImageOptions)
	})
}

func (system *AnimatedSpriteSystem) DrawDebug(e *ecs.ECS, screen *ebiten.Image) {
	system.drawQuery.Each(e.World, func(entry *donburi.Entry) {
		transform := geometry.Transform.Get(entry)

		spriteCenter := system.camera.WorldToScreenPosition(transform.Center)
		draw.BigDot(screen, spriteCenter, colornames.Green)
	})
}
