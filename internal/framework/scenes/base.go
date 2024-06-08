package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/framework/resources"
	"sync"
)

type Scene interface {
	Load(resourceRegistry *resources.Registry) error
	Update() error
	Draw(screen *ebiten.Image)
}

type BaseScene struct {
	once sync.Once

	loaded    bool
	Resources *resources.Library

	Settings         *internal.Settings
	ResourceRegistry *resources.Registry
}

func NewBaseScene(settings *internal.Settings, resources *resources.Library) *BaseScene {
	return &BaseScene{
		Settings:  settings,
		Resources: resources,
	}
}

func (scene *BaseScene) Load(resourceRegistry *resources.Registry) error {
	scene.ResourceRegistry = resourceRegistry

	if err := resourceRegistry.RegisterResources(scene.Resources); err == nil {
		scene.loaded = true
		return nil
	} else {
		return err
	}
}

func (scene *BaseScene) Update() {
	if !scene.loaded {
		panic("scene not loaded")
	}
}

func (scene *BaseScene) Draw(_ *ebiten.Image) {
	panic("must be overridden")
}
