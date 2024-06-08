package game_system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/framework/ecs/injector"
	"github.com/ubootgame/ubootgame/internal/framework/ecs/systems"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/yohamta/donburi/ecs"
	"gonum.org/v1/gonum/spatial/r2"
)

type CursorSystem struct {
	systems.BaseSystem

	cursor *game_system.CursorData
	camera *game_system.CameraData
}

func NewCursorSystem() *CursorSystem {
	system := &CursorSystem{}
	system.Injector = injector.NewInjector([]injector.Injection{
		injector.Once([]injector.Injection{
			injector.Component(&system.cursor, game_system.Cursor),
			injector.Component(&system.camera, game_system.Camera),
		}),
	})
	return system
}

func (system *CursorSystem) Update(e *ecs.ECS) {
	system.BaseSystem.Update(e)

	screenX, screenY := ebiten.CursorPosition()
	screenPosition := r2.Vec{X: float64(screenX), Y: float64(screenY)}
	system.cursor.ScreenPosition = screenPosition
	system.cursor.WorldPosition = system.camera.ScreenToWorldPosition(screenPosition)
}
