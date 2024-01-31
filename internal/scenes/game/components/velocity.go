package components

import (
	"github.com/yohamta/donburi"
	"gonum.org/v1/gonum/spatial/r2"
)

var Velocity = donburi.NewComponentType[r2.Vec]()
