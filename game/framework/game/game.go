package game

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/do"
	"github.com/ubootgame/ubootgame/framework/graphics/display"
	"github.com/ubootgame/ubootgame/framework/settings"
	"gonum.org/v1/gonum/spatial/r2"
	"runtime/debug"
)

type Game[S interface{}] struct {
	injector *do.Injector

	settingsProvider settings.Provider[S]
	display          display.Display

	scenes      SceneMap
	activeScene Scene
}

func NewGame[S interface{}](i *do.Injector, scenes SceneMap) *Game[S] {
	return &Game[S]{
		injector:         i,
		settingsProvider: do.MustInvoke[settings.Provider[S]](i),
		display:          do.MustInvoke[display.Display](i),
		scenes:           scenes,
	}
}

func (g *Game[S]) Run(initialScene SceneID) error {
	g.ApplySettings()
	if err := g.LoadScene(initialScene); err != nil {
		return err
	}
	return ebiten.RunGame(g)
}

func (g *Game[S]) ApplySettings() {
	ebiten.SetWindowTitle("U-Boot")
	ebiten.SetWindowSize(int(g.settingsProvider.Settings().Window.DefaultSize.X), int(g.settingsProvider.Settings().Window.DefaultSize.Y))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(g.settingsProvider.Settings().Internals.TPS)
	ebiten.SetVsyncEnabled(g.settingsProvider.Settings().Graphics.VSync)
	debug.SetGCPercent(100)
}

func (g *Game[S]) LoadScene(sceneID SceneID) error {
	sceneConstructor, ok := g.scenes[sceneID]
	if !ok {
		return errors.New("scene not found")
	}
	scene := sceneConstructor(g.injector)
	if err := scene.Load(); err != nil {
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

	if g.settingsProvider.Settings().Debug.FontScale != scaleFactor || g.settingsProvider.Settings().Debug.FontFace == nil {
		g.settingsProvider.UpdateDebugFontFace(scaleFactor)
	}

	virtualResolution := g.display.UpdateVirtualResolution(r2.Vec{X: float64(outsideWidth), Y: float64(outsideHeight)}, scaleFactor)

	return int(virtualResolution.X), int(virtualResolution.Y)
}
