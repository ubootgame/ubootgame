package game_system

import (
	"github.com/ubootgame/ubootgame/internal/framework/ecs/archetypes"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
)

var CameraTag = donburi.NewTag().SetName("Camera")

var Camera = archetypes.NewArchetype(
	CameraTag,
	transform.Transform,
	game_system.Camera,
)

func CreateCamera(ecs *ecs.ECS) *donburi.Entry {
	entry := Camera.Spawn(ecs, layers.Game)

	return entry
}
