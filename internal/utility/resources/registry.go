package resources

import (
	"bufio"
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/quasilyte/ebitengine-resource"
	"github.com/samber/lo"
	"image"
	"io"
	"os"
	"path"
)

type ImageInfo = resource.ImageInfo
type AudioInfo = resource.AudioInfo
type TileSheetInfo struct {
	Path string
}

type Library struct {
	Data       map[string][]byte
	Images     map[string]ImageInfo
	Audio      map[string]AudioInfo
	TileSheets map[string]TileSheetInfo
}

type Registry struct {
	l              *resource.Loader
	highestImageID int
	highestAudioID int
	imageIDs       map[string]resource.ImageID
	audioIDs       map[string]resource.AudioID
	tileSheets     map[string]TileSheet
	data           map[string][]byte
}

func NewRegistry(audioContext *audio.Context) *Registry {
	l := resource.NewLoader(audioContext)

	registry := Registry{l: l,
		imageIDs:   map[string]resource.ImageID{},
		audioIDs:   map[string]resource.AudioID{},
		tileSheets: map[string]TileSheet{}}

	l.OpenAssetFunc = registry.openAsset

	return &registry
}

func (registry *Registry) openAsset(path string) io.ReadCloser {
	if data, ok := registry.data[path]; ok {
		return io.NopCloser(bytes.NewReader(data))
	}

	f := lo.Must(os.Open(path))
	reader := bufio.NewReader(f)

	return io.NopCloser(reader)
}

func (registry *Registry) RegisterResources(library *Library) error {
	registry.data = library.Data

	if library.Images != nil {
		for key, info := range library.Images {
			registry.RegisterImage(key, info)
		}
	}
	if library.Audio != nil {
		for key, info := range library.Audio {
			registry.RegisterAudio(key, info)
		}
	}
	if library.TileSheets != nil {
		for key, info := range library.TileSheets {
			tileSheet, err := LoadTileSheet(info.Path, library)
			if err != nil {
				return err
			}
			registry.tileSheets[key] = tileSheet

			dir := path.Dir(info.Path)
			imageInfo := ImageInfo{
				Path: path.Join(dir, tileSheet.Path),
			}

			registry.RegisterImage(key, imageInfo)
		}
	}
	return nil
}

func (registry *Registry) RegisterImage(key string, info ImageInfo) {
	imageID := resource.ImageID(registry.highestImageID)
	registry.imageIDs[key] = imageID
	registry.l.ImageRegistry.Set(imageID, info)
	registry.highestImageID++
}

func (registry *Registry) RegisterAudio(key string, info AudioInfo) {
	audioID := resource.AudioID(registry.highestAudioID)
	registry.audioIDs[key] = audioID
	registry.l.AudioRegistry.Set(audioID, info)
	registry.highestAudioID++
}

func (registry *Registry) LoadImage(id string) resource.Image {
	return registry.l.LoadImage(registry.imageIDs[id])
}

func (registry *Registry) LoadAudio(id string) resource.Audio {
	return registry.l.LoadAudio(registry.audioIDs[id])
}

func (registry *Registry) LoadTile(id string, name string) resource.Image {
	tileSheet := registry.tileSheets[id]
	tile := tileSheet.Tiles[name]
	img := registry.l.LoadImage(registry.imageIDs[id])
	tileImage := img.Data.SubImage(image.Rect(tile.X, tile.Y, tile.X+tile.Width, tile.Y+tile.Height)).(*ebiten.Image)

	return resource.Image{
		ID:                 registry.imageIDs[id],
		Data:               tileImage,
		DefaultFrameWidth:  float64(tileImage.Bounds().Size().X),
		DefaultFrameHeight: float64(tileImage.Bounds().Size().Y),
	}
}

func (registry *Registry) Preload() {
	for _, id := range registry.imageIDs {
		registry.l.LoadImage(id)
	}
	for _, id := range registry.audioIDs {
		registry.l.LoadAudio(id)
	}
}
