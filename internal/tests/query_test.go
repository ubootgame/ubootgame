package tests

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"testing"
)

var component = donburi.NewComponentType[bool]()

func BenchmarkCreateQueryEveryTime(b *testing.B) {
	world := donburi.NewWorld()
	e := ecs.NewECS(world)
	_ = e.World.Entry(e.World.Create(component))

	for i := 0; i < b.N; i++ {
		donburi.NewQuery(filter.Contains(component)).Each(e.World, func(entry *donburi.Entry) {})
	}
}

func BenchmarkCreateQueryOnce(b *testing.B) {
	world := donburi.NewWorld()
	e := ecs.NewECS(world)
	_ = e.World.Entry(e.World.Create(component))
	query := donburi.NewQuery(filter.Contains(component))

	for i := 0; i < b.N; i++ {
		query.Each(e.World, func(entry *donburi.Entry) {})
	}
}
