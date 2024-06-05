package geometry

import (
	"github.com/yohamta/donburi"
	"gonum.org/v1/gonum/spatial/r2"
)

type TransformData struct {
	Center       r2.Vec
	Size         r2.Vec
	Rotate       float64
	FlipX, FlipY bool
}

var Transform = donburi.NewComponentType[TransformData]()
