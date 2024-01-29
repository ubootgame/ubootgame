package systems

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/yohamta/donburi/ecs"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/opentype"
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

type debug struct {
	keys             []ebiten.Key
	resolvLinesImage *ebiten.Image
	memStats         *runtime.MemStats
	ticks            uint64
	scale            float64
	fontFace         font.Face
	debugText        *strings.Builder
}

var Debug = &debug{
	keys:      make([]ebiten.Key, 0),
	memStats:  &runtime.MemStats{},
	scale:     1.0,
	debugText: &strings.Builder{},
}

func (system *debug) Update(e *ecs.ECS) {
	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		config.C.Debug = !config.C.Debug
	}

	if !config.C.Debug {
		return
	}

	debugEntry, _ := components.Debug.First(e.World)
	debugData := components.Debug.Get(debugEntry)

	system.keys = inpututil.AppendPressedKeys(system.keys[:0])

	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		debugData.DrawGrid = !debugData.DrawGrid
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		debugData.DrawResolvLines = !debugData.DrawResolvLines
		if !debugData.DrawResolvLines && system.resolvLinesImage != nil {
			system.resolvLinesImage.Dispose()
			system.resolvLinesImage = nil
		}
	}

	if system.ticks%uint64(config.C.TargetTPS*2) == 0 {
		runtime.ReadMemStats(system.memStats)
	}
	system.ticks++

	scale := utility.CalculateScalingFactor()

	if scale != system.scale || system.fontFace == nil {
		system.scale = scale
		system.updateFontFace(scale)
	}
}

func (system *debug) Draw(e *ecs.ECS, screen *ebiten.Image) {
	if !config.C.Debug {
		return
	}

	debugEntry, _ := components.Debug.First(e.World)
	debugData := components.Debug.Get(debugEntry)

	cameraEntry, _ := components.Camera.First(e.World)
	cameraData := components.Camera.Get(cameraEntry)

	system.createDebugText(debugData, cameraData, system.debugText)
	system.printDebugTextAt(screen, system.debugText.String(), &ebiten.DrawImageOptions{})

	if config.C.Debug && debugData.DrawResolvLines {
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

func (system *debug) updateFontFace(scale float64) {
	var err error

	system.fontFace, err = opentype.NewFace(debugFont, &opentype.FaceOptions{
		Size:    defaultFontSize * scale,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (system *debug) createDebugText(debugData *components.DebugData, cameraData *components.CameraData, builder *strings.Builder) {
	builder.Reset()

	ms := system.memStats

	_, _ = fmt.Fprintf(builder, `(/ to toggle debug)
FPS: %.1f
TPS: %.1f
VSync: %v
Keys: %v
Device scale factor: %.2f
Draw grid (F1): %v
Draw resolv (F2): %v
Camera position: %.2f, %.2f
Camera zoom: %.2f
Camera rotation: %.2f
Alloc: %s
Total: %s
Sys: %s
NextGC: %s
NumGC: %d`,
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
		ebiten.IsVsyncEnabled(),
		strings.Join(lo.Map(system.keys, func(item ebiten.Key, index int) string {
			return item.String()
		}), ", "),
		ebiten.DeviceScaleFactor(),
		debugData.DrawGrid,
		debugData.DrawResolvLines,
		cameraData.Position.X, cameraData.Position.Y,
		cameraData.ZoomFactor,
		cameraData.Rotation,
		utility.FormatBytes(ms.Alloc), utility.FormatBytes(ms.TotalAlloc), utility.FormatBytes(ms.Sys),
		utility.FormatBytes(ms.NextGC), ms.NumGC)
}

func (system *debug) printDebugTextAt(screen *ebiten.Image, debugText string, opts *ebiten.DrawImageOptions) {
	utility.DrawTextAtWithOptions(screen, debugText, system.fontFace, opts)
}
