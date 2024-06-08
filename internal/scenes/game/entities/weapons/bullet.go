package weapons

import (
	"github.com/ubootgame/ubootgame/internal/framework"
	"github.com/ubootgame/ubootgame/internal/framework/ecs/archetypes"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"gonum.org/v1/gonum/spatial/r2"
)

var BulletTag = donburi.NewTag().SetName("Bullet")

var Bullet = archetypes.NewArchetype(
	BulletTag,
	transform.Transform,
	geometry.Velocity,
)

func CreateBullet(ecs *ecs.ECS, from, to r2.Vec) *donburi.Entry {
	entry := Bullet.Spawn(ecs, layers.Game)

	direction := r2.Sub(to, from)
	velocity := r2.Unit(direction)

	transform.Transform.SetValue(entry, transform.TransformData{
		LocalPosition: math.Vec2(from),
		LocalScale:    math.NewVec2(1, 1),
		LocalRotation: 0,
	})
	geometry.Velocity.SetValue(entry, r2.Scale(framework.WorldSizeBase, velocity))

	return entry
}
