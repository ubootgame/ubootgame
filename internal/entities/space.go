package entities

import (
	"github.com/solarlune/resolv"
	"github.com/ubootgame/ubootgame/internal/components"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var Space = utility.NewArchetype(
	components.Space,
)

func CreateSpace(ecs *ecs.ECS) *donburi.Entry {
	space := Space.Spawn(ecs)

	cfg := config.C
	spaceData := resolv.NewSpace(cfg.Width, cfg.Height, 16, 16)
	components.Space.Set(space, spaceData)

	return space
}
