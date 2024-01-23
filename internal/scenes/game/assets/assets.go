package assets

import (
	_ "embed"
	"github.com/quasilyte/ebitengine-resource"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
)

const ImageBattleship = "battleship"
const TileSheetWater = "water"

var Assets = &resources.Library{
	Images: map[string]resource.ImageInfo{
		ImageBattleship: {Path: "military-boats-collection/ship1.png"},
	},
	TileSheets: map[string]resources.TileSheetInfo{
		TileSheetWater: {Path: "water/fishSpritesheet.xml"},
	},
}
