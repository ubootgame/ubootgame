package visuals

import (
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/visuals"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/systems"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type AsepriteSystem struct {
	systems.BaseSystem
}

func NewAsepriteSystem() *AsepriteSystem {
	return &AsepriteSystem{}
}

func (system *AsepriteSystem) Update(e *ecs.ECS) {
	system.BaseSystem.Update(e)

	visuals.Aseprite.Each(e.World, func(entry *donburi.Entry) {
		aseprite := visuals.Aseprite.Get(entry)
		aseprite.Aseprite.Player.Update(1.0 / float32(config.C.TargetTPS) * aseprite.Speed)
	})
}
