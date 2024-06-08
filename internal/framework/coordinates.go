package framework

import (
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/yohamta/donburi/features/transform"
	"gonum.org/v1/gonum/spatial/r2"
	"math"
)

var WorldSizeBase = 1000.0

func UpdateCameraMatrix(display *internal.Display, camera *game_system.CameraData, transform *transform.TransformData) {
	camera.Matrix.Reset()

	// Position camera in world coordinates
	camera.Matrix.Rotate(float64(transform.LocalRotation) * 2 * math.Pi / 360)
	camera.Matrix.Scale(transform.LocalScale.X, transform.LocalScale.Y)
	camera.Matrix.Translate(transform.LocalPosition.X, transform.LocalPosition.Y)

	// Move mid-point to center of screen
	camera.Matrix.Scale(WorldSizeBase/display.VirtualResolution.X, WorldSizeBase/display.VirtualResolution.X)
	camera.Matrix.Translate(display.VirtualResolution.X/2, display.VirtualResolution.Y/2)
}

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
