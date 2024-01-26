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
	var width, height int
	outsideRatio := float64(outsideWidth) / float64(outsideHeight)
	desiredRatio := config.C.Ratio

	if outsideRatio >= desiredRatio {
		width = int(float64(outsideHeight) * desiredRatio)
		height = outsideHeight
	} else {
		width = outsideWidth
		height = int(float64(outsideWidth) / desiredRatio)
	}
	return width, height
}
