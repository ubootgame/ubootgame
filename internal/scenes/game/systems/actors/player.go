package actors

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/actors"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/weapons"
	"github.com/ubootgame/ubootgame/internal/scenes/game/tags"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/injector"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/systems"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
	"gonum.org/v1/gonum/spatial/r2"
)

type PlayerSystem struct {
	systems.BaseSystem

	cursor    *game_system.CursorData
	transform *transform.TransformData
	velocity  *r2.Vec

	fireTick uint64
}

func NewPlayerSystem() *PlayerSystem {
	system := &PlayerSystem{}
	system.Injector = injector.NewInjector([]injector.Injection{
		injector.Once([]injector.Injection{
			injector.Component(&system.cursor, game_system.Cursor),
		}),
		injector.WithTag(actors.PlayerTag, []injector.Injection{
			injector.Component(&system.velocity, geometry.Velocity),
			injector.Component(&system.transform, transform.Transform),
		}),
	})
	return system
}

func (system *PlayerSystem) Update(e *ecs.ECS) {
	system.BaseSystem.Update(e)

	acceleration := 0.01
	friction := 0.05
	maxSpeed := 0.25

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if system.velocity.X > 0 {
			system.velocity.X *= 1 - friction
		}
		system.velocity.X -= acceleration
		system.velocity.X = max(system.velocity.X, -maxSpeed)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if system.velocity.X < 0 {
			system.velocity.X *= 1 - friction
		}
		system.velocity.X += acceleration
		system.velocity.X = min(system.velocity.X, maxSpeed)
	} else {
		system.velocity.X *= 1 - friction
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if system.fireTick%uint64(config.C.TargetTPS/8) == 0 {
			player, _ := actors.PlayerTag.First(e.World)
			playerWorld := transform.WorldPosition(player)
			bullet := weapons.CreateBullet(e, r2.Vec{X: playerWorld.X, Y: playerWorld.Y}, system.cursor.WorldPosition)

			projectilesGroup, _ := tags.ProjectilesTag.First(e.World)
			transform.AppendChild(projectilesGroup, bullet, true)
		}
		system.fireTick++
	} else {
		system.fireTick = 0
	}
}
