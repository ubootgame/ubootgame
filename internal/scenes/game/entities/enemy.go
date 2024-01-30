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

var EnemyTag = donburi.NewTag().SetName("Enemy")

var Enemy = utility.NewArchetype(
	EnemyTag,
	components.Object,
	components.Sprite,
	components.Position,
	components.Velocity,
)

func CreateEnemy(ecs *ecs.ECS, registry *resources.Registry, scaler utility.Scaler, position, velocity r2.Vec) *donburi.Entry {
	entry := Enemy.Spawn(ecs, layers.Foreground)

	sprite := registry.LoadImage(assets.Submarine)

	size, scale := scaler.GetNormalSizeAndScale(r2.Vec{X: float64(sprite.Data.Bounds().Size().X), Y: float64(sprite.Data.Bounds().Size().Y)})

	components.Sprite.SetValue(entry, components.SpriteData{
		Image: sprite.Data,
		Scale: scale,
	})
	components.Position.SetValue(entry, components.PositionData{
		Center: position,
		Size:   size,
	})
	components.Velocity.SetValue(entry, velocity)

	// TODO: Convert from world coordinates
	obj := resolv.NewObject(0, 0, 64, 32)
	dresolv.SetObject(entry, obj)

	return entry
}
