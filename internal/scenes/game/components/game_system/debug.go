package game_system

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
)

type DebugData struct {
	Enabled                                 bool
	DrawCollisions, DrawGrid, DrawPositions bool
	FontScale                               float64
	FontFace                                text.Face
}

var Debug = donburi.NewComponentType[DebugData]()
