package environment

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/archetypes"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"gonum.org/v1/gonum/spatial/r2"
)

var AnimatedWaterTag = donburi.NewTag().SetName("Animated Game")

var AnimatedWater = archetypes.NewArchetype(
	AnimatedWaterTag,
	visuals.AnimatedSprite,
	geometry.Transform,
)

func CreateAnimatedWater(ecs *ecs.ECS, registry *resources.Registry, scaler utility.Scaler, position r2.Vec) *donburi.Entry {
	entry := AnimatedWater.Spawn(ecs, layers.Game)

	aseprite := registry.LoadAseprite(assets.AnimatedWater)

	size, scale := scaler.GetNormalSizeAndScale(r2.Vec{X: float64(aseprite.Player.File.FrameWidth), Y: float64(aseprite.Player.File.Height)})

	visuals.AnimatedSprite.SetValue(entry, visuals.AnimatedSpriteData{Aseprite: aseprite, Speed: 0.5, Scale: scale})

	geometry.Transform.SetValue(entry, geometry.TransformData{
		Center: position,
		Size:   size,
	})

	_ = aseprite.Player.Play("")

	return entry
}
