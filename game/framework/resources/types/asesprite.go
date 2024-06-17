package types

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
	ImageID   ImageID
	ImageInfo ImageInfo
	Player    *goaseprite.Player
}

type Aseprite struct {
	ID     AsepriteID
	Image  *ebiten.Image
	Player *goaseprite.Player
}

func LoadAseprite(info AsepriteInfo) (AsepriteEntry, error) {
	json, _ := assets.FS.ReadFile(info.Path)
	file := goaseprite.Read(json)
	player := file.CreatePlayer()

	dir := path.Dir(info.Path)
	imagePath := path.Join(dir, file.ImagePath)
	imageInfo := ImageInfo{Path: imagePath}

	aseprite := AsepriteEntry{
		ImageInfo: imageInfo,
		Player:    player,
	}

	return aseprite, nil
}
