package ecs

import (
	"github.com/samber/do"
	"github.com/yohamta/donburi"
)

type EntityFactory[P any] func(i *do.Injector, parameters P) *donburi.Entry

//func Spawn[P any](i *do.Injector, factory EntityFactory[P], parameters P) *donburi.Entry {
//	return factory(i, parameters)
//}

func (f EntityFactory[P]) Spawn(i *do.Injector, parameters P) *donburi.Entry {
	return f(i, parameters)
}
