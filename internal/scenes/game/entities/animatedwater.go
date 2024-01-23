package entities

import (
	"github.com/solarlune/goaseprite"
	"github.com/ubootgame/ubootgame/assets"
	assets2 "github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var AnimatedWaterTag = donburi.NewTag().SetName("Animated Water")

var AnimatedWater = utility.NewArchetype(
	AnimatedWaterTag,
	components.AnimatedSprite,
)

func CreateAnimatedWater(ecs *ecs.ECS, registry *resources.Registry) *donburi.Entry {
	animatedWater := AnimatedWater.Spawn(ecs)
	json, _ := assets.FS.ReadFile("water/water.json")
	sprite := goaseprite.Read(json)
	player := sprite.CreatePlayer()
	components.AnimatedSprite.SetValue(animatedWater, components.AnimatedSpriteData{
		Sprite: sprite,
		Player: player,
		Image:  registry.LoadImage(assets2.ImageAnimatedWater).Data,
	})

	player.Play("")

	return animatedWater
}
