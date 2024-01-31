package systems

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
	"gonum.org/v1/gonum/spatial/r2"
)

type displaySystem struct {
	entry *donburi.Entry
}

var Display = &displaySystem{}

func (system *displaySystem) UpdateScreen(w donburi.World, event ScreenUpdatedEventData) {
	if system.entry == nil {
		var ok bool
		if system.entry, ok = components.Display.First(w); !ok {
			panic("no display found")
		}
	}

	display := components.Display.Get(system.entry)

	display.WindowSize = event.WindowSize
	display.VirtualResolution = event.VirtualResolution
}

type ScreenUpdatedEventData struct {
	WindowSize, VirtualResolution r2.Vec
}

var ScreenUpdatedEvent = events.NewEventType[ScreenUpdatedEventData]()
