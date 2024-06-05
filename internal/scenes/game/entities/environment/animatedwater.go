package environment

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/archetypes"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var AnimatedWaterTag = donburi.NewTag().SetName("Animated Game")

var AnimatedWater = archetypes.NewArchetype(
	AnimatedWaterTag,
	visuals.Aseprite,
)

func CreateAnimatedWater(ecs *ecs.ECS, registry *resources.Registry) *donburi.Entry {
	entry := AnimatedWater.Spawn(ecs, layers.Game)

	aseprite := registry.LoadAseprite(assets.AnimatedWater)
	visuals.Aseprite.SetValue(entry, visuals.AsepriteData{Aseprite: aseprite, Speed: 0.5})

	_ = aseprite.Player.Play("")

	return entry
}
