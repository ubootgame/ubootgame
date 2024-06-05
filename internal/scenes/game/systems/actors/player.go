package actors

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/actors"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/weapons"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type playerSystem struct {
	playerEntry, cursorEntry *donburi.Entry
	fireTick                 uint64
}

var Player = &playerSystem{}

func (system *playerSystem) Update(e *ecs.ECS) {
	var ok bool
	if system.playerEntry == nil {
		if system.playerEntry, ok = actors.PlayerTag.First(e.World); !ok {
			panic("no player found")
		}
	}
	if system.cursorEntry == nil {
		if system.cursorEntry, ok = game_system.Cursor.First(e.World); !ok {
			panic("no cursor found")
		}
	}

	velocity := geometry.Velocity.Get(system.playerEntry)
	transform := geometry.Transform.Get(system.playerEntry)
	cursor := game_system.Cursor.Get(system.cursorEntry)

	acceleration := 0.01
	friction := 0.05
	maxSpeed := 0.25

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

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if system.fireTick%uint64(config.C.TargetTPS/8) == 0 {
			weapons.CreateBullet(e, transform.Center, cursor.WorldPosition)
		}
		system.fireTick++
	} else {
		system.fireTick = 0
	}
}
