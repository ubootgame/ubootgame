package internal

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"gonum.org/v1/gonum/spatial/r2"
)

type Settings struct {
	DefaultWindowSize r2.Vec
	Ratio             float64
	TargetTPS         int
	Debug             Debug
	Display           Display
}

type Debug struct {
	Enabled                                 bool
	DrawCollisions, DrawGrid, DrawPositions bool
	FontScale                               float64
	FontFace                                text.Face
}

type Display struct {
	WindowSize, VirtualResolution r2.Vec
	ScalingFactor                 float64
}

func NewSettings() *Settings {
	return &Settings{
		DefaultWindowSize: r2.Vec{X: 1280, Y: 720},
		Ratio:             1280.0 / 720.0,
		TargetTPS:         60,
		Debug: Debug{
			Enabled:        true,
			DrawGrid:       true,
			DrawCollisions: true,
			DrawPositions:  true,
			FontScale:      1.0,
		},
	}
}
