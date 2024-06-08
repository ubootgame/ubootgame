package systems

import (
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/systems"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

type MovementSystem struct {
	systems.BaseSystem

	query *donburi.Query
}

func NewMovementSystem() *MovementSystem {
	system := &MovementSystem{
		query: donburi.NewQuery(filter.Contains(transform.Transform, geometry.Velocity)),
	}
	return system
}

func (system *MovementSystem) Update(e *ecs.ECS) {
	system.BaseSystem.Update(e)

	system.query.Each(e.World, func(entry *donburi.Entry) {
		velocity := geometry.Velocity.Get(entry)
		t := transform.Transform.Get(entry)

		t.LocalPosition.X += velocity.X / float64(config.C.TargetTPS)
		t.LocalPosition.Y += velocity.Y / float64(config.C.TargetTPS)

		//if velocity.X < 0 {
		//	transform.FlipY = true
		//} else if velocity.X > 0 {
		//	transform.FlipY = false
		//}
		//
		//if velocity.Y < 0 {
		//	transform.FlipX = true
		//} else if velocity.Y > 0 {
		//	transform.FlipX = false
		//}
	})
}
