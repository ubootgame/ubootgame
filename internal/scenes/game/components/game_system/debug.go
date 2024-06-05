package game_system

import (
	"github.com/yohamta/donburi"
	"golang.org/x/image/font"
)

type DebugData struct {
	Enabled                                 bool
	DrawCollisions, DrawGrid, DrawPositions bool
	FontScale                               float64
	FontFace                                font.Face
}

var Debug = donburi.NewComponentType[DebugData]()
