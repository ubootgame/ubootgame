package components

import (
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	"github.com/yohamta/donburi"
)

type AsepriteData struct {
	Aseprite resources.Aseprite
	Speed    float32
}

var Aseprite = donburi.NewComponentType[AsepriteData]()
