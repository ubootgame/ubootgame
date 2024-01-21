package main

import (
	"github.com/ubootgame/ubootgame/internal/scenes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
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

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("U-Boot")

	gameScene := scenes.NewGameScene()

	if err := ebiten.RunGame(NewGame(gameScene)); err != nil {
		log.Fatal(err)
	}
}
