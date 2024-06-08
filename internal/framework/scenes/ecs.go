package scenes

import (
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/internal"
	ecsSystems "github.com/ubootgame/ubootgame/internal/framework/ecs"
	"github.com/ubootgame/ubootgame/internal/framework/resources"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type System interface {
	Update(e *ecs.ECS)
	Layers() []lo.Tuple2[ecs.LayerID, ecsSystems.Renderer]
}

type ECSScene struct {
	*BaseScene
	ECS *ecs.ECS
}

func NewECSScene(settings *internal.Settings, resources *resources.Library) *ECSScene {
	return &ECSScene{
		BaseScene: NewBaseScene(settings, resources),
		ECS:       ecs.NewECS(donburi.NewWorld()),
	}
}

func (scene *ECSScene) RegisterSystem(system System) {
	scene.ECS.AddSystem(system.Update)
	for _, layer := range system.Layers() {
		scene.ECS.AddRenderer(layer.A, layer.B)
	}
}
