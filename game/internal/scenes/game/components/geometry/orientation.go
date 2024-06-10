package geometry

import (
	"github.com/yohamta/donburi"
)

type DirectionData struct {
	HorizontalBase DirectionHorizontal
	VerticalBase   DirectionVertical
	Horizontal     DirectionHorizontal
	Vertical       DirectionVertical
}

var Direction = donburi.NewComponentType[DirectionData]()

type DirectionHorizontal int

const (
	Left DirectionHorizontal = iota
	Right
)

type DirectionVertical int

const (
	Up DirectionVertical = iota
	Down
)
