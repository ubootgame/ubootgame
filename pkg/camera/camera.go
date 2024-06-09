package camera

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/pkg/game"
	"github.com/ubootgame/ubootgame/pkg/world"
	"gonum.org/v1/gonum/spatial/r2"
	"math"
)

type Camera struct {
	displayInfo *game.DisplayInfo

	matrix ebiten.GeoM

	Position r2.Vec
	Scale    r2.Vec
	Rotation float64
}

func NewCamera(displayInfo *game.DisplayInfo) *Camera {
	return &Camera{
		displayInfo: displayInfo,
		matrix:      ebiten.GeoM{},
		Scale:       r2.Vec{X: 1, Y: 1},
	}
}

func (camera *Camera) WorldToScreenPosition(position r2.Vec) r2.Vec {
	x, y := camera.matrix.Apply(position.X, position.Y)
	return r2.Vec{X: x, Y: y}
}

func (camera *Camera) ScreenToWorldPosition(position r2.Vec) r2.Vec {
	camera.matrix.Invert()
	x, y := camera.matrix.Apply(position.X, position.Y)
	return r2.Vec{X: x, Y: y}
}

func (camera *Camera) Apply(geom *ebiten.GeoM) {
	geom.Concat(camera.matrix)
}

func (camera *Camera) UpdateCameraMatrix() {
	camera.matrix.Reset()

	// Position camera in world coordinates
	camera.matrix.Rotate(float64(camera.Rotation) * 2 * math.Pi / 360)
	camera.matrix.Scale(camera.Scale.X, camera.Scale.Y)
	camera.matrix.Translate(camera.Position.X, camera.Position.Y)

	// Move mid-point to center of screen
	camera.matrix.Scale(world.WorldSizeBase/camera.displayInfo.VirtualResolution.X, world.WorldSizeBase/camera.displayInfo.VirtualResolution.X)
	camera.matrix.Translate(camera.displayInfo.VirtualResolution.X/2, camera.displayInfo.VirtualResolution.Y/2)
}
