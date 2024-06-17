package enemy

import (
	"github.com/samber/do"
	"github.com/samber/lo"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/internal/components/graphics"
	"github.com/ubootgame/ubootgame/internal/components/physics"
	"github.com/ubootgame/ubootgame/internal/entities/actors"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type enemySystem struct {
	query *donburi.Query
}

func NewEnemySystem(_ *do.Injector) ecsFramework.System {
	return &enemySystem{query: donburi.NewQuery(filter.Contains(actors.EnemyTag))}
}

func (system *enemySystem) Layers() []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer]{}
}

func (system *enemySystem) Update(e *ecs.ECS) {
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
