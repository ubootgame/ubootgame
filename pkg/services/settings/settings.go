package settings

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"gonum.org/v1/gonum/spatial/r2"
)

type Settings[S any] struct {
	Debug     Debug
	Window    Window
	Graphics  Graphics
	Internals Internals
	Game      S
}

type Debug struct {
	Enabled                                 bool
	DrawCollisions, DrawGrid, DrawPositions bool
	FontScale                               float64
	FontFace                                text.Face
}

type Window struct {
	Title        string
	ResizingMode ebiten.WindowResizingModeType
	DefaultSize  r2.Vec
	Ratio        float64
}

type Graphics struct {
	VSync bool
}

type Internals struct {
	TPS, GCPercent int
}
