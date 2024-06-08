package environment

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/archetypes"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"gonum.org/v1/gonum/spatial/r2"
)

var WaterTag = donburi.NewTag().SetName("Game")

var Water = archetypes.NewArchetype(
	WaterTag,
	geometry.Scale,
	visuals.Sprite,
)

func CreateWater(ecs *ecs.ECS, registry *resources.Registry, scaler utility.Scaler) *donburi.Entry {
	entry := Water.Spawn(ecs, layers.Game)

	sprite := registry.LoadTile(assets.Water, "fishTile_088.png")
	visuals.Sprite.SetValue(entry, visuals.SpriteData{Image: sprite.Data})

	normalizedSize, normalizedScale, _ := scaler.GetNormalizedSizeAndScale(r2.Vec{X: float64(sprite.Data.Bounds().Size().X), Y: float64(sprite.Data.Bounds().Size().Y)})

	geometry.Scale.SetValue(entry, geometry.ScaleData{
		NormalizedSize:  normalizedSize,
		NormalizedScale: normalizedScale,
	})

	return entry
}
