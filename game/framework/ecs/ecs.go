package ecs

import (
	"github.com/samber/do"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type ECS struct {
	*ecs.ECS
}

func NewECS(_ *do.Injector) (*ECS, error) {
	return &ECS{ECS: ecs.NewECS(donburi.NewWorld())}, nil
}

func (e ECS) RegisterSystem(i *do.Injector, systemFn func(i *do.Injector) System) {
	system := systemFn(i)
	e.AddSystem(system.Update)
	for _, layer := range system.Layers() {
		e.AddRenderer(layer.A, layer.B)
	}
}
