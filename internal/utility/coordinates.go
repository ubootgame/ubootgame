package utility

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"gonum.org/v1/gonum/spatial/r2"
	"math"
)

func CameraMatrix(camera *components.CameraData) ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-camera.Position.X, -camera.Position.Y)
	m.Translate(-camera.ViewportCenter().X, -camera.ViewportCenter().Y)
	m.Scale(camera.ZoomFactor, camera.ZoomFactor)
	m.Rotate(float64(camera.Rotation) * 2 * math.Pi / 360)
	m.Translate(camera.ViewportCenter().X, camera.ViewportCenter().Y)
	return m
}

func PositionMatrix(position, actualSize r2.Vec, scale float64, direction components.ScaleDirection) ebiten.GeoM {
	sw, sh := config.C.VirtualResolution.X, config.C.VirtualResolution.Y

	sizeScale := ebiten.DeviceScaleFactor()
	switch direction {
	case components.Horizontal:
		sizeScale *= scale * (sw / actualSize.X)
	case components.Vertical:
		sizeScale *= scale * (sh / actualSize.Y)
	}

	m := ebiten.GeoM{}

	m.Translate(float64(-actualSize.X)/2, float64(-actualSize.Y)/2)
	m.Scale(sizeScale, sizeScale)
	m.Translate(float64(sw)/2+position.X, float64(sh)/2+position.Y)

	return m
}
