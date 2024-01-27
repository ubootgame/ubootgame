package components

import (
	"github.com/yohamta/donburi"
	"gonum.org/v1/gonum/spatial/r2"
)

type PositionData struct {
	Center         r2.Vec
	Scale          float64
	ScaleDirection ScaleDirection
}

var Position = donburi.NewComponentType[PositionData]()

type ScaleDirection int

const (
	Horizontal ScaleDirection = iota
	Vertical
)
