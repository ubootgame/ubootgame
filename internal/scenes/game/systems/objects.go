package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	dresolv "github.com/ubootgame/ubootgame/internal/resolv"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type objects struct{}

var Objects = &objects{}

func (system *objects) Update(e *ecs.ECS) {
	donburi.NewQuery(filter.Contains(components.Position, components.Object)).Each(e.World, func(entry *donburi.Entry) {
		positionData := components.Position.Get(entry)
		object := dresolv.GetObject(entry)

		object.X = positionData.Center.X
		object.Y = positionData.Center.Y

		object.Update()
	})
}

func (system *objects) Draw(e *ecs.ECS, screen *ebiten.Image) {
	//donburi.NewQuery(filter.Contains(components.Position, components.Object)).Each(e.World, func(entry *donburi.Entry) {
	//	positionData := components.Position.Get(entry)
	//	object := dresolv.GetObject(entry)
	//
	//	object.X = positionData.Center.X
	//	object.Y = positionData.Center.Y
	//
	//	object.Update()
	//})
}
