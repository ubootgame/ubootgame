package entities

import (
	"github.com/solarlune/resolv"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var Space = utility.NewArchetype(
	components.Space,
)

func CreateSpace(ecs *ecs.ECS) *donburi.Entry {
	entry := Space.Spawn(ecs, layers.Foreground)

	cfg := config.C
	space := resolv.NewSpace(int(cfg.VirtualResolution.X), int(cfg.VirtualResolution.Y), 16, 16)
	components.Space.Set(entry, space)

	return entry
}
