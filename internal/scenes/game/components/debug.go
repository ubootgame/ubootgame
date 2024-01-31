package components

import "github.com/yohamta/donburi"

type DebugData struct {
	Enabled                                  bool
	DrawResolvLines, DrawGrid, DrawPositions bool
}

var Debug = donburi.NewComponentType[DebugData]()
