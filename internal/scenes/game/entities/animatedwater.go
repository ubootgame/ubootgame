package entities

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var AnimatedWaterTag = donburi.NewTag().SetName("Animated Background")

var AnimatedWater = utility.NewArchetype(
	AnimatedWaterTag,
	components.Aseprite,
)

func CreateAnimatedWater(ecs *ecs.ECS, registry *resources.Registry) *donburi.Entry {
	entry := AnimatedWater.Spawn(ecs, layers.Background)

	aseprite := registry.LoadAseprite(assets.AnimatedWater)
	components.Aseprite.SetValue(entry, components.AsepriteData{Aseprite: aseprite, Speed: 0.5})

	_ = aseprite.Player.Play("")

	return entry
}
