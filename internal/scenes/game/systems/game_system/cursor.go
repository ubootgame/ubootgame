package game_system

import (
	"github.com/hajimehoshi/ebiten/v2"
	ecs2 "github.com/ubootgame/ubootgame/internal/framework/ecs"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/yohamta/donburi/ecs"
	"gonum.org/v1/gonum/spatial/r2"
)

type CursorSystem struct {
	ecs2.System

	cursor *game_system.CursorData
	camera *game_system.CameraData
}

func NewCursorSystem() *CursorSystem {
	system := &CursorSystem{}
	system.Injector = ecs2.NewInjector([]ecs2.Injection{
		ecs2.Once([]ecs2.Injection{
			ecs2.Component(&system.cursor, game_system.Cursor),
			ecs2.Component(&system.camera, game_system.Camera),
		}),
	})
	return system
}

func (system *CursorSystem) Update(e *ecs.ECS) {
	system.System.Update(e)

	screenX, screenY := ebiten.CursorPosition()
	screenPosition := r2.Vec{X: float64(screenX), Y: float64(screenY)}
	system.cursor.ScreenPosition = screenPosition
	system.cursor.WorldPosition = system.camera.ScreenToWorldPosition(screenPosition)
}
