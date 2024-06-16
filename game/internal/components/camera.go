package components

import (
	"github.com/setanarut/kamera/v2"
	"github.com/yohamta/donburi"
	"gonum.org/v1/gonum/spatial/r2"
)

type CameraData struct {
	Camera                              *kamera.Camera
	MoveSpeed, ZoomSpeed, RotationSpeed float64
	MinZoom, MaxZoom                    float64
	Target, Delta                       r2.Vec
}

var Camera = donburi.NewComponentType[CameraData]()
