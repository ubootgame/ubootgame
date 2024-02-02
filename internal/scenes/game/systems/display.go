package systems

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
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
		if system.entry, ok = components.Display.First(w); !ok {
			panic("no display found")
		}
	}

	display := components.Display.Get(system.entry)

	display.WindowSize = event.WindowSize
	display.VirtualResolution = event.VirtualResolution
	display.ScalingFactor = event.ScalingFactor
}
