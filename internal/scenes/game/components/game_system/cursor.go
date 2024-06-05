package game_system

import (
	"github.com/yohamta/donburi"
	"gonum.org/v1/gonum/spatial/r2"
)

type CursorData struct {
	ScreenPosition, WorldPosition r2.Vec
}

var Cursor = donburi.NewComponentType[CursorData]()
