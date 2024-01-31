package components

import (
	"github.com/yohamta/donburi"
	"gonum.org/v1/gonum/spatial/r2"
)

type DisplayData struct {
	WindowSize, VirtualResolution r2.Vec
	ScalingFactor                 float64
}

var Display = donburi.NewComponentType[DisplayData]()

func (display *DisplayData) Ratio() float64 {
	return display.VirtualResolution.X / display.VirtualResolution.Y
}
