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

	sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
	w, h := float64(spriteData.Image.Bounds().Dx()), float64(spriteData.Image.Bounds().Dy())
	sizeScale := 0.1 * (sw / w)
	deviceScale := ebiten.DeviceScaleFactor()

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(-w)/2, float64(-h)/2)
	op.GeoM.Scale(sizeScale*deviceScale, sizeScale*deviceScale)
	op.GeoM.Translate(float64(sw)/2+object.X, float64(sh)/2+object.Y)

	screen.DrawImage(spriteData.Image, op)
}
