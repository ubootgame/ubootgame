package weapons

import (
	"github.com/jakecoffman/cp"
	"github.com/samber/do"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/internal/components/physics"
	"github.com/ubootgame/ubootgame/internal/layers"
	"github.com/yohamta/donburi"
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

var BulletFactory ecsFramework.EntityFactory[NewBulletParams] = func(i *do.Injector, params NewBulletParams) *donburi.Entry {
	ecs := do.MustInvoke[ecsFramework.Service](i)

	entry := Bullet.Spawn(ecs.ECS(), layers.Game)

	direction := params.To.Sub(params.From)
	velocity := direction.Normalize()

	body := params.Space.AddBody(cp.NewBody(1e8, cp.MomentForBox(1e8, 0.002, 0.002)))
	body.SetPosition(params.From)
	body.SetVelocityVector(velocity)
	params.Space.AddShape(cp.NewBox(body, 0.02, 0.02, 0))
	physics.Body.Set(entry, body)

	return entry
}
