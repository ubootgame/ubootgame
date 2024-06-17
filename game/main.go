package main

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/pkg/profile"
	"github.com/samber/do"
	"github.com/ubootgame/ubootgame/config"
	"github.com/ubootgame/ubootgame/framework/cli"
	"github.com/ubootgame/ubootgame/framework/game"
	"github.com/ubootgame/ubootgame/framework/graphics/display"
	"github.com/ubootgame/ubootgame/framework/input"
	"github.com/ubootgame/ubootgame/framework/resources"
	"github.com/ubootgame/ubootgame/framework/settings"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/scenes"
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

	injector := do.New()
	defer cleanUpServices(injector)

	prepareServices(injector)

	g := game.NewGame[internal.Settings](injector, scenes.Scenes)

	if err := g.Run("game"); err != nil {
		log.Fatal(err)
	}
}

func prepareServices(injector *do.Injector) {
	do.Provide(injector, display.NewDisplay[internal.Settings])
	do.Provide(injector, input.NewInput)
	do.Provide(injector, func(i *do.Injector) (settings.Provider[internal.Settings], error) {
		return settings.NewProvider[internal.Settings](i, config.DefaultSettings[internal.Settings]())
	})
	do.Provide(injector, func(i *do.Injector) (resources.Registry, error) {
		audioContext := audio.NewContext(44100)
		return resources.NewRegistry(i, audioContext)
	})
}

func cleanUpServices(injector *do.Injector) {
	if err := injector.Shutdown(); err != nil {
		log.Fatal(err)
	}
}
