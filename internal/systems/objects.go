package systems

import (
	"github.com/ubootgame/ubootgame/internal/components"
	dresolv "github.com/ubootgame/ubootgame/internal/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateObjects(ecs *ecs.ECS) {
	components.Object.Each(ecs.World, func(e *donburi.Entry) {
		obj := dresolv.GetObject(e)
		obj.Update()
	})
}
