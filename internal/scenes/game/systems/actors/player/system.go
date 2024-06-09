package player

import (
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/actors"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/weapons"
	"github.com/ubootgame/ubootgame/internal/scenes/game/tags"
	ecsFramework "github.com/ubootgame/ubootgame/pkg/ecs"
	"github.com/ubootgame/ubootgame/pkg/input"
	"github.com/ubootgame/ubootgame/pkg/settings"
	"github.com/ubootgame/ubootgame/pkg/world"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
	"go/types"
	"gonum.org/v1/gonum/spatial/r2"
)

const (
	acceleration = 10.0
	friction     = 50.0
	maxSpeed     = 250.0
)

type System struct {
	ecsFramework.System

	ecs      *ecs.ECS
	settings *settings.Settings[internal.Settings]
	cursor   *input.Cursor

	transform *transform.TransformData
	velocity  *r2.Vec

	moving, shooting bool
	fireTick         uint64
}

func NewPlayerSystem(ecs *ecs.ECS, settings *settings.Settings[internal.Settings], cursor *input.Cursor) *System {
	system := &System{ecs: ecs, settings: settings, cursor: cursor}
	system.Injector = ecsFramework.NewInjector([]ecsFramework.Injection{
		ecsFramework.WithTag(actors.PlayerTag, []ecsFramework.Injection{
			ecsFramework.Component(&system.velocity, geometry.Velocity),
			ecsFramework.Component(&system.transform, transform.Transform),
		}),
	})

	MoveLeftEvent.Subscribe(ecs.World, system.MoveLeft)
	MoveRightEvent.Subscribe(ecs.World, system.MoveRight)
	ShootEvent.Subscribe(ecs.World, system.Shoot)

	return system
}

func (system *System) Update(e *ecs.ECS) {
	system.System.Update(e)

	if system.moving {
		system.moving = false
	} else {
		system.velocity.X *= 1 - friction/world.WorldSizeBase
	}

	if system.shooting {
		system.shooting = false
	} else {
		system.fireTick = 0
	}
}

func (system *System) MoveLeft(_ donburi.World, _ types.Nil) {
	if system.velocity.X > 0 {
		system.velocity.X *= 1 - friction/world.WorldSizeBase
	}
	system.velocity.X -= acceleration
	system.velocity.X = max(system.velocity.X, -maxSpeed)

	system.moving = true
}

func (system *System) MoveRight(_ donburi.World, _ types.Nil) {
	if system.velocity.X < 0 {
		system.velocity.X *= 1 - friction/world.WorldSizeBase
	}
	system.velocity.X += acceleration
	system.velocity.X = min(system.velocity.X, maxSpeed)

	system.moving = true
}

func (system *System) Shoot(w donburi.World, _ types.Nil) {
	if system.fireTick%uint64(system.settings.TargetTPS/8) == 0 {
		player, _ := actors.PlayerTag.First(w)
		playerWorld := transform.WorldPosition(player)
		bullet := weapons.CreateBullet(system.ecs, r2.Vec(playerWorld), system.cursor.WorldPosition)

		projectilesGroup, _ := tags.ProjectilesTag.First(w)
		transform.AppendChild(projectilesGroup, bullet, true)
	}
	system.fireTick++

	system.shooting = true
}
