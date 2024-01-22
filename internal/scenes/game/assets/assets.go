package assets

import (
	"github.com/quasilyte/ebitengine-resource"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
)

const ImageBattleship = "battleship"
const TileSheetWater = "water"

var Assets = resources.Library{
	Images: map[string]resource.ImageInfo{
		ImageBattleship: {Path: "assets/battleship.png"},
	},
	TileSheets: map[string]resources.TileSheetInfo{
		TileSheetWater: {Path: "assets/water/fishSpritesheet.xml"},
	},
}
