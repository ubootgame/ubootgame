package assets

import (
	_ "embed"
	"github.com/quasilyte/ebitengine-resource"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
)

const (
	Battleship resources.ImageID = iota
)

const (
	Water resources.TilesheetID = iota
)

const (
	AnimatedWater resources.AsepriteID = iota
)

var Assets = &resources.Library{
	Images: map[resources.ImageID]resource.ImageInfo{
		Battleship: {Path: "military-boats-collection/ship1.png"},
	},
	Tilesheets: map[resources.TilesheetID]resources.TilesheetInfo{
		Water: {Path: "water/fishSpritesheet.xml"},
	},
	Aseprites: map[resources.AsepriteID]resources.AsepriteInfo{
		AnimatedWater: {Path: "water/water.json"},
	},
}
