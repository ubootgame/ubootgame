package types

import (
	"encoding/xml"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/assets"
	"path"
)

type TilesheetID int

type TilesheetInfo struct {
	Path string
}

type TilesheetEntry struct {
	ImageID   ImageID
	ImageInfo ImageInfo
	Tiles     map[string]Tile
}

type Tile struct {
	X, Y, Width, Height int
}

func LoadTileSheet(info TilesheetInfo) (TilesheetEntry, error) {
	xmlFile := lo.Must(assets.FS.ReadFile(info.Path))

	var textureAtlas TextureAtlas
	err := xml.Unmarshal(xmlFile, &textureAtlas)
	if err != nil {
		return TilesheetEntry{}, err
	}

	dir := path.Dir(info.Path)
	imagePath := path.Join(dir, textureAtlas.ImagePath)
	imageInfo := ImageInfo{Path: imagePath}

	tileSheet := TilesheetEntry{
		ImageInfo: imageInfo,
		Tiles: lo.Associate(textureAtlas.SubTextures, func(item SubTexture) (string, Tile) {
			return item.Name, Tile{
				X:      item.X,
				Y:      item.Y,
				Width:  item.Width,
				Height: item.Height,
			}
		}),
	}

	return tileSheet, nil
}

type TextureAtlas struct {
	XMLName     xml.Name     `xml:"TextureAtlas"`
	ImagePath   string       `xml:"imagePath,attr"`
	SubTextures []SubTexture `xml:"SubTexture"`
}

type SubTexture struct {
	Name   string `xml:"name,attr"`
	X      int    `xml:"x,attr"`
	Y      int    `xml:"y,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}
