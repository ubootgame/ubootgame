package display

import (
	"github.com/ubootgame/ubootgame/framework"
)

type Service[S any] struct {
	settings framework.SettingsService[S]

	virtualResolution virtualResolution
}

type virtualResolution struct {
	x, y float64
}

func NewService[S any](settings framework.SettingsService[S]) *Service[S] {
	return &Service[S]{settings: settings}
}

func (service *Service[S]) VirtualResolution() (float64, float64) {
	return service.virtualResolution.x, service.virtualResolution.y
}

func (service *Service[S]) UpdateVirtualResolution(width, height int, scaleFactor float64) (float64, float64) {
	outerRatio := float64(width) / float64(height)

	if outerRatio <= service.settings.Settings().Window.Ratio {
		service.virtualResolution.x = float64(width) * scaleFactor
		service.virtualResolution.y = service.virtualResolution.x / service.settings.Settings().Window.Ratio
	} else {
		service.virtualResolution.y = float64(height) * scaleFactor
		service.virtualResolution.x = service.virtualResolution.y * service.settings.Settings().Window.Ratio
	}

	return service.virtualResolution.x, service.virtualResolution.y
}
