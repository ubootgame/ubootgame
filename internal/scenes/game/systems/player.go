package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type playerSystem struct {
	entry *donburi.Entry
}

var Player = &playerSystem{}

func (system *playerSystem) Update(e *ecs.ECS) {
	if system.entry == nil {
		var ok bool
		if system.entry, ok = entities.PlayerTag.First(e.World); !ok {
			panic("no player found")
		}
	}

	velocity := components.Velocity.Get(system.entry)

	acceleration := 0.0001
	friction := 0.05
	maxSpeed := 0.005

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if velocity.X > 0 {
			velocity.X *= 1 - friction
		}
		velocity.X -= acceleration
		velocity.X = max(velocity.X, -maxSpeed)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if velocity.X < 0 {
			velocity.X *= 1 - friction
		}
		velocity.X += acceleration
		velocity.X = min(velocity.X, maxSpeed)
	} else {
		velocity.X *= 1 - friction
	}
}
