package camera

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/framework"
	"gonum.org/v1/gonum/spatial/r2"
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
	matrix := &camera.matrix
	matrix.Invert()
	x, y := matrix.Apply(position.X, position.Y)
	return r2.Vec{X: x, Y: y}
}

func (camera *Camera) Apply(geom *ebiten.GeoM) {
	geom.Concat(camera.matrix)
}
