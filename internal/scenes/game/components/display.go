package components

import (
	"github.com/yohamta/donburi"
	"gonum.org/v1/gonum/spatial/r2"
)

type DisplayData struct {
	WindowSize, VirtualResolution r2.Vec
}

var Display = donburi.NewComponentType[DisplayData]()

func (screen *DisplayData) ScalingFactor() float64 {
	desiredRatio := screen.VirtualResolution.X / screen.VirtualResolution.Y
	outerRatio := screen.WindowSize.X / screen.WindowSize.Y
	scale := screen.VirtualResolution.Y / screen.WindowSize.Y
	if desiredRatio > outerRatio {
		scale *= desiredRatio / outerRatio
	}
	return scale
}

func (screen *DisplayData) Ratio() float64 {
	return screen.VirtualResolution.X / screen.VirtualResolution.Y
}
