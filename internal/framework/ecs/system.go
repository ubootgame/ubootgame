package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"github.com/yohamta/donburi/ecs"
)

type System struct {
	*Injector
}

func (system *System) Update(e *ecs.ECS) {
	if system.Injector != nil {
		system.Inject(e.World)
	}
}

func (system *System) Layers() []lo.Tuple2[ecs.LayerID, Renderer] {
	return []lo.Tuple2[ecs.LayerID, Renderer]{}
}

type Renderer func(e *ecs.ECS, screen *ebiten.Image)