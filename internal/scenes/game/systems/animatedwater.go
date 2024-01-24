package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities"
	"github.com/yohamta/donburi/ecs"
	"image"
)

func DrawAnimatedWater(e *ecs.ECS, screen *ebiten.Image) {
	waterEntry, _ := entities.AnimatedWaterTag.First(e.World)

	spriteData := components.Aseprite.Get(waterEntry)

	sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(sw/2, sh/2)

	sub := spriteData.Aseprite.Image.SubImage(image.Rect(spriteData.Aseprite.Player.CurrentFrameCoords()))

	screen.DrawImage(sub.(*ebiten.Image), opts)
}
