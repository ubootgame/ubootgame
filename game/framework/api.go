package framework

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ebitengine-resource"
	"github.com/ubootgame/ubootgame/framework/services/resources"
	"github.com/ubootgame/ubootgame/framework/services/settings"
)

type ResourceService interface {
	RegisterResources(library *resources.Library) error
	NextImageID() resources.ImageID
	NextAudioID() resources.AudioID
	RegisterImage(id resources.ImageID, info resources.ImageInfo)
	RegisterAudio(id resources.AudioID, info resources.AudioInfo)
	LoadImage(id resources.ImageID) resource.Image
	LoadAudio(id resources.AudioID) resource.Audio
	LoadTile(id resources.TilesheetID, name string) resource.Image
	LoadAseprite(id resources.AsepriteID) resources.Aseprite
	Preload()
}

type SettingsService[S any] interface {
	Settings() *settings.Settings[S]
	UpdateDebugFontFace(scalingFactor float64)
}

type SceneService interface {
	Get(id SceneID) (Scene, error)
}

type DisplayService interface {
	VirtualResolution() (float64, float64)
	UpdateVirtualResolution(width, height int, scaleFactor float64) (float64, float64)
}

type Scene interface {
	Load() error
	Update() error
	Draw(screen *ebiten.Image)
}

type SceneID string
