package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"github.com/yohamta/donburi/ecs"
)

type Renderer ecs.RendererWithArg[ebiten.Image]

type System interface {
	Update(e *ecs.ECS)
	Layers() []lo.Tuple2[ecs.LayerID, Renderer]
}

func RegisterSystem(e *ecs.ECS, system System) {
	e.AddSystem(system.Update)
	for _, layer := range system.Layers() {
		e.AddRenderer(layer.A, layer.B)
	}
}
