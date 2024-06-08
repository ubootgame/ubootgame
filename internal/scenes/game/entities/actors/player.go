package actors

import (
	"github.com/ubootgame/ubootgame/internal/framework"
	"github.com/ubootgame/ubootgame/internal/framework/ecs/archetypes"
	"github.com/ubootgame/ubootgame/internal/framework/resources"
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

var PlayerTag = donburi.NewTag().SetName("Player")

var Player = archetypes.NewArchetype(
	PlayerTag,
	transform.Transform,
	geometry.Scale,
	visuals.Sprite,
	geometry.Velocity,
	geometry.Direction,
)

func CreatePlayer(ecs *ecs.ECS, registry *resources.Registry, scaler framework.Scaler) *donburi.Entry {
	entry := Player.Spawn(ecs, layers.Game)

	sprite := registry.LoadImage(assets.Battleship)

	normalizedSize, normalizedScale, localScale := scaler.GetNormalizedSizeAndScale(r2.Vec{X: float64(sprite.Data.Bounds().Size().X), Y: float64(sprite.Data.Bounds().Size().Y)})

	visuals.Sprite.SetValue(entry, visuals.SpriteData{
		Image: sprite.Data,
	})
	geometry.Scale.SetValue(entry, geometry.ScaleData{
		NormalizedSize:  normalizedSize,
		NormalizedScale: normalizedScale,
	})
	transform.Transform.SetValue(entry, transform.TransformData{
		LocalScale: math.NewVec2(localScale, localScale),
	})
	geometry.Direction.SetValue(entry, geometry.DirectionData{
		HorizontalBase: geometry.Right,
		VerticalBase:   geometry.Up,
		Horizontal:     geometry.Right,
		Vertical:       geometry.Up,
	})

	return entry
}
