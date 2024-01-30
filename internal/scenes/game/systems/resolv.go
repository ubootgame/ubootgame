package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	dresolv "github.com/ubootgame/ubootgame/internal/resolv"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type resolvSystem struct {
	query *donburi.Query
}

var Resolv = &resolvSystem{
	query: donburi.NewQuery(filter.Contains(components.Position, components.Object)),
}

func (system *resolvSystem) Update(e *ecs.ECS) {
	system.query.Each(e.World, func(entry *donburi.Entry) {
		position := components.Position.Get(entry)
		object := dresolv.GetObject(entry)

		object.X = position.Center.X
		object.Y = position.Center.Y

		object.Update()
	})
}

func (system *resolvSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	//donburi.NewQuery(filter.Contains(components.Position, components.Object)).Each(e.World, func(debugEntry *donburi.Entry) {
	//	positionData := components.Position.Get(debugEntry)
	//	object := dresolv.GetObject(debugEntry)
	//
	//	object.X = positionData.Center.X
	//	object.Y = positionData.Center.Y
	//
	//	object.Update()
	//})
}
