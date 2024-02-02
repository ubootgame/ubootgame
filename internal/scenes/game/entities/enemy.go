package entities

import (
	"github.com/solarlune/resolv"
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
	components.Shape,
	components.Sprite,
	components.Transform,
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
	components.Transform.SetValue(entry, components.TransformData{
		Center: position,
		Size:   size,
	})
	components.Velocity.SetValue(entry, velocity)

	shapePosition := r2.Vec{X: position.X - size.X/2, Y: position.Y - size.Y/2}
	components.Shape.SetValue(entry, *resolv.NewRectangle(shapePosition.X, shapePosition.Y, size.X, size.Y))

	return entry
}
