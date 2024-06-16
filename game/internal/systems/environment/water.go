package environment

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/framework"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/components/graphics"
	"github.com/ubootgame/ubootgame/internal/layers"
	"github.com/yohamta/donburi/ecs"
)

type WaterSystem struct {
	settings framework.SettingsService[internal.Settings]

	sprite *graphics.SpriteData
}

func NewWaterSystem(settings framework.SettingsService[internal.Settings]) *WaterSystem {
	system := &WaterSystem{settings: settings}
	return system
}

func (system *WaterSystem) Layers() []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer]{
		{A: layers.Game, B: system.Draw},
	}
}

func (system *WaterSystem) Update(_ *ecs.ECS) {}

func (system *WaterSystem) Draw(_ *ecs.ECS, _ *ebiten.Image) {
	//sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
	//w, h := float64(system.sprite.Image.Bounds().Size().X), float64(system.sprite.Image.Bounds().Dy())
	//y := sh / 2
	//
	//sizeScale := 0.1 * (sh / h)
	//
	//op := &ebiten.DrawImageOptions{}
	//
	//op.GeoM.Translate(float64(-w)/2, float64(-h)/2)
	//op.GeoM.Zoom(sizeScale*system.display.ScalingFactor, sizeScale*system.display.ScalingFactor)
	//op.GeoM.Translate(0, y)
	//op.ColorScale.ScaleAlpha(0.1)
	//
	//for x := 0; x <= int(sw/(w*sizeScale)+1); x++ {
	//	screen.DrawImage(system.sprite.Image, op)
	//	op.GeoM.Translate(w*sizeScale, 0)
	//}
}
