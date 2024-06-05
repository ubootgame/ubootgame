package geometry

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

var Shape = donburi.NewComponentType[resolv.ConvexPolygon]()
