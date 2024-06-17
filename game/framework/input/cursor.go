package input

import "gonum.org/v1/gonum/spatial/r2"

type Cursor struct {
	ScreenPosition, WorldPosition r2.Vec
}
