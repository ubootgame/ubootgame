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
	"github.com/ubootgame/ubootgame/internal/entities"
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
}

var EnemyFactory ecsFramework.EntityFactory[NewEnemyParams] = func(i *do.Injector, params NewEnemyParams) *donburi.Entry {
	resourceRegistry := do.MustInvoke[resources.Registry](i)
	ecs := do.MustInvoke[ecsFramework.Service](i)

	entry := Enemy.SpawnOnLayer(ecs.ECS(), layers.Game)

	image := resourceRegistry.LoadImage(params.ImageID)

	sprite := graphics.NewSprite(image, params.Scale, false, false)
	graphics.Sprite.SetValue(entry, sprite)

	worldSize := sprite.WorldSize()

	spaceEntry, _ := entities.SpaceTag.First(ecs.World())
	space := physics.Space.Get(spaceEntry)

	body := space.AddBody(cp.NewBody(1e9, cp.MomentForBox(1e9, worldSize.X, worldSize.Y)))
	body.SetPosition(params.Position)
	body.SetVelocityVector(params.Velocity)
	space.AddShape(cp.NewBox(body, worldSize.X, worldSize.Y, 0))
	physics.Body.Set(entry, body)

	return entry
}
