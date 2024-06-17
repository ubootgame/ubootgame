package resources

import "github.com/ubootgame/ubootgame/framework/resources/types"

type Library struct {
	Images     map[types.ImageID]types.ImageInfo
	Audio      map[types.AudioID]types.AudioInfo
	Tilesheets map[types.TilesheetID]types.TilesheetInfo
	Aseprites  map[types.AsepriteID]types.AsepriteInfo
}
