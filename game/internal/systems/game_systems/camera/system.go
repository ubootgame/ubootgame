package camera

import (
	"github.com/samber/do"
	"github.com/samber/lo"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/framework/graphics/display"
	"github.com/ubootgame/ubootgame/framework/settings"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"go/types"
	"math"
)

type cameraSystem struct {
	settingsProvider settings.Provider[internal.Settings]
	display          display.Display

	camera *components.CameraData
}

func NewCameraSystem(i *do.Injector) ecsFramework.System {
	system := &cameraSystem{
		settingsProvider: do.MustInvoke[settings.Provider[internal.Settings]](i),
		display:          do.MustInvoke[display.Display](i),
	}

	e := do.MustInvoke[*ecsFramework.ECS](i)

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

func (system *cameraSystem) Layers() []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer]{}
}

func (system *cameraSystem) Update(e *ecs.ECS) {
	if entry, found := components.Camera.First(e.World); found {
		system.camera = components.Camera.Get(entry)
	}

	virtualResolution := system.display.VirtualResolution()
	if system.camera.Camera.Width() != virtualResolution.X || system.camera.Camera.Height() != virtualResolution.Y {
		system.camera.Camera.SetSize(virtualResolution.X, virtualResolution.Y)
	}

	if system.camera.Delta.X != 0 && system.camera.Delta.Y != 0 {
		factor := system.camera.MoveSpeed / float64(system.settingsProvider.Settings().Internals.TPS) / math.Sqrt(system.camera.Delta.X*system.camera.Delta.X+system.camera.Delta.Y*system.camera.Delta.Y)
		system.camera.Delta.X *= factor
		system.camera.Delta.Y *= factor
	}

	system.camera.Target.X += system.camera.Delta.X
	system.camera.Target.Y += system.camera.Delta.Y

	system.camera.Camera.LookAt(system.camera.Target.X, system.camera.Target.Y)

	system.camera.Delta.X = 0
	system.camera.Delta.Y = 0
}

func (system *cameraSystem) PanLeft(_ donburi.World, _ types.Nil) {
	system.camera.Delta.X = -(system.camera.MoveSpeed / float64(system.settingsProvider.Settings().Internals.TPS))
}

func (system *cameraSystem) PanRight(_ donburi.World, _ types.Nil) {
	system.camera.Delta.X = system.camera.MoveSpeed / float64(system.settingsProvider.Settings().Internals.TPS)
}

func (system *cameraSystem) PanUp(_ donburi.World, _ types.Nil) {
	system.camera.Delta.Y = -(system.camera.MoveSpeed / float64(system.settingsProvider.Settings().Internals.TPS))
}

func (system *cameraSystem) PanDown(_ donburi.World, _ types.Nil) {
	system.camera.Delta.Y = system.camera.MoveSpeed / float64(system.settingsProvider.Settings().Internals.TPS)
}

func (system *cameraSystem) ZoomIn(_ donburi.World, _ types.Nil) {
	system.camera.Camera.ZoomFactor = min(system.camera.MaxZoom, system.camera.Camera.ZoomFactor+system.camera.ZoomSpeed)
}

func (system *cameraSystem) ZoomOut(_ donburi.World, _ types.Nil) {
	system.camera.Camera.ZoomFactor = max(system.camera.MinZoom, system.camera.Camera.ZoomFactor-system.camera.ZoomSpeed)
}

func (system *cameraSystem) RotateLeft(_ donburi.World, _ types.Nil) {
	system.camera.Camera.SetRotation(system.camera.Camera.Rotation() - system.camera.RotationSpeed)
}

func (system *cameraSystem) RotateRight(_ donburi.World, _ types.Nil) {
	system.camera.Camera.SetRotation(system.camera.Camera.Rotation() + system.camera.RotationSpeed)
}
