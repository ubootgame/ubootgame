package entities

import (
	"github.com/solarlune/resolv"
	dresolv "github.com/ubootgame/ubootgame/internal/resolv"
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"gonum.org/v1/gonum/spatial/r2"
)

var ShipTag = donburi.NewTag().SetName("Ship")

var Ship = utility.NewArchetype(
	ShipTag,
	components.Object,
	components.Sprite,
	components.Position,
	components.Velocity,
)

func CreateShip(ecs *ecs.ECS, registry *resources.Registry) *donburi.Entry {
	ship := Ship.Spawn(ecs, layers.Foreground)

	sprite := registry.LoadImage(assets.Battleship)
	components.Sprite.SetValue(ship, components.SpriteData{Image: sprite.Data})

	positionData := components.PositionData{
		Center:         r2.Vec{},
		Scale:          0.1,
		ScaleDirection: components.Horizontal,
	}
	components.Position.SetValue(ship, positionData)

	obj := resolv.NewObject(0, 0, 64, 32)
	dresolv.SetObject(ship, obj)

	return ship
}
