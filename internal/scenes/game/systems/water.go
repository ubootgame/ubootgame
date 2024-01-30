package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type waterSystem struct {
	entry *donburi.Entry
}

var Water = &waterSystem{}

func (system *waterSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	if system.entry == nil {
		var ok bool
		if system.entry, ok = entities.WaterTag.First(e.World); !ok {
			panic("no water found")
		}
	}

	sprite := components.Sprite.Get(system.entry)

	sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
	w, h := float64(sprite.Image.Bounds().Size().X), float64(sprite.Image.Bounds().Dy())
	y := sh / 2

	sizeScale := 0.1 * (sh / h)
	deviceScale := ebiten.DeviceScaleFactor()

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(-w)/2, float64(-h)/2)
	op.GeoM.Scale(sizeScale*deviceScale, sizeScale*deviceScale)
	op.GeoM.Translate(0, y)
	op.ColorScale.ScaleAlpha(0.1)

	for x := 0; x <= int(sw/(w*sizeScale)+1); x++ {
		screen.DrawImage(sprite.Image, op)
		op.GeoM.Translate(w*sizeScale, 0)
	}
}
