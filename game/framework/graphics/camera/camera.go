package camera

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/framework"
	"github.com/ubootgame/ubootgame/framework/game/world"
	"gonum.org/v1/gonum/spatial/r2"
	"math"
)

type Camera struct {
	display framework.DisplayService

	matrix ebiten.GeoM

	Position r2.Vec
	Zoom     r2.Vec
	Rotation float64
}

func NewCamera(display framework.DisplayService) *Camera {
	return &Camera{
		display: display,
		matrix:  ebiten.GeoM{},
		Zoom:    r2.Vec{X: 1, Y: 1},
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
	camera.matrix.Scale(camera.Zoom.X, camera.Zoom.Y)
	camera.matrix.Translate(camera.Position.X, camera.Position.Y)

	// Move mid-point to center of screen
	vx, vy := camera.display.VirtualResolution()
	camera.matrix.Scale(vx/world.BaseSize, vx/world.BaseSize)
	camera.matrix.Translate(vx/2, vy/2)
}