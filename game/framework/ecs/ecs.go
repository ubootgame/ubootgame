package ecs

import (
	"github.com/samber/do"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type Service interface {
	ECS() *ecs.ECS
	World() donburi.World
	RegisterSystem(systemFn func(i *do.Injector) System)
}

type service struct {
	injector *do.Injector
	ecs      *ecs.ECS
}

func NewECSService(i *do.Injector) (Service, error) {
	return &service{
		injector: i,
		ecs:      ecs.NewECS(donburi.NewWorld()),
	}, nil
}

func (s *service) ECS() *ecs.ECS {
	return s.ecs
}

func (s *service) World() donburi.World {
	return s.ecs.World
}

func (s *service) RegisterSystem(systemFn func(i *do.Injector) System) {
	system := systemFn(s.injector)
	s.ecs.AddSystem(system.Update)
	for _, layer := range system.Layers() {
		s.ecs.AddRenderer(layer.A, layer.B)
	}
}
