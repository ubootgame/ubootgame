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

var PlayerTag = donburi.NewTag().SetName("Player")

var Player = utility.NewArchetype(
	PlayerTag,
	components.Object,
	components.Sprite,
	components.Transform,
	components.Velocity,
)

func CreatePlayer(ecs *ecs.ECS, registry *resources.Registry, scaler utility.Scaler) *donburi.Entry {
	entry := Player.Spawn(ecs, layers.Foreground)

	sprite := registry.LoadImage(assets.Battleship)

	size, scale := scaler.GetNormalSizeAndScale(r2.Vec{X: float64(sprite.Data.Bounds().Size().X), Y: float64(sprite.Data.Bounds().Size().Y)})

	components.Sprite.SetValue(entry, components.SpriteData{
		Image: sprite.Data,
		Scale: scale,
	})
	components.Transform.SetValue(entry, components.TransformData{
		Center: r2.Vec{},
		Size:   size,
	})

	// TODO: Convert from world coordinates
	obj := resolv.NewObject(0, 0, 64, 32)
	dresolv.SetObject(entry, obj)

	return entry
}
