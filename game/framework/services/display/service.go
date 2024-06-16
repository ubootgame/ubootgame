package display

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/framework"
	"gonum.org/v1/gonum/spatial/r2"
)

type Service[S any] struct {
	settings framework.SettingsService[S]

	windowSize, virtualResolution r2.Vec
}

func NewService[S any](settings framework.SettingsService[S]) *Service[S] {
	service := &Service[S]{settings: settings}
	service.UpdateVirtualResolution(settings.Settings().Window.DefaultSize, ebiten.Monitor().DeviceScaleFactor())
	return service
}

func (service *Service[S]) WindowSize() r2.Vec {
	return service.windowSize
}

func (service *Service[S]) VirtualResolution() r2.Vec {
	return service.virtualResolution
}

func (service *Service[S]) UpdateVirtualResolution(windowSize r2.Vec, scaleFactor float64) r2.Vec {
	outerRatio := windowSize.X / windowSize.Y

	service.windowSize = windowSize

	if outerRatio <= service.settings.Settings().Window.Ratio {
		service.virtualResolution.X = windowSize.X * scaleFactor
		service.virtualResolution.Y = service.virtualResolution.X / service.settings.Settings().Window.Ratio
	} else {
		service.virtualResolution.Y = windowSize.Y * scaleFactor
		service.virtualResolution.X = service.virtualResolution.Y * service.settings.Settings().Window.Ratio
	}

	return service.virtualResolution
}

func (service *Service[S]) WorldToScreen(v r2.Vec) (float64, float64) {
	return v.X * service.virtualResolution.X, v.Y * service.virtualResolution.X
}
