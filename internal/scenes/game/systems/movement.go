package systems

import (
	dresolv "github.com/ubootgame/ubootgame/internal/resolv"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateMovement(e *ecs.ECS) {
	donburi.NewQuery(filter.Contains(components.Object, components.Velocity)).Each(e.World, func(entry *donburi.Entry) {
		velocityData := components.Velocity.Get(entry)
		object := dresolv.GetObject(entry)

		object.X += velocityData.X
	})
}
