package game

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/pkg"
	"runtime/debug"
)

type Game[S any] struct {
	settings pkg.SettingsService[S]
	scenes   pkg.SceneService
	display  pkg.DisplayService

	activeScene pkg.Scene
}

func NewGame[S any](settings pkg.SettingsService[S], scenes pkg.SceneService, display pkg.DisplayService) *Game[S] {
	return &Game[S]{
		settings: settings,
		scenes:   scenes,
		display:  display,
	}
}

func (g *Game[S]) ApplySettings() {
	ebiten.SetWindowTitle("U-Boot")
	ebiten.SetWindowSize(int(g.settings.Settings().Window.DefaultSize.X), int(g.settings.Settings().Window.DefaultSize.Y))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(g.settings.Settings().Internals.TPS)
	ebiten.SetVsyncEnabled(g.settings.Settings().Graphics.VSync)
	debug.SetGCPercent(100)
}

func (g *Game[S]) LoadScene(sceneID pkg.SceneID) error {
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

	vx, vy := g.display.UpdateVirtualResolution(outsideWidth, outsideHeight, scaleFactor)
	return int(vx), int(vy)
}
