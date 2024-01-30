package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities"
	"github.com/yohamta/donburi/ecs"
)

func UpdatePlayer(ecs *ecs.ECS) {
	entry, _ := entities.PlayerTag.First(ecs.World)

	acceleration := 0.0001
	friction := 0.05
	maxSpeed := 0.005

	velocityData := components.Velocity.Get(entry)

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if velocityData.X > 0 {
			velocityData.X *= 1 - friction
		}
		velocityData.X -= acceleration
		velocityData.X = max(velocityData.X, -maxSpeed)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if velocityData.X < 0 {
			velocityData.X *= 1 - friction
		}
		velocityData.X += acceleration
		velocityData.X = min(velocityData.X, maxSpeed)
	} else {
		velocityData.X *= 1 - friction
	}
}
