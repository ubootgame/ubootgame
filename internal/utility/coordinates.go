package utility

import (
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"gonum.org/v1/gonum/spatial/r2"
	"math"
)

func SetCameraMatrix(camera *components.CameraData) {
	camera.Matrix.Reset()
	camera.Matrix.Translate(-(camera.Position.X), -(camera.Position.Y))
	camera.Matrix.Scale(camera.ZoomFactor, camera.ZoomFactor)
	camera.Matrix.Rotate(float64(camera.Rotation) * 2 * math.Pi / 360)
	camera.Matrix.Translate(0.5, 0.5/config.C.Ratio)
	camera.Matrix.Scale(config.C.VirtualResolution.X, config.C.VirtualResolution.X)
}

func CalculateScreenScalingFactor() float64 {
	desiredRatio := config.C.VirtualResolution.X / config.C.VirtualResolution.Y
	outerRatio := config.C.ActualOuterSize.X / config.C.ActualOuterSize.Y
	scale := config.C.VirtualResolution.Y / config.C.ActualOuterSize.Y
	if desiredRatio > outerRatio {
		scale *= desiredRatio / outerRatio
	}
	return scale
}

type Scaler interface {
	GetNormalSizeAndScale(original r2.Vec) (r2.Vec, float64)
}

type hScaler struct{ scale float64 }

func (s hScaler) GetNormalSizeAndScale(original r2.Vec) (r2.Vec, float64) {
	ratio := original.X / original.Y
	scale := 1.0 / original.X
	return r2.Vec{X: s.scale, Y: s.scale / ratio}, scale
}

func HScaler(scale float64) Scaler {
	return hScaler{scale}
}

type vScaler struct{ scale float64 }

func (s vScaler) GetNormalSizeAndScale(original r2.Vec) (r2.Vec, float64) {
	ratio := original.X / original.Y
	scale := 1.0 / original.Y
	return r2.Vec{X: s.scale * ratio, Y: s.scale}, scale
}

func VScaler(scale float64) Scaler {
	return vScaler{scale}
}
