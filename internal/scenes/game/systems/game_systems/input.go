package game_systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/ubootgame/ubootgame/internal/framework"
	ecsFramework "github.com/ubootgame/ubootgame/internal/framework/ecs"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/actors/player"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/game_systems/camera"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/game_systems/debug"
	"github.com/yohamta/donburi/ecs"
	"go/types"
	"gonum.org/v1/gonum/spatial/r2"
)

type InputSystem struct {
	ecsFramework.System

	cursor *framework.Cursor
	camera *framework.Camera
}

func NewInputSystem(cursor *framework.Cursor, camera *framework.Camera) *InputSystem {
	return &InputSystem{cursor: cursor, camera: camera}
}

func (system *InputSystem) Update(e *ecs.ECS) {
	screenX, screenY := ebiten.CursorPosition()
	screenPosition := r2.Vec{X: float64(screenX), Y: float64(screenY)}
	system.cursor.ScreenPosition = screenPosition
	system.cursor.WorldPosition = system.camera.ScreenToWorldPosition(screenPosition)

	// Camera
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		camera.PanLeftEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		camera.PanRightEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		camera.PanUpEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		camera.PanDownEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		camera.ZoomInEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		camera.ZoomOutEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		camera.RotateLeftEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		camera.RotateRightEvent.Publish(e.World, types.Nil{})
	}

	// Player
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		player.MoveLeftEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		player.MoveRightEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		player.ShootEvent.Publish(e.World, types.Nil{})
	}

	// Debug
	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		debug.ToggleDebugEvent.Publish(e.World, types.Nil{})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		debug.ToggleDrawGrid.Publish(e.World, types.Nil{})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		debug.ToggleDrawCollisions.Publish(e.World, types.Nil{})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		debug.ToggleDrawPositions.Publish(e.World, types.Nil{})
	}
}
