package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/components"
	"github.com/ubootgame/ubootgame/internal/entities"
	dresolv "github.com/ubootgame/ubootgame/internal/resolv"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func DrawShip(e *ecs.ECS, screen *ebiten.Image) {
	shipEntry, _ := ecs.NewQuery(ecs.LayerDefault, filter.Contains(entities.ShipTag)).First(e.World)

	spriteData := components.Sprite.Get(shipEntry)

	object := dresolv.GetObject(shipEntry)

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(object.X, object.Y)
	screen.DrawImage(spriteData.Image, options)
}
