package environment

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/environment"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"image"
)

type animatedWaterSystem struct {
	entry *donburi.Entry
}

var AnimatedWater = &animatedWaterSystem{}

func (system *animatedWaterSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	if system.entry == nil {
		var ok bool
		if system.entry, ok = environment.AnimatedWaterTag.First(e.World); !ok {
			panic("no animated water found")
		}
	}

	aseprite := visuals.Aseprite.Get(system.entry)

	sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(sw/2, sh/2)

	sub := aseprite.Aseprite.Image.SubImage(image.Rect(aseprite.Aseprite.Player.CurrentFrameCoords()))

	screen.DrawImage(sub.(*ebiten.Image), opts)
}
