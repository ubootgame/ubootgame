package config

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/framework/settings"
	"gonum.org/v1/gonum/spatial/r2"
)

func DefaultSettings[S interface{}]() settings.Settings[S] {
	return settings.Settings[S]{
		Window: settings.Window{
			Title:        "U-Boot",
			ResizingMode: ebiten.WindowResizingModeEnabled,
			DefaultSize:  r2.Vec{X: 1280, Y: 720},
			Ratio:        1280.0 / 720.0,
		},
		Debug: settings.Debug{
			Enabled:        true,
			DrawGrid:       true,
			DrawCollisions: true,
			DrawPositions:  true,
			FontScale:      1.0,
		},
		Graphics: settings.Graphics{
			VSync: true,
		},
		Internals: settings.Internals{
			TPS:       60,
			GCPercent: 100,
		},
	}
}
