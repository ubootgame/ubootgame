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

type Service struct {
	l              *resource.Loader
	highestImageID ImageID
	highestAudioID AudioID
	images         []ImageID
	audio          []AudioID
	tilesheets     map[TilesheetID]TilesheetEntry
	aseprites      map[AsepriteID]AsepriteEntry
}

func NewService(audioContext *audio.Context) *Service {
	l := resource.NewLoader(audioContext)

	registry := Service{l: l,
		tilesheets: map[TilesheetID]TilesheetEntry{},
		aseprites:  map[AsepriteID]AsepriteEntry{},
	}

	l.OpenAssetFunc = registry.openAsset

	return &registry
}

func (service *Service) openAsset(path string) io.ReadCloser {
	f := lo.Must(assets.FS.Open(path))
	reader := bufio.NewReader(f)

	return io.NopCloser(reader)
}

func (service *Service) RegisterResources(library *Library) error {
	if library.Images != nil {
		for key, info := range library.Images {
			service.RegisterImage(key, info)
		}
	}
	if library.Audio != nil {
		for key, info := range library.Audio {
			service.RegisterAudio(key, info)
		}
	}
	if library.Tilesheets != nil {
		for key, info := range library.Tilesheets {
			tileSheet, err := LoadTileSheet(info, service)
			if err != nil {
				return err
			}
			service.tilesheets[key] = tileSheet
		}
	}
	if library.Aseprites != nil {
		for key, info := range library.Aseprites {
			aseprite, err := LoadAseprite(info, service)
			if err != nil {
				return err
			}
			service.aseprites[key] = aseprite
		}
	}
	return nil
}

func (service *Service) NextImageID() ImageID {
	return ImageID(int(service.highestImageID) + 1)
}

func (service *Service) NextAudioID() AudioID {
	return AudioID(int(service.highestAudioID) + 1)
}

func (service *Service) RegisterImage(id ImageID, info ImageInfo) {
	if int(id) > int(service.highestImageID) {
		service.highestImageID = id
	}
	service.l.ImageRegistry.Set(id, info)
}

func (service *Service) RegisterAudio(id AudioID, info AudioInfo) {
	if int(id) > int(service.highestAudioID) {
		service.highestAudioID = id
	}
	service.l.AudioRegistry.Set(id, info)
}

func (service *Service) LoadImage(id ImageID) resource.Image {
	return service.l.LoadImage(id)
}

func (service *Service) LoadAudio(id AudioID) resource.Audio {
	return service.l.LoadAudio(id)
}

func (service *Service) LoadTile(id TilesheetID, name string) resource.Image {
	tileSheet := service.tilesheets[id]
	tile := tileSheet.Tiles[name]
	img := service.l.LoadImage(tileSheet.ImageID)
	tileImage := img.Data.SubImage(image.Rect(tile.X, tile.Y, tile.X+tile.Width, tile.Y+tile.Height)).(*ebiten.Image)

	return resource.Image{
		ID:                 tileSheet.ImageID, // FIXME: this is now the ID of the source image
		Data:               tileImage,
		DefaultFrameWidth:  float64(tileImage.Bounds().Size().X),
		DefaultFrameHeight: float64(tileImage.Bounds().Size().Y),
	}
}

func (service *Service) LoadAseprite(id AsepriteID) Aseprite {
	asepriteEntry := service.aseprites[id]
	img := service.l.LoadImage(asepriteEntry.ImageID)

	return Aseprite{
		ID:     id,
		Image:  img.Data,
		Player: asepriteEntry.Player,
	}
}

func (service *Service) Preload() {
	for _, id := range service.images {
		service.l.LoadImage(id)
	}
	for _, id := range service.audio {
		service.l.LoadAudio(id)
	}
}
