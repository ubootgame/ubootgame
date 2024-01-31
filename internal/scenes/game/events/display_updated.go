package events

import (
	"github.com/yohamta/donburi/features/events"
	"gonum.org/v1/gonum/spatial/r2"
)

type DisplayUpdatedEventData struct {
	WindowSize, VirtualResolution r2.Vec
	ScalingFactor                 float64
}

var DisplayUpdatedEvent = events.NewEventType[DisplayUpdatedEventData]()
