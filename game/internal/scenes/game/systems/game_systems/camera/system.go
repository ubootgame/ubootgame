package camera

import (
	"github.com/ubootgame/ubootgame/framework"
	"github.com/ubootgame/ubootgame/framework/camera"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"go/types"
	"gonum.org/v1/gonum/spatial/r2"
)

const translationSpeed, zoomSpeed = 500.0, 0.1 // world unit
const rotationSpeed = 2                        // degrees
const minZoom, maxZoom = 0.5, 2.0

type System struct {
	ecsFramework.System

	settings framework.SettingsService[internal.Settings]
	camera   *camera.Camera
}

func NewCameraSystem(settings framework.SettingsService[internal.Settings], e *ecs.ECS, camera *camera.Camera) *System {
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
	system.camera.Position.X -= translationSpeed / float64(system.settings.Settings().Internals.TPS)
}

func (system *System) PanRight(_ donburi.World, _ types.Nil) {
	system.camera.Position.X += translationSpeed / float64(system.settings.Settings().Internals.TPS)
}

func (system *System) PanUp(_ donburi.World, _ types.Nil) {
	system.camera.Position.Y -= translationSpeed / float64(system.settings.Settings().Internals.TPS)
}

func (system *System) PanDown(_ donburi.World, _ types.Nil) {
	system.camera.Position.Y += translationSpeed / float64(system.settings.Settings().Internals.TPS)
}

func (system *System) ZoomIn(_ donburi.World, _ types.Nil) {
	scale := min(maxZoom, system.camera.Zoom.X+zoomSpeed)
	system.camera.Zoom = r2.Vec{X: scale, Y: scale}
}

func (system *System) ZoomOut(_ donburi.World, _ types.Nil) {
	scale := max(minZoom, system.camera.Zoom.X-zoomSpeed)
	system.camera.Zoom = r2.Vec{X: scale, Y: scale}
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
