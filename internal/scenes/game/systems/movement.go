package systems

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type movementSystem struct {
	query *donburi.Query
}

var Movement = movementSystem{
	query: donburi.NewQuery(filter.Contains(components.Position, components.Velocity)),
}

func (system *movementSystem) Update(e *ecs.ECS) {
	system.query.Each(e.World, func(entry *donburi.Entry) {
		velocity := components.Velocity.Get(entry)
		position := components.Position.Get(entry)

		position.Center.X += velocity.X
		position.Center.Y += velocity.Y
	})
}
