package environment

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/environment"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type waterSystem struct {
	displayEntry, waterEntry *donburi.Entry
}

var Water = &waterSystem{}

func (system *waterSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	var ok bool
	if system.displayEntry == nil {
		if system.displayEntry, ok = game_system.Display.First(e.World); !ok {
			panic("no display found")
		}
	}
	if system.waterEntry == nil {
		if system.waterEntry, ok = environment.WaterTag.First(e.World); !ok {
			panic("no water found")
		}
	}

	display := game_system.Display.Get(system.displayEntry)
	sprite := visuals.Sprite.Get(system.waterEntry)

	sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
	w, h := float64(sprite.Image.Bounds().Size().X), float64(sprite.Image.Bounds().Dy())
	y := sh / 2

	sizeScale := 0.1 * (sh / h)

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(-w)/2, float64(-h)/2)
	op.GeoM.Scale(sizeScale*display.ScalingFactor, sizeScale*display.ScalingFactor)
	op.GeoM.Translate(0, y)
	op.ColorScale.ScaleAlpha(0.1)

	for x := 0; x <= int(sw/(w*sizeScale)+1); x++ {
		screen.DrawImage(sprite.Image, op)
		op.GeoM.Translate(w*sizeScale, 0)
	}
}
