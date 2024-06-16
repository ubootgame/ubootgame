package weapons

import (
	"github.com/jakecoffman/cp"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/internal/components/physics"
	"github.com/ubootgame/ubootgame/internal/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var BulletTag = donburi.NewTag().SetName("Bullet")

var Bullet = ecsFramework.NewArchetype(
	BulletTag,
	physics.Body,
)

func CreateBullet(ecs *ecs.ECS, from, to cp.Vector, space *cp.Space) *donburi.Entry {
	entry := Bullet.Spawn(ecs, layers.Game)

	direction := to.Sub(from)
	velocity := direction.Normalize()

	body := space.AddBody(cp.NewBody(1e8, cp.MomentForBox(1e8, 0.002, 0.002)))
	body.SetPosition(from)
	body.SetVelocityVector(velocity)
	space.AddShape(cp.NewBox(body, 0.02, 0.02, 0))
	physics.Body.Set(entry, body)

	return entry
}
