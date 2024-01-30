package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/pkg/profile"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	_ "image/png"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
)

func main() {
	if ok, _ := utility.GetEnvBool("DEBUG"); ok {
		defer profile.Start(profile.CPUProfile, profile.MemProfile).Stop()

		go func() {
			_ = http.ListenAndServe("localhost:6060", nil)
		}()
	}

	ebiten.SetWindowTitle("U-Boot")
	ebiten.SetWindowSize(int(config.C.DefaultOuterSize.X), int(config.C.DefaultOuterSize.Y))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(config.C.TargetTPS)
	ebiten.SetVsyncEnabled(true)

	debug.SetGCPercent(100)

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
