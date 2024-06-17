package debug

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/samber/do"
	"github.com/samber/lo"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/framework/graphics/d2d"
	"github.com/ubootgame/ubootgame/framework/input"
	"github.com/ubootgame/ubootgame/framework/settings"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/components"
	"github.com/ubootgame/ubootgame/internal/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"go/types"
	"runtime"
	"strings"
)

type debugSystem struct {
	settingsProvider settings.Provider[internal.Settings]
	input            input.Input

	camera *components.CameraData

	keys     []ebiten.Key
	memStats *runtime.MemStats
	ticks    uint64

	debugText        string
	debugTextBuilder *strings.Builder
	debugTextOptions *text.DrawOptions
}

func NewDebugSystem(i *do.Injector) ecsFramework.System {
	system := &debugSystem{
		settingsProvider: do.MustInvoke[settings.Provider[internal.Settings]](i),
		input:            do.MustInvoke[input.Input](i),
		keys:             make([]ebiten.Key, 0),
		memStats:         &runtime.MemStats{},
		debugTextBuilder: &strings.Builder{},
		debugTextOptions: &text.DrawOptions{
			DrawImageOptions: ebiten.DrawImageOptions{
				Filter: ebiten.FilterLinear,
			},
		},
	}

	e := do.MustInvoke[*ecsFramework.ECS](i)

	ToggleDebugEvent.Subscribe(e.World, system.ToggleDebug)
	ToggleDrawGrid.Subscribe(e.World, system.ToggleDrawGrid)
	ToggleDrawCollisions.Subscribe(e.World, system.ToggleDrawCollisions)
	ToggleDrawPositions.Subscribe(e.World, system.ToggleDrawPositions)

	return system
}

func (system *debugSystem) Layers() []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer]{
		{A: layers.Debug, B: system.DrawDebug},
	}
}

func (system *debugSystem) Update(e *ecs.ECS) {
	if !system.settingsProvider.Settings().Debug.Enabled {
		return
	}

	if entry, found := components.Camera.First(e.World); found {
		system.camera = components.Camera.Get(entry)
	}

	system.keys = inpututil.AppendPressedKeys(system.keys[:0])

	if system.ticks%uint64(system.settingsProvider.Settings().Internals.TPS*2) == 0 {
		runtime.ReadMemStats(system.memStats)
	}

	system.debugText = system.generateDebugText()

	system.ticks++
}

func (system *debugSystem) DrawDebug(_ *ecs.ECS, screen *ebiten.Image) {
	metrics := system.settingsProvider.Settings().Debug.FontFace.Metrics()
	system.debugTextOptions.LineSpacing = metrics.HAscent + metrics.HDescent + metrics.HLineGap

	text.Draw(screen, system.debugText, system.settingsProvider.Settings().Debug.FontFace, system.debugTextOptions)
}

func (system *debugSystem) ToggleDebug(_ donburi.World, _ types.Nil) {
	system.settingsProvider.Settings().Debug.Enabled = !system.settingsProvider.Settings().Debug.Enabled
}

func (system *debugSystem) ToggleDrawGrid(_ donburi.World, _ types.Nil) {
	system.settingsProvider.Settings().Debug.DrawGrid = !system.settingsProvider.Settings().Debug.DrawGrid
}

func (system *debugSystem) ToggleDrawCollisions(_ donburi.World, _ types.Nil) {
	system.settingsProvider.Settings().Debug.DrawCollisions = !system.settingsProvider.Settings().Debug.DrawCollisions
}

func (system *debugSystem) ToggleDrawPositions(_ donburi.World, _ types.Nil) {
	system.settingsProvider.Settings().Debug.DrawPositions = !system.settingsProvider.Settings().Debug.DrawPositions
}

func (system *debugSystem) generateDebugText() string {
	system.debugTextBuilder.Reset()

	ms := system.memStats

	cursor := system.input.Cursor()
	cursorWorldPosition := cursor.WorldPosition
	cursorScreenPosition := cursor.ScreenPosition

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
		system.settingsProvider.Settings().Debug.DrawGrid,
		system.settingsProvider.Settings().Debug.DrawCollisions,
		system.settingsProvider.Settings().Debug.DrawPositions,
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
