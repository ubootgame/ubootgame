package debug

import (
	"github.com/yohamta/donburi/features/events"
	"go/types"
)

var ToggleDebugEvent = events.NewEventType[types.Nil]()
var ToggleDrawGrid = events.NewEventType[types.Nil]()
var ToggleDrawCollisions = events.NewEventType[types.Nil]()
var ToggleDrawPositions = events.NewEventType[types.Nil]()
