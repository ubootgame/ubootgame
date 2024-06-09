package scenes

import (
	"github.com/samber/lo"
	ecsFramework "github.com/ubootgame/ubootgame/pkg/ecs"
	"github.com/ubootgame/ubootgame/pkg/resources"
	framework "github.com/ubootgame/ubootgame/pkg/settings"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type System interface {
	Update(e *ecs.ECS)
	Layers() []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer]
}

type ECSScene[S any] struct {
	*BaseScene[S]
	ECS *ecs.ECS
}

func NewECSScene[S any](settings *framework.Settings[S], resources *resources.Library) *ECSScene[S] {
	return &ECSScene[S]{
		BaseScene: NewBaseScene(settings, resources),
		ECS:       ecs.NewECS(donburi.NewWorld()),
	}
}

func (scene *ECSScene[S]) RegisterSystem(system System) {
	scene.ECS.AddSystem(system.Update)
	for _, layer := range system.Layers() {
		scene.ECS.AddRenderer(layer.A, layer.B)
	}
}
