package game_system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"gonum.org/v1/gonum/spatial/r2"
)

type CameraData struct {
	Matrix *ebiten.GeoM
}

var Camera = donburi.NewComponentType[CameraData](CameraData{
	Matrix: &ebiten.GeoM{},
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
