package game_system

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/weapons"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/injector"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/systems"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"runtime"
	"strings"
)

type DebugSystem struct {
	systems.BaseSystem

	debug   *game_system.DebugData
	camera  *game_system.CameraData
	display *game_system.DisplayData
	cursor  *game_system.CursorData

	keys      []ebiten.Key
	memStats  *runtime.MemStats
	ticks     uint64
	debugText *strings.Builder
}

func NewDebugSystem() *DebugSystem {
	system := &DebugSystem{
		keys:      make([]ebiten.Key, 0),
		memStats:  &runtime.MemStats{},
		debugText: &strings.Builder{},
	}
	system.Injector = injector.NewInjector([]injector.Injection{
		injector.Once([]injector.Injection{
			injector.Component(&system.debug, game_system.Debug),
			injector.Component(&system.camera, game_system.Camera),
			injector.Component(&system.display, game_system.Display),
			injector.Component(&system.cursor, game_system.Cursor),
		}),
	})
	return system
}

func (system *DebugSystem) Update(e *ecs.ECS) {
	system.BaseSystem.Update(e)

	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		system.debug.Enabled = !system.debug.Enabled
	}

	if !system.debug.Enabled {
		return
	}

	system.keys = inpututil.AppendPressedKeys(system.keys[:0])

	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		system.debug.DrawGrid = !system.debug.DrawGrid
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		system.debug.DrawCollisions = !system.debug.DrawCollisions
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		system.debug.DrawPositions = !system.debug.DrawPositions
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
		ecs.NewQuery(layers.Foreground, filter.Contains(weapons.BulletTag)).Each(e.World, func(entry *donburi.Entry) {
			transform := geometry.Transform.Get(entry)
			velocity := geometry.Velocity.Get(entry)
			fmt.Printf("%v %v\n", transform, velocity)
		})
	}

	if system.ticks%uint64(config.C.TargetTPS*2) == 0 {
		runtime.ReadMemStats(system.memStats)
	}
	system.ticks++

	system.updateDebugText(system.debugText)
}

func (system *DebugSystem) Draw(_ *ecs.ECS, screen *ebiten.Image) {
	utility.Debug.PrintDebugTextAt(screen, system.debugText.String(), &ebiten.DrawImageOptions{})
}

func (system *DebugSystem) updateDebugText(builder *strings.Builder) {
	builder.Reset()

	ms := system.memStats

	worldPosition := system.cursor.WorldPosition
	screenPosition := system.cursor.ScreenPosition

	_, _ = fmt.Fprintf(builder, `(/ to toggle debugSystem)
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
		system.debug.DrawGrid,
		system.debug.DrawCollisions,
		system.debug.DrawPositions,
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
		ebiten.IsVsyncEnabled(),
		system.display.ScalingFactor,
		strings.Join(lo.Map(system.keys, func(item ebiten.Key, index int) string {
			return item.String()
		}), ", "),
		screenPosition.X, screenPosition.Y,
		worldPosition.X, worldPosition.Y,
		system.camera.Position.X, system.camera.Position.Y,
		system.camera.ZoomFactor,
		system.camera.Rotation,
		utility.FormatBytes(ms.Alloc), utility.FormatBytes(ms.TotalAlloc), utility.FormatBytes(ms.Sys),
		utility.FormatBytes(ms.NextGC), ms.NumGC)
}
