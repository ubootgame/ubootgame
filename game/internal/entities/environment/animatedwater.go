package environment

import (
	"github.com/ubootgame/ubootgame/framework"
	ecsFramework "github.com/ubootgame/ubootgame/framework/game/ecs"
	"github.com/ubootgame/ubootgame/framework/services/display"
	"github.com/ubootgame/ubootgame/framework/services/resources"
	"github.com/ubootgame/ubootgame/internal/components/graphics"
	"github.com/ubootgame/ubootgame/internal/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
	"gonum.org/v1/gonum/spatial/r2"
)

var AnimatedWaterTag = donburi.NewTag().SetName("Animated Game")

var AnimatedWater = ecsFramework.NewArchetype(
	AnimatedWaterTag,
	transform.Transform,
	graphics.AnimatedSprite,
)

func CreateAnimatedWater(ecs *ecs.ECS, resources framework.ResourceService, asepriteID resources.AsepriteID, scale display.Scale, position r2.Vec) *donburi.Entry {
	entry := AnimatedWater.Spawn(ecs, layers.Game)

	aseprite := resources.LoadAseprite(asepriteID)

	graphics.AnimatedSprite.SetValue(entry, graphics.AnimatedSpriteData{Aseprite: aseprite, Speed: 0.5})

	//_, _, localScale := scale.GetNormalizedSizeAndScale(r2.Vec{X: float64(aseprite.Player.File.FrameWidth), Y: float64(aseprite.Player.File.Height)})
	//
	//transform.Transform.SetValue(entry, transform.TransformData{
	//	LocalPosition: math.Vec2(position),
	//	LocalScale:    math.NewVec2(localScale, localScale),
	//	LocalRotation: 0,
	//})

	_ = aseprite.Player.Play("")

	return entry
}
