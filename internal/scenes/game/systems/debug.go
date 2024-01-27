package systems

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"image/color"
	"strings"
)

var keys = make([]ebiten.Key, 0)

func UpdateDebug(e *ecs.ECS) {
	debugEntry, _ := components.Debug.First(e.World)
	debugData := components.Debug.Get(debugEntry)

	keys = inpututil.AppendPressedKeys(keys[:0])

	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		config.C.Debug = !config.C.Debug
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		debugData.DrawGrid = !debugData.DrawGrid
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		debugData.DrawResolvLines = !debugData.DrawResolvLines
	}
}

func DrawDebug(e *ecs.ECS, screen *ebiten.Image) {
	if !config.C.Debug {
		return
	}

	debugEntry, _ := components.Debug.First(e.World)
	debugData := components.Debug.Get(debugEntry)

	cameraEntry, _ := donburi.NewQuery(filter.Contains(components.Camera)).First(e.World)
	cameraData := components.Camera.Get(cameraEntry)

	debugText := fmt.Sprintf(`(/ to toggle debug)
FPS: %.1f
TPS: %.1f
VSync: %v
Keys: %v
Device scale factor: %.2f
Draw grid (F1): %v
Draw resolv (F2): %v
Camera position: %.2f, %.2f
Camera zoom: %.2f
Camera rotation: %.2f`,
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
		ebiten.IsVsyncEnabled(),
		strings.Join(lo.Map(keys, func(item ebiten.Key, index int) string {
			return item.String()
		}), ", "),
		ebiten.DeviceScaleFactor(),
		debugData.DrawGrid,
		debugData.DrawResolvLines,
		cameraData.Position.X, cameraData.Position.Y,
		cameraData.ZoomFactor,
		cameraData.Rotation)

	printDebugTextAt(screen, debugText, ebiten.GeoM{})

	if debugData.DrawResolvLines {
		spaceEntry, _ := components.Space.First(e.World)
		space := components.Space.Get(spaceEntry)

		for x := 0; x < space.Width(); x++ {
			for y := 0; y < space.Height(); y++ {
				sx, sy := space.SpaceToWorld(x, y)
				vector.StrokeRect(screen, float32(sx), float32(sy), float32(space.CellWidth), float32(space.CellHeight), 1, color.RGBA{R: 10, G: 10, B: 10, A: 10}, false)
			}
		}
	}
}

func printDebugTextAt(screen *ebiten.Image, debugText string, geom ebiten.GeoM) {
	debugImage := ebiten.NewImage(int(config.C.VirtualResolution.X/4), int(config.C.VirtualResolution.Y/4))

	desiredRatio := config.C.VirtualResolution.X / config.C.VirtualResolution.Y
	outerRatio := config.C.ActualOuterSize.X / config.C.ActualOuterSize.Y

	scale := config.C.VirtualResolution.Y / config.C.ActualOuterSize.Y

	if desiredRatio > outerRatio {
		scale *= desiredRatio / outerRatio
	}

	ebitenutil.DebugPrint(debugImage, fmt.Sprintf(debugText))

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(scale, scale)
	opts.GeoM.Concat(geom)
	opts.Filter = ebiten.FilterLinear

	screen.DrawImage(debugImage, opts)
}
