package game_system

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/framework"
	"github.com/ubootgame/ubootgame/internal/framework/draw"
	ecs2 "github.com/ubootgame/ubootgame/internal/framework/ecs"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/yohamta/donburi/ecs"
	"runtime"
	"strings"
)

type DebugSystem struct {
	ecs2.System

	settings *internal.Settings
	cursor   *framework.Cursor
	camera   *framework.Camera

	keys             []ebiten.Key
	memStats         *runtime.MemStats
	ticks            uint64
	debugTextBuilder *strings.Builder
	debugTextOptions *text.DrawOptions
}

func NewDebugSystem(settings *internal.Settings, cursor *framework.Cursor, camera *framework.Camera) *DebugSystem {
	return &DebugSystem{
		settings:         settings,
		cursor:           cursor,
		camera:           camera,
		keys:             make([]ebiten.Key, 0),
		memStats:         &runtime.MemStats{},
		debugTextBuilder: &strings.Builder{},
		debugTextOptions: &text.DrawOptions{
			DrawImageOptions: ebiten.DrawImageOptions{
				Filter: ebiten.FilterLinear,
			},
		},
	}
}

func (system *DebugSystem) Layers() []lo.Tuple2[ecs.LayerID, ecs2.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecs2.Renderer]{
		{A: layers.Debug, B: system.DrawDebug},
	}
}

func (system *DebugSystem) Update(e *ecs.ECS) {
	system.System.Update(e)

	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		system.settings.Debug.Enabled = !system.settings.Debug.Enabled
	}

	if !system.settings.Debug.Enabled {
		return
	}

	system.keys = inpututil.AppendPressedKeys(system.keys[:0])

	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		system.settings.Debug.DrawGrid = !system.settings.Debug.DrawGrid
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		system.settings.Debug.DrawCollisions = !system.settings.Debug.DrawCollisions
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		system.settings.Debug.DrawPositions = !system.settings.Debug.DrawPositions
	}

	if system.ticks%uint64(system.settings.TargetTPS*2) == 0 {
		runtime.ReadMemStats(system.memStats)
	}
	system.ticks++

}

func (system *DebugSystem) DrawDebug(_ *ecs.ECS, screen *ebiten.Image) {
	debugText := system.generateDebugText()

	metrics := system.settings.Debug.FontFace.Metrics()
	system.debugTextOptions.LineSpacing = metrics.HAscent + metrics.HDescent + metrics.HLineGap

	text.Draw(screen, debugText, system.settings.Debug.FontFace, system.debugTextOptions)
}

func (system *DebugSystem) generateDebugText() string {
	system.debugTextBuilder.Reset()

	ms := system.memStats

	cursorWorldPosition := system.cursor.WorldPosition
	cursorScreenPosition := system.cursor.ScreenPosition

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
		system.settings.Debug.DrawGrid,
		system.settings.Debug.DrawCollisions,
		system.settings.Debug.DrawPositions,
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
		ebiten.IsVsyncEnabled(),
		system.settings.Display.ScalingFactor,
		strings.Join(lo.Map(system.keys, func(item ebiten.Key, index int) string {
			return item.String()
		}), ", "),
		cursorScreenPosition.X, cursorScreenPosition.Y,
		cursorWorldPosition.X, cursorWorldPosition.Y,
		system.camera.Position.X, system.camera.Position.Y,
		system.camera.Scale,
		system.camera.Rotation,
		draw.FormatBytes(ms.Alloc), draw.FormatBytes(ms.TotalAlloc), draw.FormatBytes(ms.Sys),
		draw.FormatBytes(ms.NextGC), ms.NumGC)

	return system.debugTextBuilder.String()
}
