package environment

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/environment"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/injector"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/systems"
	"github.com/yohamta/donburi/ecs"
	"image"
)

type AnimatedWaterSystem struct {
	systems.BaseSystem

	aseprite *visuals.AsepriteData
}

func NewAnimatedWaterSystem() *AnimatedWaterSystem {
	system := &AnimatedWaterSystem{}
	system.Injector = injector.NewInjector([]injector.Injection{
		injector.WithTag(environment.AnimatedWaterTag, []injector.Injection{
			injector.Component(&system.aseprite, visuals.Aseprite),
		}),
	})
	return system
}

func (system *AnimatedWaterSystem) Layers() []lo.Tuple2[ecs.LayerID, systems.Renderer] {
	return []lo.Tuple2[ecs.LayerID, systems.Renderer]{
		{A: layers.Game, B: system.Draw},
	}
}

func (system *AnimatedWaterSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	system.BaseSystem.Update(e)

	sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(sw/2, sh/2)

	sub := system.aseprite.Aseprite.Image.SubImage(image.Rect(system.aseprite.Aseprite.Player.CurrentFrameCoords()))

	screen.DrawImage(sub.(*ebiten.Image), opts)
}
