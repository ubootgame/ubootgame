package systems

import (
	dresolv "github.com/ubootgame/ubootgame/internal/resolv"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateObjects(e *ecs.ECS) {
	donburi.NewQuery(filter.Contains(components.Position, components.Object)).Each(e.World, func(entry *donburi.Entry) {
		positionData := components.Position.Get(entry)
		object := dresolv.GetObject(entry)

		object.X = positionData.Center.X
		object.Y = positionData.Center.Y

		object.Update()
	})
}
