package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/samber/lo"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/internal/components/physics"
	"github.com/yohamta/donburi/ecs"
)

type PhysicsSystem struct {
	space *cp.Space
}

func NewPhysicsSystem() *PhysicsSystem {
	return &PhysicsSystem{}
}

func (system *PhysicsSystem) Layers() []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer]{}
}

func (system *PhysicsSystem) Update(e *ecs.ECS) {
	if entry, found := physics.Space.First(e.World); found {
		system.space = physics.Space.Get(entry)
	}

	system.space.Step(1.0 / float64(ebiten.TPS()))
}
