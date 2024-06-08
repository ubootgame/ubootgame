package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/pkg/profile"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/framework"
	"github.com/ubootgame/ubootgame/internal/framework/resources"
	"github.com/ubootgame/ubootgame/internal/scenes/game"
	_ "image/png"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
)

func main() {
	if ok, _ := framework.GetEnvBool("DEBUG"); ok {
		defer profile.Start(profile.CPUProfile, profile.MemProfile).Stop()

		go func() {
			_ = http.ListenAndServe("localhost:6060", nil)
		}()
	}

	settings := internal.NewSettings()

	ebiten.SetWindowTitle("U-Boot")
	ebiten.SetWindowSize(int(settings.DefaultWindowSize.X), int(settings.DefaultWindowSize.Y))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(settings.TargetTPS)
	ebiten.SetVsyncEnabled(true)

	debug.SetGCPercent(100)

	audioContext := audio.NewContext(44100)
	resourceRegistry := resources.NewRegistry(audioContext)

	g := framework.NewGame(settings, resourceRegistry)

	if err := g.LoadScene(game.NewScene(settings)); err != nil {
		panic(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
