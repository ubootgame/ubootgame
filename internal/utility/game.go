package utility

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
)

type Scene interface {
	Assets() *resources.Library
	Update()
	Draw(screen *ebiten.Image)
}

type Game struct {
	scene Scene
}

func NewGame(scene Scene) *Game {
	g := &Game{
		scene: scene,
	}

	return g
}

func (g *Game) Update() error {
	g.scene.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.scene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	config.C.ActualOuterSize.X, config.C.ActualOuterSize.Y = float64(outsideWidth), float64(outsideHeight)
	return int(config.C.VirtualResolution.X), int(config.C.VirtualResolution.Y)
}
