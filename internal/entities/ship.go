package entities

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/solarlune/resolv"
	"github.com/ubootgame/ubootgame/internal/assets"
	"github.com/ubootgame/ubootgame/internal/components"
	dresolv "github.com/ubootgame/ubootgame/internal/resolv"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var ShipTag = donburi.NewTag().SetName("Ship")

var Ship = utility.NewArchetype(
	ShipTag,
	components.Object,
	components.Sprite,
	components.Position,
	components.Velocity,
)

func CreateShip(ecs *ecs.ECS, resourceLoader *resource.Loader) *donburi.Entry {
	ship := Ship.Spawn(ecs)

	sprite := resourceLoader.LoadImage(assets.ImageBattleship)
	components.Sprite.SetValue(ship, components.SpriteData{Image: sprite.Data})

	obj := resolv.NewObject(0, 0, 64, 32)
	dresolv.SetObject(ship, obj)

	return ship
}
