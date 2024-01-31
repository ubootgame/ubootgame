package components

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
