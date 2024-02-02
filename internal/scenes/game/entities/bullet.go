package entities

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"gonum.org/v1/gonum/spatial/r2"
)

var BulletTag = donburi.NewTag().SetName("Bullet")

var Bullet = utility.NewArchetype(
	BulletTag,
	components.Shape,
	components.Transform,
	components.Velocity,
)

func CreateBullet(ecs *ecs.ECS, from, to r2.Vec) *donburi.Entry {
	entry := Bullet.Spawn(ecs, layers.Foreground)

	direction := r2.Sub(to, from)
	velocity := r2.Unit(direction)

	components.Transform.SetValue(entry, components.TransformData{
		Center: from,
	})
	components.Velocity.SetValue(entry, r2.Scale(1, velocity))

	return entry
}
