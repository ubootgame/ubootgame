package actors

import (
	"github.com/solarlune/resolv"
	"github.com/ubootgame/ubootgame/framework"
	ecsFramework "github.com/ubootgame/ubootgame/framework/game/ecs"
	"github.com/ubootgame/ubootgame/framework/game/world"
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
	"gonum.org/v1/gonum/spatial/r2"
)

var EnemyTag = donburi.NewTag().SetName("Enemy")

var Enemy = ecsFramework.NewArchetype(
	EnemyTag,
	transform.Transform,
	geometry.Scale,
	geometry.Bounds,
	geometry.Direction,
	visuals.Sprite,
	geometry.Velocity,
)

func CreateEnemy(resources framework.ResourceService, ecs *ecs.ECS, scaler world.Scaler, position, velocity r2.Vec) *donburi.Entry {
	entry := Enemy.Spawn(ecs, layers.Game)

	sprite := resources.LoadImage(assets.Submarine)

	normalizedSize, normalizedScale, localScale := scaler.GetNormalizedSizeAndScale(r2.Vec{X: float64(sprite.Data.Bounds().Size().X), Y: float64(sprite.Data.Bounds().Size().Y)})

	visuals.Sprite.SetValue(entry, visuals.SpriteData{
		Image: sprite.Data,
	})
	geometry.Scale.SetValue(entry, geometry.ScaleData{
		NormalizedSize:  normalizedSize,
		NormalizedScale: normalizedScale,
	})
	transform.Transform.SetValue(entry, transform.TransformData{
		LocalPosition: math.Vec2(position),
		LocalScale:    math.NewVec2(localScale, localScale),
	})
	geometry.Velocity.SetValue(entry, velocity)

	var directionHorizontal geometry.DirectionHorizontal
	if velocity.X < 0 {
		directionHorizontal = geometry.Left
	} else {
		directionHorizontal = geometry.Right
	}

	geometry.Direction.SetValue(entry, geometry.DirectionData{
		HorizontalBase: geometry.Right,
		VerticalBase:   geometry.Up,
		Horizontal:     directionHorizontal,
		Vertical:       geometry.Up,
	})

	shape := resolv.NewRectangle(normalizedSize.X/2, normalizedSize.Y/2, normalizedSize.X, normalizedSize.Y)
	shape.RecenterPoints()
	geometry.Bounds.Set(entry, shape)

	return entry
}
