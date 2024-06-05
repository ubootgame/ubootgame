package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/injector"
	"github.com/yohamta/donburi/ecs"
)

type BaseSystem struct {
	*injector.Injector
}

func (system *BaseSystem) Update(e *ecs.ECS) {
	if system.Injector != nil {
		system.Inject(e.World)
	}
}

func (system *BaseSystem) Draw(_ *ecs.ECS, _ *ebiten.Image) {}
