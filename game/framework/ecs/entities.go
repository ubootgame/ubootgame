package ecs

import (
	"github.com/samber/do"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type EntityFactory[P any] func(i *do.Injector, e *ecs.ECS, parameters P) *donburi.Entry

func Spawn[P any](i *do.Injector, e *ecs.ECS, factory EntityFactory[P], parameters P) *donburi.Entry {
	return factory(i, e, parameters)
}
