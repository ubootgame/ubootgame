package utility

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	"image"
)

type Scene interface {
	Assets() *resources.Library
	Update()
	Draw(screen *ebiten.Image)
}

type Game struct {
	bounds image.Rectangle
	scene  Scene
}

func NewGame(scene Scene) *Game {
	g := &Game{
		bounds: image.Rectangle{},
		scene:  scene,
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
	g.bounds = image.Rect(0, 0, outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}
