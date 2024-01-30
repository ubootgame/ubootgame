package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

const cameraSpeed = 0.01

func UpdateCamera(e *ecs.ECS) {
	cameraEntry, _ := donburi.NewQuery(filter.Contains(components.Camera)).First(e.World)
	cameraData := components.Camera.Get(cameraEntry)

	utility.SetCameraMatrix(cameraData)

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		cameraData.Position.X -= cameraSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		cameraData.Position.X += cameraSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		cameraData.Position.Y -= cameraSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		cameraData.Position.Y += cameraSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		cameraData.ZoomFactor = max(0.5, cameraData.ZoomFactor-0.1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		cameraData.ZoomFactor = min(2.0, cameraData.ZoomFactor+0.1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		newCameraRotation := cameraData.Rotation - 1
		if newCameraRotation < 0 {
			newCameraRotation = 360 - newCameraRotation
		}
		cameraData.Rotation = newCameraRotation
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		newCameraRotation := cameraData.Rotation + 1
		if newCameraRotation >= 360 {
			newCameraRotation = newCameraRotation - 360
		}
		cameraData.Rotation = newCameraRotation
	}
}
