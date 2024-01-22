package resources

import (
	"encoding/xml"
	"github.com/samber/lo"
	"io"
	"os"
)

type TileSheet struct {
	Path  string
	Tiles map[string]Tile
}

type Tile struct {
	X, Y, Width, Height int
}

func LoadTileSheet(path string) (TileSheet, error) {
	xmlFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		return TileSheet{}, err
	}
	defer func(xmlFile *os.File) {
		_ = xmlFile.Close()
	}(xmlFile)

	byteValue, _ := io.ReadAll(xmlFile)

	var textureAtlas TextureAtlas
	err = xml.Unmarshal(byteValue, &textureAtlas)
	if err != nil {
		return TileSheet{}, err
	}

	tileSheet := TileSheet{
		Path: textureAtlas.ImagePath,
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
