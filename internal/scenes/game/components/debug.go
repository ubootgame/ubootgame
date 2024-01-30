package components

import "github.com/yohamta/donburi"

type DebugData struct {
	DrawResolvLines, DrawGrid, DrawPositions bool
}

var Debug = donburi.NewComponentType[DebugData]()
