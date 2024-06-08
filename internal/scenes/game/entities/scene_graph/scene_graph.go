package scene_graph

import (
	ecs2 "github.com/ubootgame/ubootgame/internal/framework/ecs"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

var SceneGraphTag = donburi.NewTag().SetName("SceneGraph")

var SceneGraph = ecs2.NewArchetype(
	SceneGraphTag,
	transform.Transform,
)

func CreateSceneGraph(ecs *ecs.ECS) *donburi.Entry {
	entry := SceneGraph.Spawn(ecs, layers.Game)

	transform.Transform.SetValue(entry, transform.TransformData{
		LocalPosition: math.NewVec2(0, 0),
		LocalScale:    math.NewVec2(1, 1),
		LocalRotation: 0,
	})

	return entry
}
