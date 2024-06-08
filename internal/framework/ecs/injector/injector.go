package injector

import (
	"github.com/yohamta/donburi"
	"log"
)

type Injector struct {
	injections []Injection
}

func NewInjector(injections []Injection) *Injector {
	return &Injector{injections: injections}
}

func (i *Injector) Inject(world donburi.World) {
	for _, injection := range i.injections {
		injection.Run(world, nil)
	}
}

type Injection interface {
	Run(world donburi.World, entry *donburi.Entry)
}

type component[D any] struct {
	Injection
	target        **D
	componentType *donburi.ComponentType[D]
}

func Component[D any](target **D, componentType *donburi.ComponentType[D]) Injection {
	return &component[D]{
		target:        target,
		componentType: componentType,
	}
}

func (component *component[D]) Run(world donburi.World, entry *donburi.Entry) {
	if entry == nil {
		var ok bool
		entry, ok = component.componentType.First(world)
		if !ok {
			log.Panicf("no %v found", component.componentType)
		}
	}
	*component.target = component.componentType.Get(entry)
}

type once struct {
	Injection
	done       bool
	injections []Injection
}

func Once(injections []Injection) Injection {
	return &once{
		done:       false,
		injections: injections,
	}
}

func (once *once) Run(world donburi.World, entry *donburi.Entry) {
	if once.done {
		return
	}
	for _, injection := range once.injections {
		injection.Run(world, entry)
	}
	once.done = true
}

type withTag struct {
	Injection
	tagType    *donburi.ComponentType[donburi.Tag]
	injections []Injection
}

func WithTag(tagType *donburi.ComponentType[donburi.Tag], injections []Injection) Injection {
	return &withTag{
		tagType:    tagType,
		injections: injections,
	}
}

func (withTag *withTag) Run(world donburi.World, _ *donburi.Entry) {
	if entry, ok := withTag.tagType.First(world); !ok {
		log.Panicf("no %v found", withTag.tagType)
	} else {
		for _, injection := range withTag.injections {
			injection.Run(world, entry)
		}
	}
}
