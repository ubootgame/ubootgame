package utility

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	"gonum.org/v1/gonum/spatial/r2"
)

type Scene interface {
	Assets() *resources.Library
	AdjustScreen(windowSize r2.Vec, virtualResolution r2.Vec, scalingFactor float64)
	Update()
	Draw(screen *ebiten.Image)
}

type Game struct {
	scene Scene
}

func NewGame(scene Scene) *Game {
	return &Game{scene: scene}
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
	var windowSize, virtualResolution r2.Vec

	windowSize.X, windowSize.Y = float64(outsideWidth), float64(outsideHeight)
	outerRatio := windowSize.X / windowSize.Y

	if outerRatio <= config.C.Ratio {
		virtualResolution.X = float64(outsideWidth) * ebiten.DeviceScaleFactor()
		virtualResolution.Y = virtualResolution.X / config.C.Ratio
	} else {
		virtualResolution.Y = float64(outsideHeight) * ebiten.DeviceScaleFactor()
		virtualResolution.X = virtualResolution.Y * config.C.Ratio
	}

	g.scene.AdjustScreen(windowSize, virtualResolution, ebiten.DeviceScaleFactor())

	return int(virtualResolution.X), int(virtualResolution.Y)
}
