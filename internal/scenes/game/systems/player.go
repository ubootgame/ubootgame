package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"gonum.org/v1/gonum/spatial/r2"
)

type playerSystem struct {
	playerEntry, cameraEntry *donburi.Entry
	fireTick                 uint64
}

var Player = &playerSystem{}

func (system *playerSystem) Update(e *ecs.ECS) {
	var ok bool
	if system.playerEntry == nil {
		if system.playerEntry, ok = entities.PlayerTag.First(e.World); !ok {
			panic("no player found")
		}
	}
	if system.cameraEntry == nil {
		if system.cameraEntry, ok = components.Camera.First(e.World); !ok {
			panic("no camera found")
		}
	}

	velocity := components.Velocity.Get(system.playerEntry)
	transform := components.Transform.Get(system.playerEntry)
	camera := components.Camera.Get(system.cameraEntry)

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

	cursorX, cursorY := ebiten.CursorPosition()
	cursorPos := camera.ScreenToWorldPosition(r2.Vec{X: float64(cursorX), Y: float64(cursorY)})

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if system.fireTick%uint64(config.C.TargetTPS/8) == 0 {
			entities.CreateBullet(e, transform.Center, cursorPos)
		}
		system.fireTick++
	} else {
		system.fireTick = 0
	}
}
