package actors

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/framework/coordinate_system"
	ecs2 "github.com/ubootgame/ubootgame/internal/framework/ecs"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/actors"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/weapons"
	"github.com/ubootgame/ubootgame/internal/scenes/game/tags"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
	"gonum.org/v1/gonum/spatial/r2"
)

type PlayerSystem struct {
	ecs2.System

	settings *internal.Settings

	cursor    *game_system.CursorData
	transform *transform.TransformData
	velocity  *r2.Vec

	fireTick uint64
}

func NewPlayerSystem(settings *internal.Settings) *PlayerSystem {
	system := &PlayerSystem{settings: settings}
	system.Injector = ecs2.NewInjector([]ecs2.Injection{
		ecs2.Once([]ecs2.Injection{
			ecs2.Component(&system.cursor, game_system.Cursor),
		}),
		ecs2.WithTag(actors.PlayerTag, []ecs2.Injection{
			ecs2.Component(&system.velocity, geometry.Velocity),
			ecs2.Component(&system.transform, transform.Transform),
		}),
	})
	return system
}

func (system *PlayerSystem) Update(e *ecs.ECS) {
	system.System.Update(e)

	acceleration := 10.0
	friction := 50.0
	maxSpeed := 250.0

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if system.velocity.X > 0 {
			system.velocity.X *= 1 - friction/coordinate_system.WorldSizeBase
		}
		system.velocity.X -= acceleration
		system.velocity.X = max(system.velocity.X, -maxSpeed)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if system.velocity.X < 0 {
			system.velocity.X *= 1 - friction/coordinate_system.WorldSizeBase
		}
		system.velocity.X += acceleration
		system.velocity.X = min(system.velocity.X, maxSpeed)
	} else {
		system.velocity.X *= 1 - friction/coordinate_system.WorldSizeBase
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if system.fireTick%uint64(system.settings.TargetTPS/8) == 0 {
			player, _ := actors.PlayerTag.First(e.World)
			playerWorld := transform.WorldPosition(player)
			bullet := weapons.CreateBullet(e, r2.Vec(playerWorld), system.cursor.WorldPosition)

			projectilesGroup, _ := tags.ProjectilesTag.First(e.World)
			transform.AppendChild(projectilesGroup, bullet, true)
		}
		system.fireTick++
	} else {
		system.fireTick = 0
	}
}
