package display

import (
	"gonum.org/v1/gonum/spatial/r2"
)

type Scale interface {
	ScreenScale(size, resolution r2.Vec) float64
	WorldSize(size r2.Vec) r2.Vec
}

type HScale float64

func (s HScale) ScreenScale(size, resolution r2.Vec) float64 {
	return float64(s) / (size.X / resolution.X)
}

func (s HScale) WorldSize(size r2.Vec) r2.Vec {
	ratio := size.X / size.Y
	return r2.Vec{X: float64(s), Y: float64(s) / ratio}
}

type VScale float64

func (s VScale) ScreenScale(size, resolution r2.Vec) float64 {
	return float64(s) / (size.Y / resolution.X)
}

func (s VScale) WorldSize(size r2.Vec) r2.Vec {
	ratio := size.X / size.Y
	return r2.Vec{X: float64(s) * ratio, Y: float64(s)}
}
