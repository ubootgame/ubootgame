package debug

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/framework"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/framework/graphics/d2d"
	"github.com/ubootgame/ubootgame/framework/input"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/components"
	"github.com/ubootgame/ubootgame/internal/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"go/types"
	"runtime"
	"strings"
)

type System struct {
	settings framework.SettingsService[internal.Settings]

	cursor *input.Cursor

	camera *components.CameraData

	keys     []ebiten.Key
	memStats *runtime.MemStats
	ticks    uint64

	debugText        string
	debugTextBuilder *strings.Builder
	debugTextOptions *text.DrawOptions
}

func NewSystem(settings framework.SettingsService[internal.Settings], e *ecs.ECS, cursor *input.Cursor) *System {
	system := &System{
		settings:         settings,
		cursor:           cursor,
		keys:             make([]ebiten.Key, 0),
		memStats:         &runtime.MemStats{},
		debugTextBuilder: &strings.Builder{},
		debugTextOptions: &text.DrawOptions{
			DrawImageOptions: ebiten.DrawImageOptions{
				Filter: ebiten.FilterLinear,
			},
		},
	}

	ToggleDebugEvent.Subscribe(e.World, system.ToggleDebug)
	ToggleDrawGrid.Subscribe(e.World, system.ToggleDrawGrid)
	ToggleDrawCollisions.Subscribe(e.World, system.ToggleDrawCollisions)
	ToggleDrawPositions.Subscribe(e.World, system.ToggleDrawPositions)

	return system
}

func (system *System) Layers() []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer]{
		{A: layers.Debug, B: system.DrawDebug},
	}
}

func (system *System) Update(e *ecs.ECS) {
	if !system.settings.Settings().Debug.Enabled {
		return
	}

	if entry, found := components.Camera.First(e.World); found {
		system.camera = components.Camera.Get(entry)
	}

	system.keys = inpututil.AppendPressedKeys(system.keys[:0])

	if system.ticks%uint64(system.settings.Settings().Internals.TPS*2) == 0 {
		runtime.ReadMemStats(system.memStats)
	}

	system.debugText = system.generateDebugText()

	system.ticks++
}

func (system *System) DrawDebug(_ *ecs.ECS, screen *ebiten.Image) {
	metrics := system.settings.Settings().Debug.FontFace.Metrics()
	system.debugTextOptions.LineSpacing = metrics.HAscent + metrics.HDescent + metrics.HLineGap

	text.Draw(screen, system.debugText, system.settings.Settings().Debug.FontFace, system.debugTextOptions)
}

func (system *System) ToggleDebug(_ donburi.World, _ types.Nil) {
	system.settings.Settings().Debug.Enabled = !system.settings.Settings().Debug.Enabled
}

func (system *System) ToggleDrawGrid(_ donburi.World, _ types.Nil) {
	system.settings.Settings().Debug.DrawGrid = !system.settings.Settings().Debug.DrawGrid
}

func (system *System) ToggleDrawCollisions(_ donburi.World, _ types.Nil) {
	system.settings.Settings().Debug.DrawCollisions = !system.settings.Settings().Debug.DrawCollisions
}

func (system *System) ToggleDrawPositions(_ donburi.World, _ types.Nil) {
	system.settings.Settings().Debug.DrawPositions = !system.settings.Settings().Debug.DrawPositions
}

func (system *System) generateDebugText() string {
	system.debugTextBuilder.Reset()

	ms := system.memStats

	cursorWorldPosition := system.cursor.WorldPosition
	cursorScreenPosition := system.cursor.ScreenPosition

	cameraWorldPositionX, cameraWorldPositionY := system.camera.Camera.Target()

	_, _ = fmt.Fprintf(system.debugTextBuilder, `(/ to toggle debugSystem)
Draw grid (F1): %v
Draw collisions (F2): %v
Draw positions (F3): %v
FPS: %.1f
TPS: %.1f
VSync: %v
Device scale factor: %.1f
Keys: %v
Cursor screen position: %.0f, %.0f
Cursor world position: %.3f, %.3f
Camera position: %.3f, %.3f
Camera zoom: %.2f
Camera rotation: %.2f
Alloc: %s
Total: %s
Sys: %s
NextGC: %s
NumGC: %d`,
		system.settings.Settings().Debug.DrawGrid,
		system.settings.Settings().Debug.DrawCollisions,
		system.settings.Settings().Debug.DrawPositions,
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
		ebiten.IsVsyncEnabled(),
		ebiten.Monitor().DeviceScaleFactor(),
		strings.Join(lo.Map(system.keys, func(item ebiten.Key, index int) string {
			return item.String()
		}), ", "),
		cursorScreenPosition.X, cursorScreenPosition.Y,
		cursorWorldPosition.X, cursorWorldPosition.Y,
		cameraWorldPositionX, cameraWorldPositionY,
		system.camera.Camera.ZoomFactor,
		system.camera.Camera.Rotation(),
		d2d.FormatBytes(ms.Alloc), d2d.FormatBytes(ms.TotalAlloc), d2d.FormatBytes(ms.Sys),
		d2d.FormatBytes(ms.NextGC), ms.NumGC)

	return system.debugTextBuilder.String()
}
