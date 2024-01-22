package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/ubootgame/ubootgame/internal/entities"
	"github.com/ubootgame/ubootgame/internal/resolv"
	"github.com/ubootgame/ubootgame/internal/systems"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"sync"
)

type GameScene struct {
	ecs            *ecs.ECS
	once           sync.Once
	resourceLoader *resource.Loader
}

func NewGameScene(resourceLoader *resource.Loader) *GameScene {
	return &GameScene{resourceLoader: resourceLoader}
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

	scene.ecs.AddSystem(systems.UpdateShip)
	scene.ecs.AddSystem(systems.UpdateObjects)
	scene.ecs.AddRenderer(ecs.LayerDefault, systems.DrawShip)
	scene.ecs.AddRenderer(ecs.LayerDefault, systems.DrawDebug)

	//gw, gh := float64(config.C.Width), float64(config.C.Height)
	space := entities.CreateSpace(scene.ecs)

	resolv.Add(space,
		entities.CreateShip(scene.ecs, scene.resourceLoader))
}
