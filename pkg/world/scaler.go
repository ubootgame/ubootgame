package world

import (
	"gonum.org/v1/gonum/spatial/r2"
)

type Scaler interface {
	GetNormalizedSizeAndScale(size r2.Vec) (normalizedSize r2.Vec, baseScale float64, localScale float64)
}

type HScale float64

func (s HScale) GetNormalizedSizeAndScale(size r2.Vec) (r2.Vec, float64, float64) {
	scale := WorldSizeBase / size.X
	return r2.Vec{X: WorldSizeBase, Y: size.Y * scale}, scale, float64(s) / WorldSizeBase
}

type VScale float64

func (s VScale) GetNormalizedSizeAndScale(size r2.Vec) (r2.Vec, float64, float64) {
	scale := WorldSizeBase / size.X
	ratio := size.X / size.Y
	return r2.Vec{X: WorldSizeBase, Y: size.Y * scale}, scale, (float64(s) * ratio) / WorldSizeBase
}
