package camera

import (
	"github.com/yohamta/donburi/features/events"
	"go/types"
)

var PanLeftEvent = events.NewEventType[types.Nil]()
var PanRightEvent = events.NewEventType[types.Nil]()
var PanUpEvent = events.NewEventType[types.Nil]()
var PanDownEvent = events.NewEventType[types.Nil]()
var ZoomInEvent = events.NewEventType[types.Nil]()
var ZoomOutEvent = events.NewEventType[types.Nil]()
var RotateLeftEvent = events.NewEventType[types.Nil]()
var RotateRightEvent = events.NewEventType[types.Nil]()
