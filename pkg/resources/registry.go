package resources

import (
	"bufio"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/quasilyte/ebitengine-resource"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/assets"
	"image"
	"io"
)

type ImageID = resource.ImageID
type AudioID = resource.AudioID

type ImageInfo = resource.ImageInfo
type AudioInfo = resource.AudioInfo

type Library struct {
	Images     map[ImageID]ImageInfo
	Audio      map[AudioID]AudioInfo
	Tilesheets map[TilesheetID]TilesheetInfo
	Aseprites  map[AsepriteID]AsepriteInfo
}

type Registry struct {
	l              *resource.Loader
	highestImageID ImageID
	highestAudioID AudioID
	images         []ImageID
	audio          []AudioID
	tilesheets     map[TilesheetID]TilesheetEntry
	aseprites      map[AsepriteID]AsepriteEntry
}

func NewRegistry(audioContext *audio.Context) *Registry {
	l := resource.NewLoader(audioContext)

	registry := Registry{l: l,
		tilesheets: map[TilesheetID]TilesheetEntry{},
		aseprites:  map[AsepriteID]AsepriteEntry{},
	}

	l.OpenAssetFunc = registry.openAsset

	return &registry
}

func (registry *Registry) openAsset(path string) io.ReadCloser {
	f := lo.Must(assets.FS.Open(path))
	reader := bufio.NewReader(f)

	return io.NopCloser(reader)
}

func (registry *Registry) RegisterResources(library *Library) error {
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
	if library.Tilesheets != nil {
		for key, info := range library.Tilesheets {
			tileSheet, err := LoadTileSheet(info, registry)
			if err != nil {
				return err
			}
			registry.tilesheets[key] = tileSheet
		}
	}
	if library.Aseprites != nil {
		for key, info := range library.Aseprites {
			aseprite, err := LoadAseprite(info, registry)
			if err != nil {
				return err
			}
			registry.aseprites[key] = aseprite
		}
	}
	return nil
}

func (registry *Registry) NextImageID() ImageID {
	return ImageID(int(registry.highestImageID) + 1)
}

func (registry *Registry) NextAudioID() AudioID {
	return AudioID(int(registry.highestAudioID) + 1)
}

func (registry *Registry) RegisterImage(id ImageID, info ImageInfo) {
	if int(id) > int(registry.highestImageID) {
		registry.highestImageID = id
	}
	registry.l.ImageRegistry.Set(id, info)
}

func (registry *Registry) RegisterAudio(id AudioID, info AudioInfo) {
	if int(id) > int(registry.highestAudioID) {
		registry.highestAudioID = id
	}
	registry.l.AudioRegistry.Set(id, info)
}

func (registry *Registry) LoadImage(id ImageID) resource.Image {
	return registry.l.LoadImage(id)
}

func (registry *Registry) LoadAudio(id AudioID) resource.Audio {
	return registry.l.LoadAudio(id)
}

func (registry *Registry) LoadTile(id TilesheetID, name string) resource.Image {
	tileSheet := registry.tilesheets[id]
	tile := tileSheet.Tiles[name]
	img := registry.l.LoadImage(tileSheet.ImageID)
	tileImage := img.Data.SubImage(image.Rect(tile.X, tile.Y, tile.X+tile.Width, tile.Y+tile.Height)).(*ebiten.Image)

	return resource.Image{
		ID:                 tileSheet.ImageID, // FIXME: this is now the ID of the source image
		Data:               tileImage,
		DefaultFrameWidth:  float64(tileImage.Bounds().Size().X),
		DefaultFrameHeight: float64(tileImage.Bounds().Size().Y),
	}
}

func (registry *Registry) LoadAseprite(id AsepriteID) Aseprite {
	asepriteEntry := registry.aseprites[id]
	img := registry.l.LoadImage(asepriteEntry.ImageID)

	return Aseprite{
		ID:     id,
		Image:  img.Data,
		Player: asepriteEntry.Player,
	}
}

func (registry *Registry) Preload() {
	for _, id := range registry.images {
		registry.l.LoadImage(id)
	}
	for _, id := range registry.audio {
		registry.l.LoadAudio(id)
	}
}
