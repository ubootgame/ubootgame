package camera

import (
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/framework"
	ecsFramework "github.com/ubootgame/ubootgame/framework/game/ecs"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"go/types"
	"math"
)

type System struct {
	settings framework.SettingsService[internal.Settings]
	display  framework.DisplayService
	camera   *components.CameraData
}

func NewSystem(settings framework.SettingsService[internal.Settings], e *ecs.ECS, display framework.DisplayService) *System {
	system := &System{settings: settings, display: display}

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

func (system *System) Layers() []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer]{}
}

func (system *System) Update(e *ecs.ECS) {
	if entry, found := components.Camera.First(e.World); found {
		system.camera = components.Camera.Get(entry)
	}

	virtualResolution := system.display.VirtualResolution()
	if system.camera.Camera.Width() != virtualResolution.X || system.camera.Camera.Height() != virtualResolution.Y {
		system.camera.Camera.SetSize(virtualResolution.X, virtualResolution.Y)
	}

	if system.camera.Delta.X != 0 && system.camera.Delta.Y != 0 {
		factor := system.camera.MoveSpeed / float64(system.settings.Settings().Internals.TPS) / math.Sqrt(system.camera.Delta.X*system.camera.Delta.X+system.camera.Delta.Y*system.camera.Delta.Y)
		system.camera.Delta.X *= factor
		system.camera.Delta.Y *= factor
	}

	system.camera.Target.X += system.camera.Delta.X
	system.camera.Target.Y += system.camera.Delta.Y

	system.camera.Camera.LookAt(system.camera.Target.X, system.camera.Target.Y)

	system.camera.Delta.X = 0
	system.camera.Delta.Y = 0
}

func (system *System) PanLeft(_ donburi.World, _ types.Nil) {
	system.camera.Delta.X = -(system.camera.MoveSpeed / float64(system.settings.Settings().Internals.TPS))
}

func (system *System) PanRight(_ donburi.World, _ types.Nil) {
	system.camera.Delta.X = system.camera.MoveSpeed / float64(system.settings.Settings().Internals.TPS)
}

func (system *System) PanUp(_ donburi.World, _ types.Nil) {
	system.camera.Delta.Y = -(system.camera.MoveSpeed / float64(system.settings.Settings().Internals.TPS))
}

func (system *System) PanDown(_ donburi.World, _ types.Nil) {
	system.camera.Delta.Y = system.camera.MoveSpeed / float64(system.settings.Settings().Internals.TPS)
}

func (system *System) ZoomIn(_ donburi.World, _ types.Nil) {
	system.camera.Camera.ZoomFactor = min(system.camera.MaxZoom, system.camera.Camera.ZoomFactor+system.camera.ZoomSpeed)
}

func (system *System) ZoomOut(_ donburi.World, _ types.Nil) {
	system.camera.Camera.ZoomFactor = max(system.camera.MinZoom, system.camera.Camera.ZoomFactor-system.camera.ZoomSpeed)
}

func (system *System) RotateLeft(_ donburi.World, _ types.Nil) {
	system.camera.Camera.SetRotation(system.camera.Camera.Rotation() - system.camera.RotationSpeed)
}

func (system *System) RotateRight(_ donburi.World, _ types.Nil) {
	system.camera.Camera.SetRotation(system.camera.Camera.Rotation() + system.camera.RotationSpeed)
}
