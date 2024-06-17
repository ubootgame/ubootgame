package scenes

import (
	"github.com/ubootgame/ubootgame/framework/game"
	gameScene "github.com/ubootgame/ubootgame/internal/scenes/game"
)

var Scenes = game.SceneMap{
	"game": gameScene.NewGameScene,
}
