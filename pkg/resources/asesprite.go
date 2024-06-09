package resources

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/goaseprite"
	"github.com/ubootgame/ubootgame/assets"
	"path"
)

type AsepriteID int

type AsepriteInfo struct {
	Path string
}

type AsepriteEntry struct {
	ImageID ImageID
	Player  *goaseprite.Player
}

type Aseprite struct {
	ID     AsepriteID
	Image  *ebiten.Image
	Player *goaseprite.Player
}

func LoadAseprite(info AsepriteInfo, registry *Registry) (AsepriteEntry, error) {
	json, _ := assets.FS.ReadFile(info.Path)
	file := goaseprite.Read(json)
	player := file.CreatePlayer()

	dir := path.Dir(info.Path)
	imagePath := path.Join(dir, file.ImagePath)
	imageInfo := ImageInfo{Path: imagePath}

	imageID := registry.NextImageID()

	registry.RegisterImage(imageID, imageInfo)

	aseprite := AsepriteEntry{
		ImageID: imageID,
		Player:  player,
	}

	return aseprite, nil
}
