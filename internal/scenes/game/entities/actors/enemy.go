package actors

import (
	"github.com/solarlune/resolv"
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

var EnemyTag = donburi.NewTag().SetName("Enemy")

var Enemy = archetypes.NewArchetype(
	EnemyTag,
	transform.Transform,
	geometry.Shape,
	visuals.Sprite,
	geometry.Velocity,
)

func CreateEnemy(ecs *ecs.ECS, registry *resources.Registry, scaler utility.Scaler, position, velocity r2.Vec) *donburi.Entry {
	entry := Enemy.Spawn(ecs, layers.Game)

	sprite := registry.LoadImage(assets.Submarine)

	size := r2.Vec{X: float64(sprite.Data.Bounds().Size().X), Y: float64(sprite.Data.Bounds().Size().Y)}

	scale := scaler.GetNormalizedScale(r2.Vec{X: float64(sprite.Data.Bounds().Size().X), Y: float64(sprite.Data.Bounds().Size().Y)})

	visuals.Sprite.SetValue(entry, visuals.SpriteData{
		Image: sprite.Data,
	})
	transform.Transform.SetValue(entry, transform.TransformData{
		LocalPosition: math.Vec2(position),
		LocalScale:    math.NewVec2(scale, scale),
		LocalRotation: 0,
	})
	geometry.Velocity.SetValue(entry, velocity)

	shape := *resolv.NewRectangle(position.X-size.X/2, position.Y-size.Y/2, size.X, size.Y)
	shape.SetScale(scale, scale)
	geometry.Shape.SetValue(entry, shape)

	return entry
}
