package game

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/framework"
	"gonum.org/v1/gonum/spatial/r2"
	"runtime/debug"
)

type Game[S any] struct {
	settings framework.SettingsService[S]
	scenes   framework.SceneService
	display  framework.DisplayService

	activeScene framework.Scene
}

func NewGame[S any](settings framework.SettingsService[S], scenes framework.SceneService, display framework.DisplayService) *Game[S] {
	game := &Game[S]{
		settings: settings,
		scenes:   scenes,
		display:  display,
	}
	game.ApplySettings()
	return game
}

func (g *Game[S]) ApplySettings() {
	ebiten.SetWindowTitle("U-Boot")
	ebiten.SetWindowSize(int(g.settings.Settings().Window.DefaultSize.X), int(g.settings.Settings().Window.DefaultSize.Y))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(g.settings.Settings().Internals.TPS)
	ebiten.SetVsyncEnabled(g.settings.Settings().Graphics.VSync)
	debug.SetGCPercent(100)
}

func (g *Game[S]) LoadScene(sceneID framework.SceneID) error {
	scene, err := g.scenes.Get(sceneID)
	if err != nil {
		return err
	}
	if err = scene.Load(); err != nil {
		return err
	}
	g.activeScene = scene
	return nil
}

func (g *Game[S]) Update() error {
	if g.activeScene == nil {
		return errors.New("no active scene")
	}

	return g.activeScene.Update()
}

func (g *Game[S]) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.activeScene.Draw(screen)
}

func (g *Game[S]) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	scaleFactor := ebiten.Monitor().DeviceScaleFactor()

	if g.settings.Settings().Debug.FontScale != scaleFactor || g.settings.Settings().Debug.FontFace == nil {
		g.settings.UpdateDebugFontFace(scaleFactor)
	}

	virtualResolution := g.display.UpdateVirtualResolution(r2.Vec{X: float64(outsideWidth), Y: float64(outsideHeight)}, scaleFactor)

	return int(virtualResolution.X), int(virtualResolution.Y)
}
