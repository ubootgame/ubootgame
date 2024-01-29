package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"runtime"
)

var normalFont font.Face

func init() {
	tt, err := opentype.Parse(gobold.TTF)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	normalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    14,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type game struct {
	memStats  *runtime.MemStats
	ticks     uint64
	debugText string
}

var Game = &game{
	memStats: &runtime.MemStats{},
}

func (g *game) Update() error {
	if g.ticks%30 == 0 {
		runtime.ReadMemStats(g.memStats)
	}
	g.ticks++

	ms := g.memStats

	g.debugText = fmt.Sprintf(`
FPS: %.1f
TPS: %.1f
VSync: %v
Device scale factor: %.2f
Alloc: %s
Total: %s
Sys: %s
NextGC: %s
NumGC: %d`,
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
		ebiten.IsVsyncEnabled(),
		ebiten.DeviceScaleFactor(),
		formatBytes(ms.Alloc), formatBytes(ms.TotalAlloc), formatBytes(ms.Sys),
		formatBytes(ms.NextGC), ms.NumGC)

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	drawTextAt(screen, g.debugText, normalFont, 0, 0, colornames.White)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetVsyncEnabled(false)
	if err := ebiten.RunGame(Game); err != nil {
		log.Fatal(err)
	}
}

func formatBytes(b uint64) string {
	if b >= 1073741824 {
		return fmt.Sprintf("%0.2f GiB", float64(b)/1073741824)
	} else if b >= 1048576 {
		return fmt.Sprintf("%0.2f MiB", float64(b)/1048576)
	} else if b >= 1024 {
		return fmt.Sprintf("%0.2f KiB", float64(b)/1024)
	} else {
		return fmt.Sprintf("%d B", b)
	}
}

func drawTextAt(screen *ebiten.Image, s string, f font.Face, x, y int, clr color.Color) {
	y2 := y + f.Metrics().Height.Round()

	text.Draw(screen, s, f, x, y2, clr)
}
