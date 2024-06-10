package player

import (
	"github.com/yohamta/donburi/features/events"
	"go/types"
)

var MoveLeftEvent = events.NewEventType[types.Nil]()
var MoveRightEvent = events.NewEventType[types.Nil]()
var ShootEvent = events.NewEventType[types.Nil]()
