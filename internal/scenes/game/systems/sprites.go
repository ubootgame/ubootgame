package systems

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func DrawSprites(e *ecs.ECS, screen *ebiten.Image) {
	ecs.NewQuery(layers.Foreground, filter.Contains(components.Sprite, components.Position)).Each(e.World, func(entry *donburi.Entry) {
		spriteData := components.Sprite.Get(entry)
		positionData := components.Position.Get(entry)

		cameraEntry, _ := components.Camera.First(e.World)
		camera := components.Camera.Get(cameraEntry)

		sw, sh := config.C.VirtualResolution.X, config.C.VirtualResolution.Y
		w, h := float64(spriteData.Image.Bounds().Dx()), float64(spriteData.Image.Bounds().Dy())

		sizeScale := ebiten.DeviceScaleFactor()
		switch positionData.ScaleDirection {
		case components.Horizontal:
			sizeScale *= positionData.Scale * (sw / w)
		case components.Vertical:
			sizeScale *= positionData.Scale * (sh / h)
		}

		op := &ebiten.DrawImageOptions{}

		op.GeoM.Translate(-w/2, -h/2)
		op.GeoM.Scale(sizeScale, sizeScale)
		op.GeoM.Translate(float64(sw)/2+(positionData.Center.X*sw), float64(sh)/2+(positionData.Center.Y*sh))
		op.GeoM.Concat(utility.CameraMatrix(camera))

		op.Filter = ebiten.FilterLinear

		screen.DrawImage(spriteData.Image, op)

		if config.C.Debug {
			debugText := fmt.Sprintf("Position: %.2f, %.2f\nScale: %.2f",
				positionData.Center.X, positionData.Center.Y,
				positionData.Scale)

			if entry.HasComponent(components.Velocity) {
				velocityData := components.Velocity.Get(entry)
				debugText += fmt.Sprintf("\nVelocity: %.3f, %.3f", velocityData.X, velocityData.Y)
			}

			geom := ebiten.GeoM{}

			geom.Translate(-config.C.VirtualResolution.X/8, -config.C.VirtualResolution.Y/8)
			geom.Translate(float64(sw)/2+(positionData.Center.X*sw), float64(sh)/2+(positionData.Center.Y*sh))
			geom.Concat(utility.CameraMatrix(camera))

			printDebugTextAt(screen, debugText, geom)
		}
	})
}
