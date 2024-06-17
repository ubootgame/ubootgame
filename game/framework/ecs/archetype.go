package ecs

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

func (a *Archetype) Spawn(e *ECS, cs ...donburi.IComponentType) *donburi.Entry {
	entry := e.World.Entry(e.World.Create(
		append(a.components, cs...)...,
	))
	return entry
}

func (a *Archetype) SpawnOnLayer(e *ECS, layer ecs.LayerID, cs ...donburi.IComponentType) *donburi.Entry {
	entry := e.World.Entry(e.Create(
		layer,
		append(a.components, cs...)...,
	))
	return entry
}
