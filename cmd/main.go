package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/pkg/profile"
	"github.com/ubootgame/ubootgame/internal"
	gameScene "github.com/ubootgame/ubootgame/internal/scenes/game"
	"github.com/ubootgame/ubootgame/pkg/cli"
	"github.com/ubootgame/ubootgame/pkg/game"
	"github.com/ubootgame/ubootgame/pkg/resources"
	"github.com/ubootgame/ubootgame/pkg/settings"
	_ "image/png"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
)

func main() {
	if ok, _ := cli.GetEnvBool("DEBUG"); ok {
		defer profile.Start(profile.CPUProfile, profile.MemProfile).Stop()

		go func() {
			_ = http.ListenAndServe("localhost:6060", nil)
		}()
	}

	s := settings.NewSettings(&internal.Settings{})

	ebiten.SetWindowTitle("U-Boot")
	ebiten.SetWindowSize(int(s.Display.DefaultWindowSize.X), int(s.Display.DefaultWindowSize.Y))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(s.TargetTPS)
	ebiten.SetVsyncEnabled(true)

	debug.SetGCPercent(100)

	audioContext := audio.NewContext(44100)
	resourceRegistry := resources.NewRegistry(audioContext)

	g := game.NewGame(s, resourceRegistry)

	if err := g.LoadScene(gameScene.NewScene(s, g.DisplayInfo())); err != nil {
		panic(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
