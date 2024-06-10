package game_systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/ubootgame/ubootgame/framework/camera"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/framework/input"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/actors/player"
	gameSystem "github.com/ubootgame/ubootgame/internal/scenes/game/systems/game_systems/camera"
	debugSystem "github.com/ubootgame/ubootgame/internal/scenes/game/systems/game_systems/debug"
	"github.com/yohamta/donburi/ecs"
	"go/types"
	"gonum.org/v1/gonum/spatial/r2"
)

type InputSystem struct {
	ecsFramework.System

	cursor *input.Cursor
	camera *camera.Camera
}

func NewInputSystem(cursor *input.Cursor, camera *camera.Camera) *InputSystem {
	return &InputSystem{cursor: cursor, camera: camera}
}

func (system *InputSystem) Update(e *ecs.ECS) {
	screenX, screenY := ebiten.CursorPosition()
	screenPosition := r2.Vec{X: float64(screenX), Y: float64(screenY)}
	system.cursor.ScreenPosition = screenPosition
	system.cursor.WorldPosition = system.camera.ScreenToWorldPosition(screenPosition)

	// Camera
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		gameSystem.PanLeftEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		gameSystem.PanRightEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		gameSystem.PanUpEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		gameSystem.PanDownEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		gameSystem.ZoomInEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		gameSystem.ZoomOutEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		gameSystem.RotateLeftEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		gameSystem.RotateRightEvent.Publish(e.World, types.Nil{})
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
		debugSystem.ToggleDebugEvent.Publish(e.World, types.Nil{})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		debugSystem.ToggleDrawGrid.Publish(e.World, types.Nil{})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		debugSystem.ToggleDrawCollisions.Publish(e.World, types.Nil{})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		debugSystem.ToggleDrawPositions.Publish(e.World, types.Nil{})
	}
}
