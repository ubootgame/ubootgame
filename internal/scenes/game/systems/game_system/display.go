package game_system

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/events"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/injector"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/systems"
	"github.com/yohamta/donburi"
)

type DisplaySystem struct {
	systems.BaseSystem

	display *game_system.DisplayData
}

func NewDisplaySystem() *DisplaySystem {
	system := &DisplaySystem{}
	system.Injector = injector.NewInjector([]injector.Injection{
		injector.Once([]injector.Injection{
			injector.Component(&system.display, game_system.Display),
		}),
	})
	return system
}

func (system *DisplaySystem) UpdateDisplay(w donburi.World, event events.DisplayUpdatedEventData) {
	system.Inject(w)

	system.display.WindowSize = event.WindowSize
	system.display.VirtualResolution = event.VirtualResolution
	system.display.ScalingFactor = event.ScalingFactor
}
