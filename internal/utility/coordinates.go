package utility

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"math"
)

func CameraMatrix(camera *components.CameraData) ebiten.GeoM {
	m := ebiten.GeoM{}
	m.Translate(-camera.Position.X*config.C.VirtualResolution.X, -camera.Position.Y*config.C.VirtualResolution.Y)
	m.Translate(-camera.ViewportCenter().X, -camera.ViewportCenter().Y)
	m.Scale(camera.ZoomFactor, camera.ZoomFactor)
	m.Rotate(float64(camera.Rotation) * 2 * math.Pi / 360)
	m.Translate(camera.ViewportCenter().X, camera.ViewportCenter().Y)
	return m
}

func CalculateScalingFactor() float64 {
	desiredRatio := config.C.VirtualResolution.X / config.C.VirtualResolution.Y
	outerRatio := config.C.ActualOuterSize.X / config.C.ActualOuterSize.Y
	scale := config.C.VirtualResolution.Y / config.C.ActualOuterSize.Y
	if desiredRatio > outerRatio {
		scale *= desiredRatio / outerRatio
	}
	return scale
}
