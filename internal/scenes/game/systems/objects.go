package systems

import (
	dresolv "github.com/ubootgame/ubootgame/internal/resolv"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateObjects(ecs *ecs.ECS) {
	components.Object.Each(ecs.World, func(e *donburi.Entry) {
		obj := dresolv.GetObject(e)
		obj.Update()
	})
}
