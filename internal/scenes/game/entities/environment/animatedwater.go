package environment

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	ecsFramework "github.com/ubootgame/ubootgame/pkg/ecs"
	"github.com/ubootgame/ubootgame/pkg/resources"
	"github.com/ubootgame/ubootgame/pkg/world"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"gonum.org/v1/gonum/spatial/r2"
)

var AnimatedWaterTag = donburi.NewTag().SetName("Animated Game")

var AnimatedWater = ecsFramework.NewArchetype(
	AnimatedWaterTag,
	transform.Transform,
	visuals.AnimatedSprite,
)

func CreateAnimatedWater(ecs *ecs.ECS, registry *resources.Registry, scaler world.Scaler, position r2.Vec) *donburi.Entry {
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
