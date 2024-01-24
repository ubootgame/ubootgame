package systems

import (
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateAseprites(e *ecs.ECS) {
	components.Aseprite.Each(e.World, func(entry *donburi.Entry) {
		asepriteData := components.Aseprite.Get(entry)
		asepriteData.Aseprite.Player.Update(1.0 / float32(config.C.TargetTPS) * asepriteData.Speed)
	})
}
