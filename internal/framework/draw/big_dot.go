package draw

import (
	"github.com/hajimehoshi/ebiten/v2"
	"gonum.org/v1/gonum/spatial/r2"
	"image/color"
)

var bigDotImg *ebiten.Image

func BigDot(screen *ebiten.Image, position r2.Vec, drawColor color.Color) {
	if bigDotImg == nil {
		bigDotImg = ebiten.NewImage(4, 4)
		bigDotImg.Fill(color.White)
	} else {
		bigDotImg.Fill(color.White)
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(position.X-2, position.Y-2)
	opt.ColorScale.ScaleWithColor(drawColor)
	screen.DrawImage(bigDotImg, opt)
}
