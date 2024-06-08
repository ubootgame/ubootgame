package game_system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	gameSystemEntities "github.com/ubootgame/ubootgame/internal/scenes/game/entities/game_system"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/injector"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/systems"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

const translationSpeed, zoomSpeed = 500.0, 0.1 // world unit
const rotationSpeed = 2                        // degrees
const minZoom, maxZoom = 0.5, 2.0

type CameraSystem struct {
	systems.BaseSystem

	camera    *game_system.CameraData
	transform *transform.TransformData
	display   *game_system.DisplayData
}

func NewCameraSystem() *CameraSystem {
	system := &CameraSystem{}
	system.Injector = injector.NewInjector([]injector.Injection{
		injector.Once([]injector.Injection{
			injector.Component(&system.display, game_system.Display),
			injector.WithTag(gameSystemEntities.CameraTag, []injector.Injection{
				injector.Component(&system.camera, game_system.Camera),
				injector.Component(&system.transform, transform.Transform),
			}),
		}),
	})
	return system
}

func (system *CameraSystem) Update(e *ecs.ECS) {
	system.BaseSystem.Update(e)

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		system.transform.LocalPosition.X -= translationSpeed / float64(config.C.TargetTPS)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		system.transform.LocalPosition.X += translationSpeed / float64(config.C.TargetTPS)
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		system.transform.LocalPosition.Y -= translationSpeed / float64(config.C.TargetTPS)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		system.transform.LocalPosition.Y += translationSpeed / float64(config.C.TargetTPS)
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		scale := max(minZoom, system.transform.LocalScale.X-zoomSpeed)
		system.transform.LocalScale = math.NewVec2(scale, scale)
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		scale := min(maxZoom, system.transform.LocalScale.X+zoomSpeed)
		system.transform.LocalScale = math.NewVec2(scale, scale)
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		newCameraRotation := system.transform.LocalRotation - rotationSpeed
		if newCameraRotation < 0 {
			newCameraRotation = 360 - newCameraRotation
		}
		system.transform.LocalRotation = newCameraRotation
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		newCameraRotation := system.transform.LocalRotation + rotationSpeed
		if newCameraRotation >= 360 {
			newCameraRotation = newCameraRotation - 360
		}
		system.transform.LocalRotation = newCameraRotation
	}

	utility.UpdateCameraMatrix(system.display, system.camera, system.transform)
}
