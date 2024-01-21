package utility

import (
	"bufio"
	"github.com/hajimehoshi/ebiten/v2/audio"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/samber/lo"
	"io"
	"os"
)

type Resources struct {
	Loader *resource.Loader
}

func CreateResourceLoader(audioContext *audio.Context) *resource.Loader {
	l := resource.NewLoader(audioContext)

	l.OpenAssetFunc = func(path string) io.ReadCloser {
		f := lo.Must(os.Open(path))
		reader := bufio.NewReader(f)

		return io.NopCloser(reader)
	}

	return l
}

func (r *Resources) LoadImages(imageResources map[resource.ImageID]resource.ImageInfo) {
	r.Loader.ImageRegistry.Assign(imageResources)

	for id := range imageResources {
		r.Loader.LoadImage(id)
	}
}
