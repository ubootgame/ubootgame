package environment

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/archetypes"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"gonum.org/v1/gonum/spatial/r2"
)

var AnimatedWaterTag = donburi.NewTag().SetName("Animated Game")

var AnimatedWater = archetypes.NewArchetype(
	AnimatedWaterTag,
	transform.Transform,
	visuals.AnimatedSprite,
)

func CreateAnimatedWater(ecs *ecs.ECS, registry *resources.Registry, scaler utility.Scaler, position r2.Vec) *donburi.Entry {
	entry := AnimatedWater.Spawn(ecs, layers.Game)

	aseprite := registry.LoadAseprite(assets.AnimatedWater)

	_, _, localScale := scaler.GetNormalizedSizeAndScale(r2.Vec{X: float64(aseprite.Player.File.FrameWidth), Y: float64(aseprite.Player.File.Height)})

	visuals.AnimatedSprite.SetValue(entry, visuals.AnimatedSpriteData{Aseprite: aseprite, Speed: 0.5})

	transform.Transform.SetValue(entry, transform.TransformData{
		LocalPosition: math.Vec2(position),
		LocalScale:    math.NewVec2(localScale, localScale),
		LocalRotation: 0,
	})

	_ = aseprite.Player.Play("")

	return entry
}
