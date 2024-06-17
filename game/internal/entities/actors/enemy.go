package actors

import (
	"github.com/jakecoffman/cp"
	"github.com/samber/do"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/framework/graphics/display"
	"github.com/ubootgame/ubootgame/framework/resources"
	"github.com/ubootgame/ubootgame/framework/resources/types"
	"github.com/ubootgame/ubootgame/internal/components/graphics"
	"github.com/ubootgame/ubootgame/internal/components/physics"
	"github.com/ubootgame/ubootgame/internal/layers"
	"github.com/yohamta/donburi"
)

var EnemyTag = donburi.NewTag().SetName("Enemy")

var Enemy = ecsFramework.NewArchetype(
	EnemyTag,
	graphics.Sprite,
	physics.Body,
)

type NewEnemyParams struct {
	ImageID            types.ImageID
	Scale              display.Scale
	Position, Velocity cp.Vector
	Space              *cp.Space
}

var EnemyFactory ecsFramework.EntityFactory[NewEnemyParams] = func(i *do.Injector, params NewEnemyParams) *donburi.Entry {
	resourceRegistry := do.MustInvoke[resources.Registry](i)
	ecs := do.MustInvoke[ecsFramework.Service](i)

	entry := Enemy.Spawn(ecs.ECS(), layers.Game)

	image := resourceRegistry.LoadImage(params.ImageID)

	sprite := graphics.NewSprite(image, params.Scale, false, false)
	graphics.Sprite.SetValue(entry, sprite)

	worldSize := sprite.WorldSize()

	body := params.Space.AddBody(cp.NewBody(1e9, cp.MomentForBox(1e9, worldSize.X, worldSize.Y)))
	body.SetPosition(params.Position)
	body.SetVelocityVector(params.Velocity)
	params.Space.AddShape(cp.NewBox(body, worldSize.X, worldSize.Y, 0))
	physics.Body.Set(entry, body)

	return entry
}
