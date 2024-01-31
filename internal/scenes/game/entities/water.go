package entities

import (
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var WaterTag = donburi.NewTag().SetName("Background")

var Water = utility.NewArchetype(
	WaterTag,
	components.Sprite,
)

func CreateWater(ecs *ecs.ECS, registry *resources.Registry) *donburi.Entry {
	entry := Water.Spawn(ecs, layers.Background)

	sprite := registry.LoadTile(assets.Water, "fishTile_088.png")
	components.Sprite.SetValue(entry, components.SpriteData{Image: sprite.Data})

	return entry
}
