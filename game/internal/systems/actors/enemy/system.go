package enemy

import (
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/framework"
	ecsFramework "github.com/ubootgame/ubootgame/framework/game/ecs"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/components/graphics"
	"github.com/ubootgame/ubootgame/internal/components/physics"
	"github.com/ubootgame/ubootgame/internal/entities/actors"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type System struct {
	settings framework.SettingsService[internal.Settings]

	query *donburi.Query
}

func NewSystem() *System {
	system := &System{query: donburi.NewQuery(filter.Contains(actors.EnemyTag))}

	return system
}

func (system *System) Layers() []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer]{}
}

func (system *System) Update(e *ecs.ECS) {
	system.query.Each(e.World, func(entry *donburi.Entry) {
		body := physics.Body.Get(entry)
		sprite := graphics.Sprite.Get(entry)

		if body.Velocity().X < 0 {
			sprite.FlipY = true
		} else if body.Velocity().X > 0 {
			sprite.FlipY = false
		}
	})
}
