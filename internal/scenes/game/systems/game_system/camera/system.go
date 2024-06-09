package camera

import (
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/framework"
	ecs2 "github.com/ubootgame/ubootgame/internal/framework/ecs"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"go/types"
	"gonum.org/v1/gonum/spatial/r2"
)

const translationSpeed, zoomSpeed = 500.0, 0.1 // world unit
const rotationSpeed = 2                        // degrees
const minZoom, maxZoom = 0.5, 2.0

type System struct {
	ecs2.System

	settings *internal.Settings
	camera   *framework.Camera
}

func NewCameraSystem(e *ecs.ECS, settings *internal.Settings, camera *framework.Camera) *System {
	system := &System{settings: settings, camera: camera}

	PanLeftEvent.Subscribe(e.World, system.PanLeft)
	PanRightEvent.Subscribe(e.World, system.PanRight)
	PanUpEvent.Subscribe(e.World, system.PanUp)
	PanDownEvent.Subscribe(e.World, system.PanDown)
	ZoomInEvent.Subscribe(e.World, system.ZoomIn)
	ZoomOutEvent.Subscribe(e.World, system.ZoomOut)
	RotateLeftEvent.Subscribe(e.World, system.RotateLeft)
	RotateRightEvent.Subscribe(e.World, system.RotateRight)

	return system
}

func (system *System) Update(e *ecs.ECS) {
	system.System.Update(e)

	system.camera.UpdateCameraMatrix()
}

func (system *System) PanLeft(_ donburi.World, _ types.Nil) {
	system.camera.Position.X -= translationSpeed / float64(system.settings.TargetTPS)
}

func (system *System) PanRight(_ donburi.World, _ types.Nil) {
	system.camera.Position.X += translationSpeed / float64(system.settings.TargetTPS)
}

func (system *System) PanUp(_ donburi.World, _ types.Nil) {
	system.camera.Position.Y -= translationSpeed / float64(system.settings.TargetTPS)
}

func (system *System) PanDown(_ donburi.World, _ types.Nil) {
	system.camera.Position.Y += translationSpeed / float64(system.settings.TargetTPS)
}

func (system *System) ZoomIn(_ donburi.World, _ types.Nil) {
	scale := min(maxZoom, system.camera.Scale.X+zoomSpeed)
	system.camera.Scale = r2.Vec{X: scale, Y: scale}
}

func (system *System) ZoomOut(_ donburi.World, _ types.Nil) {
	scale := max(minZoom, system.camera.Scale.X-zoomSpeed)
	system.camera.Scale = r2.Vec{X: scale, Y: scale}
}

func (system *System) RotateLeft(_ donburi.World, _ types.Nil) {
	newCameraRotation := system.camera.Rotation - rotationSpeed
	if newCameraRotation < 0 {
		newCameraRotation = 360 - newCameraRotation
	}
	system.camera.Rotation = newCameraRotation
}

func (system *System) RotateRight(_ donburi.World, _ types.Nil) {
	newCameraRotation := system.camera.Rotation + rotationSpeed
	if newCameraRotation >= 360 {
		newCameraRotation = newCameraRotation - 360
	}
	system.camera.Rotation = newCameraRotation
}
