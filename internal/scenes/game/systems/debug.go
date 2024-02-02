package systems

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"gonum.org/v1/gonum/spatial/r2"
	"runtime"
	"strings"
)

type debugSystem struct {
	debugEntry, cameraEntry, displayEntry *donburi.Entry
	keys                                  []ebiten.Key
	memStats                              *runtime.MemStats
	ticks                                 uint64
	debugText                             *strings.Builder
}

var Debug = &debugSystem{
	keys:      make([]ebiten.Key, 0),
	memStats:  &runtime.MemStats{},
	debugText: &strings.Builder{},
}

func (system *debugSystem) Update(e *ecs.ECS) {
	var ok bool
	if system.debugEntry == nil {
		if system.debugEntry, ok = components.Debug.First(e.World); !ok {
			panic("no debug found")
		}
	}
	if system.cameraEntry == nil {
		if system.cameraEntry, ok = components.Camera.First(e.World); !ok {
			panic("no camera found")
		}
	}
	if system.displayEntry == nil {
		if system.displayEntry, ok = components.Display.First(e.World); !ok {
			panic("no display found")
		}
	}

	debug := components.Debug.Get(system.debugEntry)
	camera := components.Camera.Get(system.cameraEntry)

	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		debug.Enabled = !debug.Enabled
	}

	if !debug.Enabled {
		return
	}

	system.keys = inpututil.AppendPressedKeys(system.keys[:0])

	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		debug.DrawGrid = !debug.DrawGrid
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		debug.DrawCollisions = !debug.DrawCollisions
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		debug.DrawPositions = !debug.DrawPositions
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
		ecs.NewQuery(layers.Foreground, filter.Contains(entities.BulletTag)).Each(e.World, func(entry *donburi.Entry) {
			transform := components.Transform.Get(entry)
			velocity := components.Velocity.Get(entry)
			fmt.Printf("%v %v\n", transform, velocity)
		})
	}

	if system.ticks%uint64(config.C.TargetTPS*2) == 0 {
		runtime.ReadMemStats(system.memStats)
	}
	system.ticks++

	system.updateDebugText(debug, camera, system.debugText)
}

func (system *debugSystem) Draw(_ *ecs.ECS, screen *ebiten.Image) {
	utility.Debug.PrintDebugTextAt(screen, system.debugText.String(), &ebiten.DrawImageOptions{})
}

func (system *debugSystem) updateDebugText(debug *components.DebugData, camera *components.CameraData, builder *strings.Builder) {
	builder.Reset()

	ms := system.memStats

	cursorX, cursorY := ebiten.CursorPosition()
	worldPosition := camera.ScreenToWorldPosition(r2.Vec{X: float64(cursorX), Y: float64(cursorY)})
	screenPosition := camera.WorldToScreenPosition(worldPosition)

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
		debug.DrawGrid,
		debug.DrawCollisions,
		debug.DrawPositions,
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
		ebiten.IsVsyncEnabled(),
		ebiten.DeviceScaleFactor(),
		strings.Join(lo.Map(system.keys, func(item ebiten.Key, index int) string {
			return item.String()
		}), ", "),
		screenPosition.X, screenPosition.Y,
		worldPosition.X, worldPosition.Y,
		camera.Position.X, camera.Position.Y,
		camera.ZoomFactor,
		camera.Rotation,
		utility.FormatBytes(ms.Alloc), utility.FormatBytes(ms.TotalAlloc), utility.FormatBytes(ms.Sys),
		utility.FormatBytes(ms.NextGC), ms.NumGC)
}
