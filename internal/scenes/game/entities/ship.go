package entities

import (
	"github.com/solarlune/resolv"
	dresolv "github.com/ubootgame/ubootgame/internal/resolv"
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var ShipTag = donburi.NewTag().SetName("Ship")

var Ship = utility.NewArchetype(
	ShipTag,
	components.Object,
	components.Sprite,
	components.Velocity,
)

func CreateShip(ecs *ecs.ECS, registry *resources.Registry) *donburi.Entry {
	ship := Ship.Spawn(ecs)

	sprite := registry.LoadImage(assets.ImageBattleship)
	components.Sprite.SetValue(ship, components.SpriteData{Image: sprite.Data})

	obj := resolv.NewObject(0, 0, 64, 32)
	dresolv.SetObject(ship, obj)

	return ship
}
