package components

import "github.com/yohamta/donburi"

type DebugData struct {
	Enabled                                 bool
	DrawCollisions, DrawGrid, DrawPositions bool
}

var Debug = donburi.NewComponentType[DebugData]()
