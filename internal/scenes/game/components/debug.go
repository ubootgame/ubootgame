package components

import "github.com/yohamta/donburi"

type DebugData struct {
	DrawResolvLines, DrawGrid bool
}

var Debug = donburi.NewComponentType[DebugData]()
