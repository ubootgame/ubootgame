package game_system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/framework"
	ecs2 "github.com/ubootgame/ubootgame/internal/framework/ecs"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	gameSystemEntities "github.com/ubootgame/ubootgame/internal/scenes/game/entities/game_system"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

const translationSpeed, zoomSpeed = 500.0, 0.1 // world unit
const rotationSpeed = 2                        // degrees
const minZoom, maxZoom = 0.5, 2.0

type CameraSystem struct {
	ecs2.System

	settings *internal.Settings

	camera    *game_system.CameraData
	transform *transform.TransformData
}

func NewCameraSystem(settings *internal.Settings) *CameraSystem {
	system := &CameraSystem{settings: settings}
	system.Injector = ecs2.NewInjector([]ecs2.Injection{
		ecs2.Once([]ecs2.Injection{
			ecs2.WithTag(gameSystemEntities.CameraTag, []ecs2.Injection{
				ecs2.Component(&system.camera, game_system.Camera),
				ecs2.Component(&system.transform, transform.Transform),
			}),
		}),
	})
	return system
}

func (system *CameraSystem) Update(e *ecs.ECS) {
	system.System.Update(e)

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		system.transform.LocalPosition.X -= translationSpeed / float64(system.settings.TargetTPS)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		system.transform.LocalPosition.X += translationSpeed / float64(system.settings.TargetTPS)
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		system.transform.LocalPosition.Y -= translationSpeed / float64(system.settings.TargetTPS)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		system.transform.LocalPosition.Y += translationSpeed / float64(system.settings.TargetTPS)
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

	framework.UpdateCameraMatrix(&system.settings.Display, system.camera, system.transform)
}
