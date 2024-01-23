package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/resolv"
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"sync"
)

type Scene struct {
	ecs              *ecs.ECS
	once             sync.Once
	resourceRegistry *resources.Registry
}

func NewGameScene(resourceRegistry *resources.Registry) *Scene {
	return &Scene{resourceRegistry: resourceRegistry}
}

func (scene *Scene) Assets() *resources.Library {
	return assets.Assets
}

func (scene *Scene) Update() {
	scene.once.Do(scene.setup)
	scene.ecs.Update()
}

func (scene *Scene) Draw(screen *ebiten.Image) {
	scene.ecs.Draw(screen)
}

func (scene *Scene) setup() {
	world := donburi.NewWorld()
	scene.ecs = ecs.NewECS(world)

	scene.ecs.AddSystem(systems.UpdateShip)
	scene.ecs.AddSystem(systems.UpdateObjects)
	scene.ecs.AddRenderer(layers.Water, systems.DrawWater)
	scene.ecs.AddRenderer(layers.Foreground, systems.DrawShip)
	scene.ecs.AddRenderer(layers.Hud, systems.DrawDebug)

	//gw, gh := float64(config.C.Width), float64(config.C.Height)

	_ = entities.CreateWater(scene.ecs, scene.resourceRegistry)

	space := entities.CreateSpace(scene.ecs)

	resolv.Add(space,
		entities.CreateShip(scene.ecs, scene.resourceRegistry))
}
