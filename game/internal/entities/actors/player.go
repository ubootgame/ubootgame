package actors

import (
	"github.com/jakecoffman/cp"
	"github.com/ubootgame/ubootgame/framework"
	ecsFramework "github.com/ubootgame/ubootgame/framework/game/ecs"
	"github.com/ubootgame/ubootgame/framework/services/display"
	"github.com/ubootgame/ubootgame/framework/services/resources"
	"github.com/ubootgame/ubootgame/internal/components/graphics"
	"github.com/ubootgame/ubootgame/internal/components/physics"
	"github.com/ubootgame/ubootgame/internal/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var PlayerTag = donburi.NewTag().SetName("Player")

var Player = ecsFramework.NewArchetype(
	PlayerTag,
	graphics.Sprite,
	physics.Body,
)

func CreatePlayer(resources framework.ResourceService, ecs *ecs.ECS, imageID resources.ImageID, scale display.Scale, space *cp.Space) *donburi.Entry {
	entry := Player.Spawn(ecs, layers.Game)

	image := resources.LoadImage(imageID)

	sprite := graphics.NewSprite(image, scale, false, false)
	graphics.Sprite.SetValue(entry, sprite)

	worldSize := sprite.WorldSize()

	body := space.AddBody(cp.NewBody(1e9, cp.INFINITY))
	space.AddShape(cp.NewBox(body, worldSize.X, worldSize.Y, 0))
	physics.Body.Set(entry, body)

	return entry
}
