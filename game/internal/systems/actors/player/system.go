package player

import (
	"github.com/jakecoffman/cp"
	"github.com/samber/do"
	"github.com/samber/lo"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/framework/input"
	"github.com/ubootgame/ubootgame/framework/settings"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/components/graphics"
	"github.com/ubootgame/ubootgame/internal/components/physics"
	"github.com/ubootgame/ubootgame/internal/entities/actors"
	"github.com/ubootgame/ubootgame/internal/entities/weapons"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"go/types"
)

const (
	acceleration = 0.01
	friction     = 0.05
	maxSpeed     = 0.25
)

type playerSystem struct {
	injector *do.Injector

	settingsProvider settings.Provider[internal.Settings]
	input            input.Input

	space  *cp.Space
	body   *cp.Body
	sprite *graphics.SpriteData

	moving, shooting bool
	fireTick         uint64
}

func NewPlayerSystem(i *do.Injector) ecsFramework.System {
	system := &playerSystem{
		injector:         i,
		settingsProvider: do.MustInvoke[settings.Provider[internal.Settings]](i),
		input:            do.MustInvoke[input.Input](i),
	}

	e := do.MustInvoke[ecsFramework.Service](i)

	MoveLeftEvent.Subscribe(e.World(), system.MoveLeft)
	MoveRightEvent.Subscribe(e.World(), system.MoveRight)
	ShootEvent.Subscribe(e.World(), system.Shoot)

	return system
}

func (system *playerSystem) Layers() []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer]{}
}

func (system *playerSystem) Update(e *ecs.ECS) {
	if entry, found := actors.PlayerTag.First(e.World); found {
		system.body = physics.Body.Get(entry)
		system.sprite = graphics.Sprite.Get(entry)
	}
	if entry, found := physics.Space.First(e.World); found {
		system.space = physics.Space.Get(entry)
	}

	if system.moving {
		system.moving = false
	} else {
		system.body.SetVelocityVector(system.body.Velocity().Mult(1 - friction))
	}

	if system.shooting {
		system.shooting = false
	} else {
		system.fireTick = 0
	}

	if system.body.Velocity().X < 0 {
		system.sprite.FlipY = true
	} else if system.body.Velocity().X > 0 {
		system.sprite.FlipY = false
	}
}

func (system *playerSystem) MoveLeft(_ donburi.World, _ types.Nil) {
	velocity := system.body.Velocity()
	velocity.X -= acceleration
	velocity.X = max(velocity.X, -maxSpeed)
	system.body.SetVelocityVector(velocity)

	system.moving = true
}

func (system *playerSystem) MoveRight(_ donburi.World, _ types.Nil) {
	velocity := system.body.Velocity()
	velocity.X += acceleration
	velocity.X = min(velocity.X, maxSpeed)
	system.body.SetVelocityVector(velocity)

	system.moving = true
}

func (system *playerSystem) Shoot(w donburi.World, _ types.Nil) {
	if system.fireTick%uint64(system.settingsProvider.Settings().Internals.TPS/8) == 0 {
		player, _ := actors.PlayerTag.First(w)
		body := physics.Body.Get(player)
		playerWorld := body.Position()

		weapons.BulletFactory.Spawn(system.injector, weapons.NewBulletParams{
			From: playerWorld,
			To:   cp.Vector(system.input.Cursor().WorldPosition),
		})
	}
	system.fireTick++

	system.shooting = true
}
