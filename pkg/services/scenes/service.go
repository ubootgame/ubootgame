package scenes

import (
	"errors"
	"github.com/ubootgame/ubootgame/pkg"
)

type SceneMap map[pkg.SceneID]func() pkg.Scene

type Service struct {
	scenes map[pkg.SceneID]func() pkg.Scene
}

func NewService(scenes SceneMap) *Service {
	return &Service{scenes: scenes}
}

func (service *Service) Get(id pkg.SceneID) (pkg.Scene, error) {
	if scene, ok := service.scenes[id]; ok {
		return scene(), nil
	} else {
		return nil, errors.New("scene not found")
	}
}
