package game_systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/samber/do"
	"github.com/samber/lo"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/framework/graphics/display"
	"github.com/ubootgame/ubootgame/framework/input"
	"github.com/ubootgame/ubootgame/internal/components"
	"github.com/ubootgame/ubootgame/internal/systems/actors/player"
	"github.com/ubootgame/ubootgame/internal/systems/game_systems/camera"
	"github.com/ubootgame/ubootgame/internal/systems/game_systems/debug"
	"github.com/yohamta/donburi/ecs"
	"go/types"
	"gonum.org/v1/gonum/spatial/r2"
)

type inputSystem struct {
	display display.Display
	input   input.Input

	camera *components.CameraData
}

func NewInputSystem(i *do.Injector) ecsFramework.System {
	return &inputSystem{
		display: do.MustInvoke[display.Display](i),
		input:   do.MustInvoke[input.Input](i),
	}
}

func (system *inputSystem) Layers() []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer]{}
}

func (system *inputSystem) Update(e *ecs.ECS) {
	if entry, found := components.Camera.First(e.World); found {
		system.camera = components.Camera.Get(entry)
	}

	cursor := system.input.Cursor()

	screenX, screenY := ebiten.CursorPosition()

	screenPosition := r2.Vec{X: float64(screenX), Y: float64(screenY)}
	cursor.ScreenPosition = screenPosition

	worldX, worldY := system.camera.Camera.ScreenToWorld(screenX, screenY)
	cursor.WorldPosition = r2.Vec{X: worldX / system.display.VirtualResolution().X, Y: worldY / system.display.VirtualResolution().X}

	// Camera
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		camera.PanLeftEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		camera.PanRightEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		camera.PanUpEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		camera.PanDownEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		camera.ZoomInEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		camera.ZoomOutEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		camera.RotateLeftEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		camera.RotateRightEvent.Publish(e.World, types.Nil{})
	}

	// Player
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		player.MoveLeftEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		player.MoveRightEvent.Publish(e.World, types.Nil{})
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		player.ShootEvent.Publish(e.World, types.Nil{})
	}

	// Debug
	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		debug.ToggleDebugEvent.Publish(e.World, types.Nil{})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		debug.ToggleDrawGrid.Publish(e.World, types.Nil{})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		debug.ToggleDrawCollisions.Publish(e.World, types.Nil{})
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		debug.ToggleDrawPositions.Publish(e.World, types.Nil{})
	}
}
