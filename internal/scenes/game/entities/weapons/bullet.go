package weapons

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/archetypes"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"gonum.org/v1/gonum/spatial/r2"
)

var BulletTag = donburi.NewTag().SetName("Bullet")

var Bullet = archetypes.NewArchetype(
	BulletTag,
	geometry.Shape,
	geometry.Transform,
	geometry.Velocity,
)

func CreateBullet(ecs *ecs.ECS, from, to r2.Vec) *donburi.Entry {
	entry := Bullet.Spawn(ecs, layers.Foreground)

	direction := r2.Sub(to, from)
	velocity := r2.Unit(direction)

	geometry.Transform.SetValue(entry, geometry.TransformData{
		Center: from,
	})
	geometry.Velocity.SetValue(entry, r2.Scale(1, velocity))

	return entry
}
