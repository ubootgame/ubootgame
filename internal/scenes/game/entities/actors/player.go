package actors

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
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"gonum.org/v1/gonum/spatial/r2"
)

var PlayerTag = donburi.NewTag().SetName("Player")

var Player = archetypes.NewArchetype(
	PlayerTag,
	transform.Transform,
	visuals.Sprite,
	geometry.Velocity,
)

func CreatePlayer(ecs *ecs.ECS, registry *resources.Registry, scaler utility.Scaler) *donburi.Entry {
	entry := Player.Spawn(ecs, layers.Game)

	sprite := registry.LoadImage(assets.Battleship)

	scale := scaler.GetNormalizedScale(r2.Vec{X: float64(sprite.Data.Bounds().Size().X), Y: float64(sprite.Data.Bounds().Size().Y)})

	visuals.Sprite.SetValue(entry, visuals.SpriteData{
		Image: sprite.Data,
	})
	transform.Transform.SetValue(entry, transform.TransformData{
		LocalPosition: math.NewVec2(0, 0),
		LocalScale:    math.NewVec2(scale, scale),
		LocalRotation: 0,
	})

	return entry
}
