package systems

import (
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/pkg"
	ecsFramework "github.com/ubootgame/ubootgame/pkg/ecs"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
)

type MovementSystem struct {
	ecsFramework.System

	settings pkg.SettingsService[internal.Settings]

	query *donburi.Query
}

func NewMovementSystem(settings pkg.SettingsService[internal.Settings]) *MovementSystem {
	system := &MovementSystem{
		settings: settings,
		query:    donburi.NewQuery(filter.Contains(transform.Transform, geometry.Velocity)),
	}
	return system
}

func (system *MovementSystem) Update(e *ecs.ECS) {
	system.System.Update(e)

	system.query.Each(e.World, func(entry *donburi.Entry) {
		velocity := geometry.Velocity.Get(entry)
		t := transform.Transform.Get(entry)

		t.LocalPosition.X += velocity.X / float64(system.settings.Settings().Internals.TPS)
		t.LocalPosition.Y += velocity.Y / float64(system.settings.Settings().Internals.TPS)

		if entry.HasComponent(geometry.Direction) {
			direction := geometry.Direction.Get(entry)

			if velocity.X < 0 {
				direction.Horizontal = geometry.Left
			} else if velocity.X > 0 {
				direction.Horizontal = geometry.Right
			}

			if velocity.Y < 0 {
				direction.Vertical = geometry.Up
			} else if velocity.Y > 0 {
				direction.Vertical = geometry.Down
			}
		}
	})
}
