package environment

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/environment"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/injector"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/systems"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type WaterSystem struct {
	systems.BaseSystem
	displayEntry, waterEntry *donburi.Entry

	display *game_system.DisplayData
	sprite  *visuals.SpriteData
}

func NewWaterSystem() *WaterSystem {
	system := &WaterSystem{}
	system.Injector = injector.NewInjector([]injector.Injection{
		injector.Once([]injector.Injection{
			injector.Component(&system.display, game_system.Display),
		}),
		injector.WithTag(environment.WaterTag, []injector.Injection{
			injector.Component(&system.sprite, visuals.Sprite),
		}),
	})
	return system
}

func (system *WaterSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	system.BaseSystem.Update(e)

	sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
	w, h := float64(system.sprite.Image.Bounds().Size().X), float64(system.sprite.Image.Bounds().Dy())
	y := sh / 2

	sizeScale := 0.1 * (sh / h)

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(-w)/2, float64(-h)/2)
	op.GeoM.Scale(sizeScale*system.display.ScalingFactor, sizeScale*system.display.ScalingFactor)
	op.GeoM.Translate(0, y)
	op.ColorScale.ScaleAlpha(0.1)

	for x := 0; x <= int(sw/(w*sizeScale)+1); x++ {
		screen.DrawImage(system.sprite.Image, op)
		op.GeoM.Translate(w*sizeScale, 0)
	}
}
