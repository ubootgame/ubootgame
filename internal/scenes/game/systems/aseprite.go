package systems

import (
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type asepriteSystem struct{}

var Aseprite = &asepriteSystem{}

func (system *asepriteSystem) Update(e *ecs.ECS) {
	components.Aseprite.Each(e.World, func(entry *donburi.Entry) {
		aseprite := components.Aseprite.Get(entry)
		aseprite.Aseprite.Player.Update(1.0 / float32(config.C.TargetTPS) * aseprite.Speed)
	})
}
