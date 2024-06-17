package entities

import (
	"github.com/jakecoffman/cp"
	"github.com/samber/do"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/internal/components/physics"
	"github.com/yohamta/donburi"
	"go/types"
)

var SpaceTag = donburi.NewTag().SetName("Space")

var Space = ecsFramework.NewArchetype(
	SpaceTag,
	physics.Space,
)

type NewSpaceParams types.Nil

var SpaceFactory ecsFramework.EntityFactory[NewSpaceParams] = func(i *do.Injector, params NewSpaceParams) *donburi.Entry {
	ecs := do.MustInvoke[ecsFramework.Service](i)

	entry := Space.Spawn(ecs.ECS())

	physics.Space.Set(entry, cp.NewSpace())

	return entry
}
