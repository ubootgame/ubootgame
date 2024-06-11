package player

import (
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/framework"
	ecsFramework "github.com/ubootgame/ubootgame/framework/game/ecs"
	"github.com/ubootgame/ubootgame/framework/game/input"
	"github.com/ubootgame/ubootgame/framework/game/world"
	"github.com/ubootgame/ubootgame/internal"
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
	settings framework.SettingsService[internal.Settings]

	ecs    *ecs.ECS
	cursor *input.Cursor

	player *donburi.Entry

	moving, shooting bool
	fireTick         uint64
}

func NewSystem(settings framework.SettingsService[internal.Settings], ecs *ecs.ECS, cursor *input.Cursor) *System {
	system := &System{ecs: ecs, settings: settings, cursor: cursor}

	MoveLeftEvent.Subscribe(ecs.World, system.MoveLeft)
	MoveRightEvent.Subscribe(ecs.World, system.MoveRight)
	ShootEvent.Subscribe(ecs.World, system.Shoot)

	return system
}

func (system *System) Layers() []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer]{}
}

func (system *System) Update(e *ecs.ECS) {
	if player, found := actors.PlayerTag.First(e.World); found {
		system.player = player
	} else {
		return
	}

	velocity := geometry.Velocity.Get(system.player)

	if system.moving {
		system.moving = false
	} else {
		velocity.X *= 1 - friction/world.BaseSize
	}

	if system.shooting {
		system.shooting = false
	} else {
		system.fireTick = 0
	}
}

func (system *System) MoveLeft(_ donburi.World, _ types.Nil) {
	velocity := geometry.Velocity.Get(system.player)

	if velocity.X > 0 {
		velocity.X *= 1 - friction/world.BaseSize
	}
	velocity.X -= acceleration
	velocity.X = max(velocity.X, -maxSpeed)

	system.moving = true
}

func (system *System) MoveRight(_ donburi.World, _ types.Nil) {
	velocity := geometry.Velocity.Get(system.player)

	if velocity.X < 0 {
		velocity.X *= 1 - friction/world.BaseSize
	}
	velocity.X += acceleration
	velocity.X = min(velocity.X, maxSpeed)

	system.moving = true
}

func (system *System) Shoot(w donburi.World, _ types.Nil) {
	if system.fireTick%uint64(system.settings.Settings().Internals.TPS/8) == 0 {
		player, _ := actors.PlayerTag.First(w)
		playerWorld := transform.WorldPosition(player)
		bullet := weapons.CreateBullet(system.ecs, r2.Vec(playerWorld), system.cursor.WorldPosition)

		projectilesGroup, _ := tags.ProjectilesTag.First(w)
		transform.AppendChild(projectilesGroup, bullet, true)
	}
	system.fireTick++

	system.shooting = true
}
