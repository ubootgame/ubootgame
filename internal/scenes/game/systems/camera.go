package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

const translationSpeed, zoomSpeed = 0.5, 0.1 // world unit
const rotationSpeed = 2                      // degrees
const minZoom, maxZoom = 0.5, 2.0

type cameraSystem struct {
	entry *donburi.Entry
}

var Camera = &cameraSystem{}

func (system *cameraSystem) Update(e *ecs.ECS) {
	if system.entry == nil {
		var ok bool
		if system.entry, ok = components.Camera.First(e.World); !ok {
			panic("no camera found")
		}
	}

	camera := components.Camera.Get(system.entry)

	utility.SetCameraMatrix(camera)

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		camera.Position.X -= translationSpeed / float64(config.C.TargetTPS)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		camera.Position.X += translationSpeed / float64(config.C.TargetTPS)
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		camera.Position.Y -= translationSpeed / float64(config.C.TargetTPS)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		camera.Position.Y += translationSpeed / float64(config.C.TargetTPS)
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		camera.ZoomFactor = max(minZoom, camera.ZoomFactor-zoomSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		camera.ZoomFactor = min(maxZoom, camera.ZoomFactor+zoomSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		newCameraRotation := camera.Rotation - rotationSpeed
		if newCameraRotation < 0 {
			newCameraRotation = 360 - newCameraRotation
		}
		camera.Rotation = newCameraRotation
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		newCameraRotation := camera.Rotation + rotationSpeed
		if newCameraRotation >= 360 {
			newCameraRotation = newCameraRotation - 360
		}
		camera.Rotation = newCameraRotation
	}
}
