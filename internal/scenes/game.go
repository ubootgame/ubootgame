package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/systems"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"sync"
)

type GameScene struct {
	ecs  *ecs.ECS
	once sync.Once
}

func NewGameScene() *GameScene {
	return &GameScene{}
}

func (scene *GameScene) Update() {
	scene.once.Do(scene.setup)
	scene.ecs.Update()
}

func (scene *GameScene) Draw(screen *ebiten.Image) {
	scene.ecs.DrawLayer(ecs.LayerDefault, screen)
}

func (scene *GameScene) setup() {
	world := donburi.NewWorld()
	scene.ecs = ecs.NewECS(world)

	scene.ecs.AddRenderer(ecs.LayerDefault, systems.DrawDebug)
}
