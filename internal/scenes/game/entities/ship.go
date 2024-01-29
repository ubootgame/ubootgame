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

func CreateShip(ecs *ecs.ECS, registry *resources.Registry, scaler Scaler) *donburi.Entry {
	ship := Ship.Spawn(ecs, layers.Foreground)

	sprite := registry.LoadImage(assets.Battleship)

	size, scale := scaler.GetNormalSizeAndScale(r2.Vec{X: float64(sprite.Data.Bounds().Size().X), Y: float64(sprite.Data.Bounds().Size().Y)})

	components.Sprite.SetValue(ship, components.SpriteData{
		Image: sprite.Data,
		Scale: scale,
	})

	positionData := components.PositionData{
		Center: r2.Vec{},
		Size:   size,
	}
	components.Position.SetValue(ship, positionData)

	// TODO: Convert from world coordinates
	obj := resolv.NewObject(0, 0, 64, 32)
	dresolv.SetObject(ship, obj)

	return ship
}

type Scaler interface {
	GetNormalSizeAndScale(original r2.Vec) (r2.Vec, float64)
}

type hScaler struct{ scale float64 }

func (s hScaler) GetNormalSizeAndScale(original r2.Vec) (r2.Vec, float64) {
	ratio := original.X / original.Y
	scale := 1.0 / original.X
	return r2.Vec{X: s.scale, Y: s.scale / ratio}, scale
}

func HScaler(scale float64) Scaler {
	return hScaler{scale}
}

type vScaler struct{ scale float64 }

func (s vScaler) GetNormalSizeAndScale(original r2.Vec) (r2.Vec, float64) {
	ratio := original.X / original.Y
	scale := 1.0 / original.Y
	return r2.Vec{X: s.scale * ratio, Y: s.scale}, scale
}

func VScaler(scale float64) Scaler {
	return vScaler{scale}
}
