package draw

import (
	"github.com/hajimehoshi/ebiten/v2"
	"gonum.org/v1/gonum/spatial/r2"
	"image/color"
)

var dotImg *ebiten.Image

func Dot(screen *ebiten.Image, position r2.Vec, drawColor color.Color) {
	if dotImg == nil {
		dotImg = ebiten.NewImage(4, 4)
		dotImg.Fill(color.White)
	} else {
		dotImg.Fill(color.White)
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(position.X-2, position.Y-2)
	opt.ColorScale.ScaleWithColor(drawColor)
	screen.DrawImage(dotImg, opt)
}
