package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/goaseprite"
	"github.com/yohamta/donburi"
)

type AnimatedSpriteData struct {
	Sprite *goaseprite.File
	Player *goaseprite.Player
	Image  *ebiten.Image
}

var AnimatedSprite = donburi.NewComponentType[AnimatedSpriteData]()
