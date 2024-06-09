package framework

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/framework/coordinate_system"
	"gonum.org/v1/gonum/spatial/r2"
	"math"
)

const translationSpeed, zoomSpeed = 500.0, 0.1 // world unit
const rotationSpeed = 2                        // degrees
const minZoom, maxZoom = 0.5, 2.0

type Camera struct {
	settings *internal.Settings

	matrix ebiten.GeoM

	Position r2.Vec
	Rotation float64
	Scale    r2.Vec
}

func NewCamera(settings *internal.Settings) *Camera {
	return &Camera{
		settings: settings,
		matrix:   ebiten.GeoM{},
		Scale:    r2.Vec{X: 1, Y: 1},
	}
}

func (camera *Camera) WorldToScreenPosition(position r2.Vec) r2.Vec {
	x, y := camera.matrix.Apply(position.X, position.Y)
	return r2.Vec{X: x, Y: y}
}

func (camera *Camera) ScreenToWorldPosition(position r2.Vec) r2.Vec {
	camera.matrix.Invert()
	x, y := camera.matrix.Apply(position.X, position.Y)
	return r2.Vec{X: x, Y: y}
}

func (camera *Camera) Apply(geom *ebiten.GeoM) {
	geom.Concat(camera.matrix)
}

func (camera *Camera) UpdateCameraMatrix() {
	camera.matrix.Reset()

	// Position camera in world coordinates
	camera.matrix.Rotate(float64(camera.Rotation) * 2 * math.Pi / 360)
	camera.matrix.Scale(camera.Scale.X, camera.Scale.Y)
	camera.matrix.Translate(camera.Position.X, camera.Position.Y)

	// Move mid-point to center of screen
	camera.matrix.Scale(coordinate_system.WorldSizeBase/camera.settings.Display.VirtualResolution.X, coordinate_system.WorldSizeBase/camera.settings.Display.VirtualResolution.X)
	camera.matrix.Translate(camera.settings.Display.VirtualResolution.X/2, camera.settings.Display.VirtualResolution.Y/2)
}
