package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/ubootgame/ubootgame/internal/assets"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes"
	"github.com/ubootgame/ubootgame/internal/utility"
	"log"
)

func main() {
	ebiten.SetWindowSize(config.C.Width, config.C.Height)
	ebiten.SetWindowTitle("U-Boot")

	audioContext := audio.NewContext(44100)
	resourceLoader := utility.CreateResourceLoader(audioContext)

	resourceLoader.ImageRegistry.Assign(assets.ImageResources)
	for id := range assets.ImageResources {
		resourceLoader.LoadImage(id)
	}

	gameScene := scenes.NewGameScene(resourceLoader)

	if err := ebiten.RunGame(NewGame(gameScene)); err != nil {
		log.Fatal(err)
	}
}
