package actors

import (
	"github.com/jakecoffman/cp"
	"github.com/ubootgame/ubootgame/framework"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/framework/services/display"
	"github.com/ubootgame/ubootgame/framework/services/resources"
	"github.com/ubootgame/ubootgame/internal/components/graphics"
	"github.com/ubootgame/ubootgame/internal/components/physics"
	"github.com/ubootgame/ubootgame/internal/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var EnemyTag = donburi.NewTag().SetName("Enemy")

var Enemy = ecsFramework.NewArchetype(
	EnemyTag,
	graphics.Sprite,
	physics.Body,
)

func CreateEnemy(resources framework.ResourceService, ecs *ecs.ECS, imageID resources.ImageID, scale display.Scale, position, velocity cp.Vector, space *cp.Space) *donburi.Entry {
	entry := Enemy.Spawn(ecs, layers.Game)

	image := resources.LoadImage(imageID)

	sprite := graphics.NewSprite(image, scale, false, false)
	graphics.Sprite.SetValue(entry, sprite)

	worldSize := sprite.WorldSize()

	body := space.AddBody(cp.NewBody(1e9, cp.MomentForBox(1e9, worldSize.X, worldSize.Y)))
	body.SetPosition(position)
	body.SetVelocityVector(velocity)
	space.AddShape(cp.NewBox(body, worldSize.X, worldSize.Y, 0))
	physics.Body.Set(entry, body)

	return entry
}
