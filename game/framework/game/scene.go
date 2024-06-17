package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/do"
)

type SceneID string

type Scene interface {
	Load() error
	Update() error
	Draw(screen *ebiten.Image)
}

type SceneMap map[SceneID]func(i *do.Injector) Scene
