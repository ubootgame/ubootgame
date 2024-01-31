package systems

import (
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type movementSystem struct {
	query *donburi.Query
}

var Movement = movementSystem{
	query: donburi.NewQuery(filter.Contains(components.Transform, components.Velocity)),
}

func (system *movementSystem) Update(e *ecs.ECS) {
	system.query.Each(e.World, func(entry *donburi.Entry) {
		velocity := components.Velocity.Get(entry)
		transform := components.Transform.Get(entry)

		transform.Center.X += velocity.X / float64(config.C.TargetTPS)
		transform.Center.Y += velocity.Y / float64(config.C.TargetTPS)

		if velocity.X < 0 {
			transform.FlipY = true
		} else if velocity.X > 0 {
			transform.FlipY = false
		}

		if velocity.Y < 0 {
			transform.FlipX = true
		} else if velocity.Y > 0 {
			transform.FlipX = false
		}
	})
}
