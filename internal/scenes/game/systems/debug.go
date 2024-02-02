package systems

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities"
	"github.com/ubootgame/ubootgame/internal/scenes/game/events"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/opentype"
	"gonum.org/v1/gonum/spatial/r2"
	"image/color"
	"log"
	"runtime"
	"strings"
)

const dpi = 72

var debugFont *opentype.Font
var defaultFontSize = 14.0

func init() {
	var err error
	debugFont, err = opentype.Parse(gobold.TTF)
	if err != nil {
		log.Fatal(err)
	}
}

type debugSystem struct {
	debugEntry, cameraEntry, displayEntry *donburi.Entry
	keys                                  []ebiten.Key
	resolvLinesImage                      *ebiten.Image
	memStats                              *runtime.MemStats
	ticks                                 uint64
	fontScale                             float64
	fontFace                              font.Face
	debugText                             *strings.Builder
}

var Debug = &debugSystem{
	keys:      make([]ebiten.Key, 0),
	memStats:  &runtime.MemStats{},
	fontScale: 1.0,
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

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cursorX, cursorY := ebiten.CursorPosition()
		worldPosition := camera.ScreenToWorldPosition(r2.Vec{X: float64(cursorX), Y: float64(cursorY)})
		screenPosition := camera.WorldToScreenPosition(worldPosition)

		fmt.Printf(`Cursor position: %v, %v,
World position: %.2f, %.2f
Screen position: %.2f, %.2f\n`,
			cursorX, cursorY,
			worldPosition.X, worldPosition.Y,
			screenPosition.X, screenPosition.Y)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		debug.DrawGrid = !debug.DrawGrid
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		debug.DrawResolvLines = !debug.DrawResolvLines
		if !debug.DrawResolvLines && system.resolvLinesImage != nil {
			system.resolvLinesImage.Dispose()
			system.resolvLinesImage = nil
		}
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

func (system *debugSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	debug := components.Debug.Get(system.debugEntry)

	system.printDebugTextAt(screen, system.debugText.String(), &ebiten.DrawImageOptions{})

	if debug.Enabled && debug.DrawResolvLines {
		if system.resolvLinesImage == nil {
			system.resolvLinesImage = ebiten.NewImage(screen.Bounds().Size().X, screen.Bounds().Size().Y)
			system.resolvLinesImage.Clear()

			spaceEntry, _ := components.Space.First(e.World)
			space := components.Space.Get(spaceEntry)

			for x := 0; x < space.Width(); x++ {
				for y := 0; y < space.Height(); y++ {
					sx, sy := space.SpaceToWorld(x, y)
					vector.StrokeRect(system.resolvLinesImage, float32(sx), float32(sy), float32(space.CellWidth), float32(space.CellHeight), 1, color.RGBA{R: 10, G: 10, B: 10, A: 10}, false)
				}
			}
		}

		screen.DrawImage(system.resolvLinesImage, nil)
	}
}

func (system *debugSystem) UpdateFontFace(_ donburi.World, event events.DisplayUpdatedEventData) {
	if system.fontScale != event.ScalingFactor || system.fontFace == nil {
		var err error

		system.fontScale = event.ScalingFactor
		system.fontFace, err = opentype.NewFace(debugFont, &opentype.FaceOptions{
			Size:    defaultFontSize * event.ScalingFactor,
			DPI:     dpi,
			Hinting: font.HintingVertical,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (system *debugSystem) updateDebugText(debug *components.DebugData, camera *components.CameraData, builder *strings.Builder) {
	builder.Reset()

	ms := system.memStats

	cursorX, cursorY := ebiten.CursorPosition()
	worldPosition := camera.ScreenToWorldPosition(r2.Vec{X: float64(cursorX), Y: float64(cursorY)})
	screenPosition := camera.WorldToScreenPosition(worldPosition)

	_, _ = fmt.Fprintf(builder, `(/ to toggle debugSystem)
Draw grid (F1): %v
Draw resolv (F2): %v
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
		debug.DrawResolvLines,
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

func (system *debugSystem) printDebugTextAt(screen *ebiten.Image, debugText string, opts *ebiten.DrawImageOptions) {
	utility.DrawTextAtWithOptions(screen, debugText, system.fontFace, opts)
}
