package entities

import (
	"github.com/samber/do"
	"github.com/setanarut/kamera/v2"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	display2 "github.com/ubootgame/ubootgame/framework/graphics/display"
	"github.com/ubootgame/ubootgame/internal/components"
	"github.com/yohamta/donburi"
	"gonum.org/v1/gonum/spatial/r2"
)

var CameraTag = donburi.NewTag().SetName("Camera")

var Camera = ecsFramework.NewArchetype(
	CameraTag,
	components.Camera,
)

type NewCameraParams struct {
	MoveSpeed, ZoomSpeed, RotationSpeed float64
	MinZoom, MaxZoom                    float64
	Target, Delta                       r2.Vec
}

var CameraFactory ecsFramework.EntityFactory[NewCameraParams] = func(i *do.Injector, params NewCameraParams) *donburi.Entry {
	display := do.MustInvoke[display2.Display](i)
	ecs := do.MustInvoke[ecsFramework.Service](i)

	entry := Camera.Spawn(ecs.ECS())

	virtualResolution := display.VirtualResolution()

	components.Camera.SetValue(entry, components.CameraData{
		Camera:        kamera.NewCamera(0, 0, virtualResolution.X, virtualResolution.Y),
		MoveSpeed:     500,
		RotationSpeed: 2,
		ZoomSpeed:     10,
		MinZoom:       -100,
		MaxZoom:       100,
	})

	return entry
}
