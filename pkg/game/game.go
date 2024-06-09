package game

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/pkg/resources"
	"github.com/ubootgame/ubootgame/pkg/scenes"
	"github.com/ubootgame/ubootgame/pkg/settings"
	"gonum.org/v1/gonum/spatial/r2"
)

type Game[S any] struct {
	settings         *settings.Settings[S]
	resourceRegistry *resources.Registry

	scene       scenes.Scene
	displayInfo DisplayInfo
}

type DisplayInfo struct {
	WindowSize, VirtualResolution r2.Vec
	ScalingFactor                 float64
}

func NewGame[S any](settings *settings.Settings[S], resourceRegistry *resources.Registry) *Game[S] {
	return &Game[S]{
		settings:         settings,
		resourceRegistry: resourceRegistry,
		displayInfo:      DisplayInfo{},
	}
}

func (g *Game[S]) LoadScene(scene scenes.Scene) error {
	if err := scene.Load(g.resourceRegistry); err == nil {
		g.scene = scene
		return nil
	} else {
		return err
	}
}

func (g *Game[S]) Update() error {
	if g.scene == nil {
		return errors.New("no scene loaded")
	}

	return g.scene.Update()
}

func (g *Game[S]) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.scene.Draw(screen)
}

func (g *Game[S]) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	var windowSize, virtualResolution r2.Vec

	windowSize.X, windowSize.Y = float64(outsideWidth), float64(outsideHeight)
	outerRatio := windowSize.X / windowSize.Y

	deviceScaleFactor := ebiten.Monitor().DeviceScaleFactor()

	if outerRatio <= g.settings.Display.Ratio {
		virtualResolution.X = float64(outsideWidth) * deviceScaleFactor
		virtualResolution.Y = virtualResolution.X / g.settings.Display.Ratio
	} else {
		virtualResolution.Y = float64(outsideHeight) * deviceScaleFactor
		virtualResolution.X = virtualResolution.Y * g.settings.Display.Ratio
	}

	if g.settings.Debug.FontScale != deviceScaleFactor || g.settings.Debug.FontFace == nil {
		g.settings.UpdateFontFace(deviceScaleFactor)
	}

	g.displayInfo.WindowSize = windowSize
	g.displayInfo.VirtualResolution = virtualResolution
	g.displayInfo.ScalingFactor = deviceScaleFactor

	return int(virtualResolution.X), int(virtualResolution.Y)
}

func (g *Game[S]) DisplayInfo() *DisplayInfo {
	return &g.displayInfo
}
