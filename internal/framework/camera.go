package framework

import (
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/framework/coordinate_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/yohamta/donburi/features/transform"
	"math"
)

func UpdateCameraMatrix(display *internal.Display, camera *game_system.CameraData, transform *transform.TransformData) {
	camera.Matrix.Reset()

	// Position camera in world coordinates
	camera.Matrix.Rotate(float64(transform.LocalRotation) * 2 * math.Pi / 360)
	camera.Matrix.Scale(transform.LocalScale.X, transform.LocalScale.Y)
	camera.Matrix.Translate(transform.LocalPosition.X, transform.LocalPosition.Y)

	// Move mid-point to center of screen
	camera.Matrix.Scale(coordinate_system.WorldSizeBase/display.VirtualResolution.X, coordinate_system.WorldSizeBase/display.VirtualResolution.X)
	camera.Matrix.Translate(display.VirtualResolution.X/2, display.VirtualResolution.Y/2)
}
