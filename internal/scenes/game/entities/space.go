package entities

import (
	"github.com/solarlune/resolv"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"gonum.org/v1/gonum/spatial/r2"
)

var Space = utility.NewArchetype(
	components.Space,
)

func CreateSpace(ecs *ecs.ECS, size r2.Vec) *donburi.Entry {
	entry := Space.Spawn(ecs, layers.Foreground)

	space := resolv.NewSpace(int(size.X), int(size.Y), 16, 16)
	components.Space.Set(entry, space)

	return entry
}
