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

var PlayerTag = donburi.NewTag().SetName("Player")

var Player = ecsFramework.NewArchetype(
	PlayerTag,
	graphics.Sprite,
	physics.Body,
)

type NewPlayerParams struct {
	ImageID types.ImageID
	Scale   display.Scale
	Space   *cp.Space
}

var PlayerFactory ecsFramework.EntityFactory[NewPlayerParams] = func(i *do.Injector, params NewPlayerParams) *donburi.Entry {
	resourceRegistry := do.MustInvoke[resources.Registry](i)
	ecs := do.MustInvoke[ecsFramework.Service](i)

	entry := Player.Spawn(ecs.ECS(), layers.Game)

	image := resourceRegistry.LoadImage(params.ImageID)

	sprite := graphics.NewSprite(image, params.Scale, false, false)
	graphics.Sprite.SetValue(entry, sprite)

	worldSize := sprite.WorldSize()

	body := params.Space.AddBody(cp.NewBody(1e9, cp.INFINITY))
	params.Space.AddShape(cp.NewBox(body, worldSize.X, worldSize.Y, 0))
	physics.Body.Set(entry, body)

	return entry
}