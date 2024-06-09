package scenes

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/pkg/resources"
	framework "github.com/ubootgame/ubootgame/pkg/settings"
	"sync"
)

type Scene interface {
	Load(resourceRegistry *resources.Registry) error
	Update() error
	Draw(screen *ebiten.Image)
}

type BaseScene[S any] struct {
	once sync.Once

	loaded    bool
	Resources *resources.Library

	Settings         *framework.Settings[S]
	ResourceRegistry *resources.Registry
}

func NewBaseScene[S any](settings *framework.Settings[S], resources *resources.Library) *BaseScene[S] {
	return &BaseScene[S]{
		Settings:  settings,
		Resources: resources,
	}
}

func (scene *BaseScene[S]) Load(resourceRegistry *resources.Registry) error {
	scene.ResourceRegistry = resourceRegistry

	if err := resourceRegistry.RegisterResources(scene.Resources); err == nil {
		scene.loaded = true
		return nil
	} else {
		return err
	}
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
