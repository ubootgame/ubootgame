package environment

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/do"
	"github.com/samber/lo"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	settings2 "github.com/ubootgame/ubootgame/framework/settings"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/components/graphics"
	"github.com/ubootgame/ubootgame/internal/layers"
	"github.com/yohamta/donburi/ecs"
)

type waterSystem struct {
	settings settings2.Provider[internal.Settings]

	sprite *graphics.SpriteData
}

func NewWaterSystem(i *do.Injector) ecsFramework.System {
	return &waterSystem{
		settings: do.MustInvoke[settings2.Provider[internal.Settings]](i),
	}
}

func (system *waterSystem) Layers() []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer]{
		{A: layers.Game, B: system.Draw},
	}
}

func (system *waterSystem) Update(_ *ecs.ECS) {}

func (system *waterSystem) Draw(_ *ecs.ECS, _ *ebiten.Image) {
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
