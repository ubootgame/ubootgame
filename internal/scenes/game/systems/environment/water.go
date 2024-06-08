package environment

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/framework/ecs/systems"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/yohamta/donburi/ecs"
)

type WaterSystem struct {
	systems.BaseSystem

	settings *internal.Settings

	sprite *visuals.SpriteData
}

func NewWaterSystem(settings *internal.Settings) *WaterSystem {
	system := &WaterSystem{settings: settings}
	return system
}

func (system *WaterSystem) Layers() []lo.Tuple2[ecs.LayerID, systems.Renderer] {
	return []lo.Tuple2[ecs.LayerID, systems.Renderer]{
		{A: layers.Game, B: system.Draw},
	}
}

func (system *WaterSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	//sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
	//w, h := float64(system.sprite.Image.Bounds().Size().X), float64(system.sprite.Image.Bounds().Dy())
	//y := sh / 2
	//
	//sizeScale := 0.1 * (sh / h)
	//
	//op := &ebiten.DrawImageOptions{}
	//
	//op.GeoM.Translate(float64(-w)/2, float64(-h)/2)
	//op.GeoM.Scale(sizeScale*system.display.ScalingFactor, sizeScale*system.display.ScalingFactor)
	//op.GeoM.Translate(0, y)
	//op.ColorScale.ScaleAlpha(0.1)
	//
	//for x := 0; x <= int(sw/(w*sizeScale)+1); x++ {
	//	screen.DrawImage(system.sprite.Image, op)
	//	op.GeoM.Translate(w*sizeScale, 0)
	//}
}
