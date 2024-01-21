package assets

import "github.com/quasilyte/ebitengine-resource"

const (
	ImageNone resource.ImageID = iota
	ImageBattleship
)

var ImageResources = map[resource.ImageID]resource.ImageInfo{
	ImageBattleship: {Path: "assets/battleship.png"},
}
