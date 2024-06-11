package ecs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"github.com/yohamta/donburi/ecs"
)

type Renderer ecs.RendererWithArg[ebiten.Image]

type System struct{}

func (system *System) Update(e *ecs.ECS) {
}

func (system *System) Layers() []lo.Tuple2[ecs.LayerID, Renderer] {
	return []lo.Tuple2[ecs.LayerID, Renderer]{}
}
