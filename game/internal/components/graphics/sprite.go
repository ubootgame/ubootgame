package graphics

import (
	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/ubootgame/ubootgame/framework/services/display"
	"github.com/yohamta/donburi"
	"gonum.org/v1/gonum/spatial/r2"
)

type SpriteData struct {
	Image        *ebiten.Image
	FlipX, FlipY bool
	size         r2.Vec
	scale        display.Scale
}

var Sprite = donburi.NewComponentType[SpriteData]()

func NewSprite(image resource.Image, scale display.Scale, flipX, flipY bool) SpriteData {
	size := image.Data.Bounds().Size()
	return SpriteData{
		Image: image.Data,
		FlipX: flipX,
		FlipY: flipY,
		size:  r2.Vec{X: float64(size.X), Y: float64(size.Y)},
		scale: scale,
	}
}

func (s *SpriteData) ScreenScale(resolution r2.Vec) float64 {
	return s.scale.ScreenScale(s.size, resolution)
}

func (s *SpriteData) ScreenSize(resolution r2.Vec) r2.Vec {
	screenScale := s.ScreenScale(resolution)
	return r2.Vec{X: screenScale * s.size.X, Y: screenScale * s.size.Y}
}

func (s *SpriteData) WorldSize() r2.Vec {
	return s.scale.WorldSize(s.size)
}
