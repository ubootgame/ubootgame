package environment

import (
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/framework/graphics/display"
	"github.com/ubootgame/ubootgame/framework/resources"
	"github.com/ubootgame/ubootgame/framework/resources/types"
	"github.com/ubootgame/ubootgame/internal/components/graphics"
	"github.com/ubootgame/ubootgame/internal/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
)

var WaterTag = donburi.NewTag().SetName("Game")

var Water = ecsFramework.NewArchetype(
	WaterTag,
	transform.Transform,
	graphics.Sprite,
)

func CreateWater(ecs *ecs.ECS, resources resources.Registry, tilesheetID types.TilesheetID, tileName string, scale display.Scale) *donburi.Entry {
	entry := Water.Spawn(ecs, layers.Game)

	//sprite := registry.LoadTile(tilesheetID, tileName)
	//
	//normalizedSize, normalizedScale, _ := scaler.GetNormalizedSizeAndScale(r2.Vec{X: float64(sprite.Data.Bounds().Size().X), Y: float64(sprite.Data.Bounds().Size().Y)})
	//
	//visuals.Sprite.SetValue(entry, visuals.SpriteData{
	//	Image:           sprite.Data,
	//	NormalizedSize:  normalizedSize,
	//	NormalizedScale: normalizedScale,
	//})

	return entry
}
