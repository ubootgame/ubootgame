package scene_graph

import (
	"github.com/ubootgame/ubootgame/internal/framework/ecs/archetypes"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

var SceneGroupTag = donburi.NewTag().SetName("SceneGroup")

var SceneGroup = archetypes.NewArchetype(
	SceneGroupTag,
	transform.Transform,
)

func CreateSceneGroup(ecs *ecs.ECS, cts ...donburi.IComponentType) *donburi.Entry {
	entry := SceneGroup.Spawn(ecs, layers.Game, cts...)

	transform.Transform.SetValue(entry, transform.TransformData{
		LocalPosition: math.NewVec2(0, 0),
		LocalScale:    math.NewVec2(1, 1),
		LocalRotation: 0,
	})

	return entry
}
