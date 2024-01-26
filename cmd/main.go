package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	_ "image/png"
	"log"
)

func main() {
	ebiten.SetWindowTitle("U-Boot")
	ebiten.SetWindowSize(config.C.DefaultWidth, config.C.DefaultHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(config.C.TargetTPS)

	audioContext := audio.NewContext(44100)
	resourceRegistry := resources.NewRegistry(audioContext)

	gameScene := game.NewGameScene(resourceRegistry)

	err := resourceRegistry.RegisterResources(gameScene.Assets())
	if err != nil {
		panic(err)
	}

	if err := ebiten.RunGame(utility.NewGame(gameScene)); err != nil {
		log.Fatal(err)
	}
}
