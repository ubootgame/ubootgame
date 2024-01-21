package utility

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type Archetype struct {
	components []donburi.IComponentType
}

func NewArchetype(cs ...donburi.IComponentType) *Archetype {
	return &Archetype{
		components: cs,
	}
}

func (a *Archetype) Spawn(e *ecs.ECS, cs ...donburi.IComponentType) *donburi.Entry {
	entry := e.World.Entry(e.Create(
		ecs.LayerDefault,
		append(a.components, cs...)...,
	))
	return entry
}
