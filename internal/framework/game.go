package framework

import (
	"bytes"
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/framework/resources"
	"github.com/ubootgame/ubootgame/internal/framework/scenes"
	"golang.org/x/image/font/gofont/gobold"
	"gonum.org/v1/gonum/spatial/r2"
	"log"
)

var defaultFontSize = 12.0

type Game struct {
	scene            scenes.Scene
	settings         *internal.Settings
	resourceRegistry *resources.Registry
}

func NewGame(settings *internal.Settings, resourceRegistry *resources.Registry) *Game {
	return &Game{settings: settings, resourceRegistry: resourceRegistry}
}

func (g *Game) LoadScene(scene scenes.Scene) error {
	if err := scene.Load(g.resourceRegistry); err == nil {
		g.scene = scene
		return nil
	} else {
		return err
	}
}

func (g *Game) Update() error {
	if g.scene == nil {
		return errors.New("no scene loaded")
	}

	return g.scene.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.scene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	var windowSize, virtualResolution r2.Vec

	windowSize.X, windowSize.Y = float64(outsideWidth), float64(outsideHeight)
	outerRatio := windowSize.X / windowSize.Y

	deviceScaleFactor := ebiten.Monitor().DeviceScaleFactor()

	if outerRatio <= g.settings.Ratio {
		virtualResolution.X = float64(outsideWidth) * deviceScaleFactor
		virtualResolution.Y = virtualResolution.X / g.settings.Ratio
	} else {
		virtualResolution.Y = float64(outsideHeight) * deviceScaleFactor
		virtualResolution.X = virtualResolution.Y * g.settings.Ratio
	}

	if g.settings.Debug.FontScale != deviceScaleFactor || g.settings.Debug.FontFace == nil {
		s, err := text.NewGoTextFaceSource(bytes.NewReader(gobold.TTF))
		if err != nil {
			log.Fatal(err)
		}

		g.settings.Debug.FontFace = &text.GoTextFace{
			Source: s,
			Size:   defaultFontSize * deviceScaleFactor,
		}
	}

	g.settings.Display.WindowSize = windowSize
	g.settings.Display.VirtualResolution = virtualResolution
	g.settings.Display.ScalingFactor = deviceScaleFactor

	return int(virtualResolution.X), int(virtualResolution.Y)
}
