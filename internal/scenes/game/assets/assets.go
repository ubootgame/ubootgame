package assets

import (
	_ "embed"
	"github.com/quasilyte/ebitengine-resource"
	"github.com/ubootgame/ubootgame/pkg/services/resources"
)

const (
	Battleship resources.ImageID = iota
	Submarine
)

const (
	Water resources.TilesheetID = iota
)

const (
	AnimatedWater resources.AsepriteID = iota
)

var Assets = &resources.Library{
	Images: map[resource.ImageID]resource.ImageInfo{
		Battleship: {Path: "military-boats-collection/ship1.png"},
		Submarine:  {Path: "military-boats-collection/submarine1.png"},
	},
	Tilesheets: map[resources.TilesheetID]resources.TilesheetInfo{
		Water: {Path: "water/fishSpritesheet.xml"},
	},
	Aseprites: map[resources.AsepriteID]resources.AsepriteInfo{
		AnimatedWater: {Path: "water/water.json"},
	},
}
