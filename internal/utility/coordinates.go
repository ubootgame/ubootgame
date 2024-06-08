package utility

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/yohamta/donburi/features/transform"
	"gonum.org/v1/gonum/spatial/r2"
	"math"
)

func UpdateCameraMatrix(display *game_system.DisplayData, camera *game_system.CameraData, transform *transform.TransformData) {
	camera.Matrix.Reset()

	// Position camera in world coordinates
	camera.Matrix.Rotate(float64(transform.LocalRotation) * 2 * math.Pi / 360)
	camera.Matrix.Scale(transform.LocalScale.X, transform.LocalScale.Y)
	camera.Matrix.Translate(transform.LocalPosition.X, transform.LocalPosition.Y)

	// Move mid-point to center of screen
	camera.Matrix.Scale(display.VirtualResolution.X, display.VirtualResolution.X)
	camera.Matrix.Translate(display.VirtualResolution.X/2, display.VirtualResolution.Y/2)
}

type Scaler interface {
	GetNormalSizeAndScale(size r2.Vec) (r2.Vec, float64)
	GetNormalizedScale(size r2.Vec) float64
}

type hScaler struct{ scale float64 }

func (s hScaler) GetNormalSizeAndScale(size r2.Vec) (r2.Vec, float64) {
	ratio := size.X / size.Y
	scale := 1.0 / size.X
	return r2.Vec{X: s.scale * scale, Y: s.scale / ratio}, scale
}

func (s hScaler) GetNormalizedScale(size r2.Vec) float64 {
	return (1.0 / size.X) * s.scale
}

func HScaler(scale float64) Scaler {
	return hScaler{scale}
}

type vScaler struct{ scale float64 }

func (s vScaler) GetNormalSizeAndScale(size r2.Vec) (r2.Vec, float64) {
	ratio := size.X / size.Y
	scale := 1.0 / size.Y
	return r2.Vec{X: s.scale / ratio, Y: s.scale}, scale
}

func (s vScaler) GetNormalizedScale(size r2.Vec) float64 {
	return (1.0 / size.Y) * s.scale
}

func VScaler(scale float64) Scaler {
	return vScaler{scale}
}
