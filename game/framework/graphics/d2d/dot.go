package d2d

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

var dotImg *ebiten.Image

func Dot(screen *ebiten.Image, options *ebiten.DrawImageOptions, drawColor color.Color) {
	if dotImg == nil {
		dotImg = ebiten.NewImage(4, 4)
		dotImg.Fill(color.White)
	} else {
		dotImg.Fill(color.White)
	}

	options.ColorScale.ScaleWithColor(drawColor)
	screen.DrawImage(dotImg, options)
}
