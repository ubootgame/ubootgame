package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities"
	"github.com/yohamta/donburi/ecs"
)

func DrawWater(e *ecs.ECS, screen *ebiten.Image) {
	waterEntry, _ := entities.WaterTag.First(e.World)

	spriteData := components.Sprite.Get(waterEntry)

	sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
	w := float64(spriteData.Image.Bounds().Size().X)
	y := sh / 2

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(0, y)
	op.ColorScale.ScaleAlpha(0.1)

	for x := 0; x <= int(sw/w); x++ {
		screen.DrawImage(spriteData.Image, op)
		op.GeoM.Translate(w, 0)
	}
}
