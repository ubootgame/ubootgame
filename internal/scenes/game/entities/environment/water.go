package environment

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/archetypes"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var WaterTag = donburi.NewTag().SetName("Background")

var Water = archetypes.NewArchetype(
	WaterTag,
	visuals.Sprite,
)

func CreateWater(ecs *ecs.ECS, registry *resources.Registry) *donburi.Entry {
	entry := Water.Spawn(ecs, layers.Background)

	sprite := registry.LoadTile(assets.Water, "fishTile_088.png")
	visuals.Sprite.SetValue(entry, visuals.SpriteData{Image: sprite.Data})

	return entry
}
