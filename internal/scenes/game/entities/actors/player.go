package actors

import (
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

var PlayerTag = donburi.NewTag().SetName("Player")

var Player = archetypes.NewArchetype(
	PlayerTag,
	visuals.Sprite,
	geometry.Transform,
	geometry.Velocity,
)

func CreatePlayer(ecs *ecs.ECS, registry *resources.Registry, scaler utility.Scaler) *donburi.Entry {
	entry := Player.Spawn(ecs, layers.Foreground)

	sprite := registry.LoadImage(assets.Battleship)

	size, scale := scaler.GetNormalSizeAndScale(r2.Vec{X: float64(sprite.Data.Bounds().Size().X), Y: float64(sprite.Data.Bounds().Size().Y)})

	visuals.Sprite.SetValue(entry, visuals.SpriteData{
		Image: sprite.Data,
		Scale: scale,
	})
	geometry.Transform.SetValue(entry, geometry.TransformData{
		Center: r2.Vec{},
		Size:   size,
	})

	return entry
}
