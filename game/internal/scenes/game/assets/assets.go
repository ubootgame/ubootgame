package assets

import (
	_ "embed"
	"github.com/quasilyte/ebitengine-resource"
	types2 "github.com/ubootgame/ubootgame/framework/resources"
	"github.com/ubootgame/ubootgame/framework/resources/types"
)

const (
	Battleship types.ImageID = iota
	Submarine
)

const (
	Water types.TilesheetID = iota
)

const (
	AnimatedWater types.AsepriteID = iota
)

var Assets = &types2.Library{
	Images: map[resource.ImageID]resource.ImageInfo{
		Battleship: {Path: "military-boats-collection/ship1.png"},
		Submarine:  {Path: "military-boats-collection/submarine1.png"},
	},
	Tilesheets: map[types.TilesheetID]types.TilesheetInfo{
		Water: {Path: "water/fishSpritesheet.xml"},
	},
	Aseprites: map[types.AsepriteID]types.AsepriteInfo{
		AnimatedWater: {Path: "water/water.json"},
	},
}
