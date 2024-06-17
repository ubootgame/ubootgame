package weapons

import (
	"github.com/jakecoffman/cp"
	"github.com/samber/do"
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

type NewBulletParams struct {
	From, To cp.Vector
	Space    *cp.Space
}

func CreateBullet(_ *do.Injector, e *ecs.ECS, params NewBulletParams) *donburi.Entry {
	entry := Bullet.Spawn(e, layers.Game)

	direction := params.To.Sub(params.From)
	velocity := direction.Normalize()

	body := params.Space.AddBody(cp.NewBody(1e8, cp.MomentForBox(1e8, 0.002, 0.002)))
	body.SetPosition(params.From)
	body.SetVelocityVector(velocity)
	params.Space.AddShape(cp.NewBox(body, 0.02, 0.02, 0))
	physics.Body.Set(entry, body)

	return entry
}
