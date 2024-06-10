package scenes

import (
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/pkg"
	ecsFramework "github.com/ubootgame/ubootgame/pkg/ecs"
	"github.com/ubootgame/ubootgame/pkg/services/resources"
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

func NewECSScene[S any](settings pkg.SettingsService[S], resources pkg.ResourceService, resourceLibrary *resources.Library) *ECSScene[S] {
	return &ECSScene[S]{
		BaseScene: NewBaseScene(settings, resources, resourceLibrary),
		ECS:       ecs.NewECS(donburi.NewWorld()),
	}
}

func (scene *ECSScene[S]) RegisterSystem(system System) {
	scene.ECS.AddSystem(system.Update)
	for _, layer := range system.Layers() {
		scene.ECS.AddRenderer(layer.A, layer.B)
	}
}
