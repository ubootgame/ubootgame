package components

import (
	"github.com/yohamta/donburi"
	"gonum.org/v1/gonum/spatial/r2"
)

type PositionData struct {
	Center r2.Vec
	Size   r2.Vec
}

var Position = donburi.NewComponentType[PositionData]()
