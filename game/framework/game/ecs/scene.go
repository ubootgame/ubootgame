package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/framework"
	"github.com/ubootgame/ubootgame/framework/game"
	"github.com/ubootgame/ubootgame/framework/services/resources"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type Renderer ecs.RendererWithArg[ebiten.Image]

type System interface {
	Update(e *ecs.ECS)
	Layers() []lo.Tuple2[ecs.LayerID, Renderer]
}

type Scene[S any] struct {
	*game.BaseScene[S]
	ECS *ecs.ECS
}

func NewECSScene[S any](settings framework.SettingsService[S], resources framework.ResourceService, resourceLibrary *resources.Library) *Scene[S] {
	return &Scene[S]{
		BaseScene: game.NewBaseScene(settings, resources, resourceLibrary),
		ECS:       ecs.NewECS(donburi.NewWorld()),
	}
}

func (scene *Scene[S]) RegisterSystem(system System) {
	scene.ECS.AddSystem(system.Update)
	for _, layer := range system.Layers() {
		scene.ECS.AddRenderer(layer.A, layer.B)
	}
}
