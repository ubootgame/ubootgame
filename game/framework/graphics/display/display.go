package display

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/do"
	"github.com/ubootgame/ubootgame/framework/settings"
	"gonum.org/v1/gonum/spatial/r2"
)

type Display interface {
	WindowSize() r2.Vec
	VirtualResolution() r2.Vec
	UpdateVirtualResolution(windowSize r2.Vec, scaleFactor float64) r2.Vec
	WorldToScreen(v r2.Vec) (float64, float64)
}

type display[S any] struct {
	settingsProvider settings.Provider[S]

	windowSize, virtualResolution r2.Vec
}

func NewDisplay[S any](i *do.Injector) (Display, error) {
	settingsService := do.MustInvoke[settings.Provider[S]](i)

	d := &display[S]{settingsProvider: settingsService}
	d.UpdateVirtualResolution(settingsService.Settings().Window.DefaultSize, ebiten.Monitor().DeviceScaleFactor())

	return d, nil
}

func (display *display[S]) WindowSize() r2.Vec {
	return display.windowSize
}

func (display *display[S]) VirtualResolution() r2.Vec {
	return display.virtualResolution
}

func (display *display[S]) UpdateVirtualResolution(windowSize r2.Vec, scaleFactor float64) r2.Vec {
	outerRatio := windowSize.X / windowSize.Y

	display.windowSize = windowSize

	if outerRatio <= display.settingsProvider.Settings().Window.Ratio {
		display.virtualResolution.X = windowSize.X * scaleFactor
		display.virtualResolution.Y = display.virtualResolution.X / display.settingsProvider.Settings().Window.Ratio
	} else {
		display.virtualResolution.Y = windowSize.Y * scaleFactor
		display.virtualResolution.X = display.virtualResolution.Y * display.settingsProvider.Settings().Window.Ratio
	}

	return display.virtualResolution
}

func (display *display[S]) WorldToScreen(v r2.Vec) (float64, float64) {
	return v.X * display.virtualResolution.X, v.Y * display.virtualResolution.X
}
