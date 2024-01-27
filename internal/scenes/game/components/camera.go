package components

import (
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/yohamta/donburi"
	"gonum.org/v1/gonum/spatial/r2"
)

type CameraData struct {
	Position   r2.Vec
	ZoomFactor float64
	Rotation   float64
}

var Camera = donburi.NewComponentType[CameraData](CameraData{
	Position:   r2.Vec{},
	ZoomFactor: 1.0,
	Rotation:   0.0,
})

func (c *CameraData) ViewportCenter() r2.Vec {
	return r2.Vec{
		X: config.C.VirtualResolution.X * 0.5,
		Y: config.C.VirtualResolution.Y * 0.5,
	}
}
