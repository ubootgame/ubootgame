package graphics

import (
	"github.com/ubootgame/ubootgame/framework/services/resources"
	"github.com/yohamta/donburi"
	"gonum.org/v1/gonum/spatial/r2"
)

type AnimatedSpriteData struct {
	Aseprite        resources.Aseprite
	Speed           float32
	NormalizedSize  r2.Vec
	NormalizedScale float64
}

var AnimatedSprite = donburi.NewComponentType[AnimatedSpriteData]()
