package entities

import (
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
	components.Aseprite,
)

func CreateAnimatedWater(ecs *ecs.ECS, registry *resources.Registry) *donburi.Entry {
	water := AnimatedWater.Spawn(ecs)

	aseprite := registry.LoadAseprite(assets2.AnimatedWater)
	components.Aseprite.SetValue(water, components.AsepriteData{Aseprite: aseprite, Speed: 0.5})

	_ = aseprite.Player.Play("")

	return water
}
