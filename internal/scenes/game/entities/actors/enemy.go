package actors

import (
	"github.com/solarlune/resolv"
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

var EnemyTag = donburi.NewTag().SetName("Enemy")

var Enemy = archetypes.NewArchetype(
	EnemyTag,
	geometry.Shape,
	visuals.Sprite,
	geometry.Transform,
	geometry.Velocity,
)

func CreateEnemy(ecs *ecs.ECS, registry *resources.Registry, scaler utility.Scaler, position, velocity r2.Vec) *donburi.Entry {
	entry := Enemy.Spawn(ecs, layers.Foreground)

	sprite := registry.LoadImage(assets.Submarine)

	size, scale := scaler.GetNormalSizeAndScale(r2.Vec{X: float64(sprite.Data.Bounds().Size().X), Y: float64(sprite.Data.Bounds().Size().Y)})

	visuals.Sprite.SetValue(entry, visuals.SpriteData{
		Image: sprite.Data,
		Scale: scale,
	})
	geometry.Transform.SetValue(entry, geometry.TransformData{
		Center: position,
		Size:   size,
	})
	geometry.Velocity.SetValue(entry, velocity)

	shapePosition := r2.Vec{X: position.X - size.X/2, Y: position.Y - size.Y/2}
	geometry.Shape.SetValue(entry, *resolv.NewRectangle(shapePosition.X, shapePosition.Y, size.X, size.Y))

	return entry
}
