package game_system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"gonum.org/v1/gonum/spatial/r2"
)

type CameraData struct {
	Position   r2.Vec
	ZoomFactor float64
	Rotation   float64
	Matrix     *ebiten.GeoM
}

var Camera = donburi.NewComponentType[CameraData](CameraData{
	ZoomFactor: 1.0,
	Matrix:     &ebiten.GeoM{},
})

func (camera *CameraData) WorldToScreenPosition(position r2.Vec) r2.Vec {
	x, y := camera.Matrix.Apply(position.X, position.Y)
	return r2.Vec{X: x, Y: y}
}

func (camera *CameraData) ScreenToWorldPosition(position r2.Vec) r2.Vec {
	matrix := *camera.Matrix
	matrix.Invert()
	x, y := matrix.Apply(position.X, position.Y)
	return r2.Vec{X: x, Y: y}
}
