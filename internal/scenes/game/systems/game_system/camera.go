package game_system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/injector"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/systems"
	"github.com/yohamta/donburi/ecs"
)

const translationSpeed, zoomSpeed = 0.5, 0.1 // world unit
const rotationSpeed = 2                      // degrees
const minZoom, maxZoom = 0.5, 2.0

type CameraSystem struct {
	systems.BaseSystem

	camera  *game_system.CameraData
	display *game_system.DisplayData
}

func NewCameraSystem() *CameraSystem {
	system := &CameraSystem{}
	system.Injector = injector.NewInjector([]injector.Injection{
		injector.Once([]injector.Injection{
			injector.Component(&system.display, game_system.Display),
			injector.Component(&system.camera, game_system.Camera),
		}),
	})
	return system
}

func (system *CameraSystem) Update(e *ecs.ECS) {
	system.BaseSystem.Update(e)

	utility.UpdateCameraMatrix(system.display, system.camera)

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		system.camera.Position.X -= translationSpeed / float64(config.C.TargetTPS)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		system.camera.Position.X += translationSpeed / float64(config.C.TargetTPS)
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		system.camera.Position.Y -= translationSpeed / float64(config.C.TargetTPS)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		system.camera.Position.Y += translationSpeed / float64(config.C.TargetTPS)
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		system.camera.ZoomFactor = max(minZoom, system.camera.ZoomFactor-zoomSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		system.camera.ZoomFactor = min(maxZoom, system.camera.ZoomFactor+zoomSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		newCameraRotation := system.camera.Rotation - rotationSpeed
		if newCameraRotation < 0 {
			newCameraRotation = 360 - newCameraRotation
		}
		system.camera.Rotation = newCameraRotation
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		newCameraRotation := system.camera.Rotation + rotationSpeed
		if newCameraRotation >= 360 {
			newCameraRotation = newCameraRotation - 360
		}
		system.camera.Rotation = newCameraRotation
	}
}
