package assets

import (
	_ "embed"
	"github.com/quasilyte/ebitengine-resource"
	"github.com/ubootgame/ubootgame/assets"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
)

const ImageBattleship = "battleship"
const TileSheetWater = "water"

var Assets = &resources.Library{
	Images: map[string]resource.ImageInfo{
		ImageBattleship: {Path: "assets/military-boats-collection/ship1.png"},
	},
	TileSheets: map[string]resources.TileSheetInfo{
		TileSheetWater: {Path: "assets/water/fishSpritesheet.xml"},
	},
	Data: map[string][]byte{
		"assets/military-boats-collection/ship1.png": assets.Ship1,
		"assets/water/fishSpritesheet.png":           assets.FishSpriteSheet,
		"assets/water/fishSpritesheet.xml":           assets.FishSpriteSheetXML,
	},
}
