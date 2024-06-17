package resources

import (
	"bufio"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/quasilyte/ebitengine-resource"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/assets"
	"github.com/ubootgame/ubootgame/framework/resources/types"
	"image"
	"io"
)

type Registry interface {
	RegisterResources(library *Library) error
	NextImageID() types.ImageID
	NextAudioID() types.AudioID
	RegisterImage(id types.ImageID, info types.ImageInfo)
	RegisterAudio(id types.AudioID, info types.AudioInfo)
	LoadImage(id types.ImageID) resource.Image
	LoadAudio(id types.AudioID) resource.Audio
	LoadTile(id types.TilesheetID, name string) resource.Image
	LoadAseprite(id types.AsepriteID) types.Aseprite
	Preload()
}

type registry struct {
	l              *resource.Loader
	highestImageID types.ImageID
	highestAudioID types.AudioID
	images         []types.ImageID
	audio          []types.AudioID
	tilesheets     map[types.TilesheetID]types.TilesheetEntry
	aseprites      map[types.AsepriteID]types.AsepriteEntry
}

func NewRegistry(_ *do.Injector, audioContext *audio.Context) (Registry, error) {
	l := resource.NewLoader(audioContext)

	registry := registry{l: l,
		tilesheets: map[types.TilesheetID]types.TilesheetEntry{},
		aseprites:  map[types.AsepriteID]types.AsepriteEntry{},
	}

	l.OpenAssetFunc = registry.openAsset

	return &registry, nil
}

func (registry *registry) openAsset(path string) io.ReadCloser {
	f := lo.Must(assets.FS.Open(path))
	reader := bufio.NewReader(f)

	return io.NopCloser(reader)
}

func (registry *registry) RegisterResources(library *Library) error {
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
			tileSheet, err := types.LoadTileSheet(info)
			if err != nil {
				return err
			}

			tileSheet.ImageID = registry.NextImageID()
			registry.tilesheets[key] = tileSheet

			registry.RegisterImage(tileSheet.ImageID, tileSheet.ImageInfo)
		}
	}
	if library.Aseprites != nil {
		for key, info := range library.Aseprites {
			aseprite, err := types.LoadAseprite(info)
			if err != nil {
				return err
			}

			aseprite.ImageID = registry.NextImageID()
			registry.aseprites[key] = aseprite

			registry.RegisterImage(aseprite.ImageID, aseprite.ImageInfo)
		}
	}
	return nil
}

func (registry *registry) NextImageID() types.ImageID {
	return types.ImageID(int(registry.highestImageID) + 1)
}

func (registry *registry) NextAudioID() types.AudioID {
	return types.AudioID(int(registry.highestAudioID) + 1)
}

func (registry *registry) RegisterImage(id types.ImageID, info types.ImageInfo) {
	if int(id) > int(registry.highestImageID) {
		registry.highestImageID = id
	}
	registry.l.ImageRegistry.Set(id, info)
}

func (registry *registry) RegisterAudio(id types.AudioID, info types.AudioInfo) {
	if int(id) > int(registry.highestAudioID) {
		registry.highestAudioID = id
	}
	registry.l.AudioRegistry.Set(id, info)
}

func (registry *registry) LoadImage(id types.ImageID) resource.Image {
	return registry.l.LoadImage(id)
}

func (registry *registry) LoadAudio(id types.AudioID) resource.Audio {
	return registry.l.LoadAudio(id)
}

func (registry *registry) LoadTile(id types.TilesheetID, name string) resource.Image {
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

func (registry *registry) LoadAseprite(id types.AsepriteID) types.Aseprite {
	asepriteEntry := registry.aseprites[id]
	img := registry.l.LoadImage(asepriteEntry.ImageID)

	return types.Aseprite{
		ID:     id,
		Image:  img.Data,
		Player: asepriteEntry.Player,
	}
}

func (registry *registry) Preload() {
	for _, id := range registry.images {
		registry.l.LoadImage(id)
	}
	for _, id := range registry.audio {
		registry.l.LoadAudio(id)
	}
}
