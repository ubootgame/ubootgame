package settings

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/gobold"
	"gonum.org/v1/gonum/spatial/r2"
	"log"
)

var defaultDebugFontSize = 12.0

type Settings[S any] struct {
	TargetTPS int
	Debug     Debug
	Display   Display
	Game      *S
}

type Debug struct {
	Enabled                                 bool
	DrawCollisions, DrawGrid, DrawPositions bool
	FontScale                               float64
	FontFace                                text.Face
}

type Display struct {
	DefaultWindowSize r2.Vec
	Ratio             float64
}

func NewSettings[S any](gameSettings *S) *Settings[S] {
	return &Settings[S]{
		TargetTPS: 60,
		Display: Display{
			DefaultWindowSize: r2.Vec{X: 1280, Y: 720},
			Ratio:             1280.0 / 720.0,
		},
		Debug: Debug{
			Enabled:        true,
			DrawGrid:       true,
			DrawCollisions: true,
			DrawPositions:  true,
			FontScale:      1.0,
		},
		Game: gameSettings,
	}
}

func (s *Settings[S]) UpdateFontFace(scalingFactor float64) {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(gobold.TTF))
	if err != nil {
		log.Fatal(err)
	}

	s.Debug.FontFace = &text.GoTextFace{
		Source: source,
		Size:   defaultDebugFontSize * scalingFactor,
	}
}
