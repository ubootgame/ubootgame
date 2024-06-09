package game_system

import (
	"github.com/hajimehoshi/ebiten/v2"
	ecs2 "github.com/ubootgame/ubootgame/internal/framework/ecs"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/actors/player"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/game_system/camera"
	"github.com/yohamta/donburi/ecs"
	"go/types"
)

type InputSystem struct {
	ecs2.System
}

func NewInputSystem() *InputSystem {
	return &InputSystem{}
}

func (system *InputSystem) Update(e *ecs.ECS) {
	//Camera
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

	//Player
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		player.MoveLeftEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		player.MoveRightEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		player.ShootEvent.Publish(e.World, types.Nil{})
	}
}
