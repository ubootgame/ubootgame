package scenes

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/framework"
	"github.com/ubootgame/ubootgame/framework/services/resources"
	"sync"
)

type BaseScene[S any] struct {
	Settings  framework.SettingsService[S]
	Resources framework.ResourceService

	ResourceLibrary *resources.Library

	once   sync.Once
	loaded bool
}

func NewBaseScene[S any](settings framework.SettingsService[S], resources framework.ResourceService, resourceLibrary *resources.Library) *BaseScene[S] {
	return &BaseScene[S]{
		Settings:        settings,
		Resources:       resources,
		ResourceLibrary: resourceLibrary,
	}
}

func (scene *BaseScene[S]) Load() error {
	if err := scene.Resources.RegisterResources(scene.ResourceLibrary); err != nil {
		return err
	}
	scene.loaded = true
	return nil
}

func (scene *BaseScene[S]) Update() error {
	if !scene.loaded {
		return errors.New("scene not loaded")
	}
	return nil
}

func (scene *BaseScene[S]) Draw(_ *ebiten.Image) {
	panic("must be overridden")
}
