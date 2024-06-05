package utility

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"gonum.org/v1/gonum/spatial/r2"
	"math"
)

func UpdateCameraMatrix(display *game_system.DisplayData, camera *game_system.CameraData) {
	camera.Matrix.Reset()
	camera.Matrix.Translate(-(camera.Position.X), -(camera.Position.Y))
	camera.Matrix.Scale(camera.ZoomFactor, camera.ZoomFactor)
	camera.Matrix.Rotate(float64(camera.Rotation) * 2 * math.Pi / 360)
	camera.Matrix.Translate(0.5, 0.5/display.Ratio())
	camera.Matrix.Scale(display.VirtualResolution.X, display.VirtualResolution.X)
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
