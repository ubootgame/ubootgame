package geometry

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

var Bounds = donburi.NewComponentType[resolv.ConvexPolygon]()
