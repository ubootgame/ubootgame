package visuals

import (
	"github.com/ubootgame/ubootgame/framework/services/resources"
	"github.com/yohamta/donburi"
)

type AnimatedSpriteData struct {
	Aseprite resources.Aseprite
	Speed    float32
}

var AnimatedSprite = donburi.NewComponentType[AnimatedSpriteData]()
