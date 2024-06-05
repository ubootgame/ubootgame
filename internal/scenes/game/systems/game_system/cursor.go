package game_system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"gonum.org/v1/gonum/spatial/r2"
)

type cursorSystem struct {
	cursorEntry, cameraEntry *donburi.Entry
}

var Cursor = &cursorSystem{}

func (system *cursorSystem) Update(e *ecs.ECS) {
	var ok bool
	if system.cursorEntry == nil {
		if system.cursorEntry, ok = game_system.Cursor.First(e.World); !ok {
			panic("no cursor found")
		}
	}
	if system.cameraEntry == nil {
		if system.cameraEntry, ok = game_system.Camera.First(e.World); !ok {
			panic("no camera found")
		}
	}

	cursor := game_system.Cursor.Get(system.cursorEntry)
	camera := game_system.Camera.Get(system.cameraEntry)

	screenX, screenY := ebiten.CursorPosition()
	screenPosition := r2.Vec{X: float64(screenX), Y: float64(screenY)}
	cursor.ScreenPosition = screenPosition
	cursor.WorldPosition = camera.ScreenToWorldPosition(screenPosition)
}
