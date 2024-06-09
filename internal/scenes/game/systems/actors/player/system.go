package player

import (
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/framework/coordinate_system"
	ecs2 "github.com/ubootgame/ubootgame/internal/framework/ecs"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/actors"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/weapons"
	"github.com/ubootgame/ubootgame/internal/scenes/game/tags"
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
	ecs2.System

	ecs      *ecs.ECS
	settings *internal.Settings

	cursor    *game_system.CursorData
	transform *transform.TransformData
	velocity  *r2.Vec

	moving, shooting bool
	fireTick         uint64
}

func NewPlayerSystem(ecs *ecs.ECS, settings *internal.Settings) *System {
	system := &System{ecs: ecs, settings: settings}
	system.Injector = ecs2.NewInjector([]ecs2.Injection{
		ecs2.Once([]ecs2.Injection{
			ecs2.Component(&system.cursor, game_system.Cursor),
		}),
		ecs2.WithTag(actors.PlayerTag, []ecs2.Injection{
			ecs2.Component(&system.velocity, geometry.Velocity),
			ecs2.Component(&system.transform, transform.Transform),
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
		system.velocity.X *= 1 - friction/coordinate_system.WorldSizeBase
	}

	if system.shooting {
		system.shooting = false
	} else {
		system.fireTick = 0
	}
}

func (system *System) MoveLeft(_ donburi.World, _ types.Nil) {
	if system.velocity.X > 0 {
		system.velocity.X *= 1 - friction/coordinate_system.WorldSizeBase
	}
	system.velocity.X -= acceleration
	system.velocity.X = max(system.velocity.X, -maxSpeed)

	system.moving = true
}

func (system *System) MoveRight(_ donburi.World, _ types.Nil) {
	if system.velocity.X < 0 {
		system.velocity.X *= 1 - friction/coordinate_system.WorldSizeBase
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
