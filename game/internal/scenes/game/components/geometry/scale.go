package geometry

import (
	"github.com/yohamta/donburi"
	"gonum.org/v1/gonum/spatial/r2"
)

type ScaleData struct {
	NormalizedSize  r2.Vec
	NormalizedScale float64
}

var Scale = donburi.NewComponentType[ScaleData]()
