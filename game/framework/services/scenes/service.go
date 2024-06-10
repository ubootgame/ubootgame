package scenes

import (
	"errors"
	"github.com/ubootgame/ubootgame/framework"
)

type SceneMap map[framework.SceneID]func() framework.Scene

type Service struct {
	scenes map[framework.SceneID]func() framework.Scene
}

func NewService(scenes SceneMap) *Service {
	return &Service{scenes: scenes}
}

func (service *Service) Get(id framework.SceneID) (framework.Scene, error) {
	if scene, ok := service.scenes[id]; ok {
		return scene(), nil
	} else {
		return nil, errors.New("scene not found")
	}
}
