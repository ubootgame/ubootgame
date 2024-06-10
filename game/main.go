package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/pkg/profile"
	"github.com/ubootgame/ubootgame/framework"
	"github.com/ubootgame/ubootgame/framework/cli"
	"github.com/ubootgame/ubootgame/framework/game"
	"github.com/ubootgame/ubootgame/framework/services/display"
	"github.com/ubootgame/ubootgame/framework/services/resources"
	"github.com/ubootgame/ubootgame/framework/services/scenes"
	"github.com/ubootgame/ubootgame/framework/services/settings"
	"github.com/ubootgame/ubootgame/internal"
	gameScene "github.com/ubootgame/ubootgame/internal/scenes/game"
	"gonum.org/v1/gonum/spatial/r2"
	_ "image/png"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	if ok, _ := cli.GetEnvBool("DEBUG"); ok {
		defer profile.Start(profile.CPUProfile, profile.MemProfile).Stop()

		go func() {
			_ = http.ListenAndServe("localhost:6060", nil)
		}()
	}

	s := &settings.Settings[internal.Settings]{
		Window: settings.Window{
			Title:        "U-Boot",
			ResizingMode: ebiten.WindowResizingModeEnabled,
			DefaultSize:  r2.Vec{X: 1280, Y: 720},
			Ratio:        1280.0 / 720.0,
		},
		Debug: settings.Debug{
			Enabled:        true,
			DrawGrid:       true,
			DrawCollisions: true,
			DrawPositions:  true,
			FontScale:      1.0,
		},
		Graphics: settings.Graphics{
			VSync: true,
		},
		Internals: settings.Internals{
			TPS:       60,
			GCPercent: 100,
		},
		Game: internal.Settings{},
	}

	settingsService := settings.NewService(s)

	audioContext := audio.NewContext(44100)
	resourceService := resources.NewService(audioContext)

	displayService := display.NewService[internal.Settings](settingsService)

	sceneService := scenes.NewService(scenes.SceneMap{
		"game": func() framework.Scene {
			return gameScene.NewScene(settingsService, displayService, resourceService)
		},
	})

	g := game.NewGame[internal.Settings](settingsService, sceneService, displayService)
	g.ApplySettings()

	if err := g.LoadScene("game"); err != nil {
		panic(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
