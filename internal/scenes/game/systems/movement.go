package systems

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateMovement(e *ecs.ECS) {
	donburi.NewQuery(filter.Contains(components.Position, components.Velocity)).Each(e.World, func(entry *donburi.Entry) {
		velocityData := components.Velocity.Get(entry)
		positionData := components.Position.Get(entry)

		positionData.X += velocityData.X
		positionData.Y += velocityData.Y
	})
}
