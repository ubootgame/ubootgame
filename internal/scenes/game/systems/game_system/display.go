package game_system

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/events"
	"github.com/yohamta/donburi"
)

type displaySystem struct {
	entry *donburi.Entry
}

var Display = &displaySystem{}

func (system *displaySystem) UpdateDisplay(w donburi.World, event events.DisplayUpdatedEventData) {
	if system.entry == nil {
		var ok bool
		if system.entry, ok = game_system.Display.First(w); !ok {
			panic("no display found")
		}
	}

	display := game_system.Display.Get(system.entry)

	display.WindowSize = event.WindowSize
	display.VirtualResolution = event.VirtualResolution
	display.ScalingFactor = event.ScalingFactor
}
