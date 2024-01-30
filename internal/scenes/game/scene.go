package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/resolv"
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"gonum.org/v1/gonum/spatial/r2"
	"image/color"
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
	screen.Fill(color.RGBA{R: 4, G: 0, B: 43, A: 255})
	scene.ecs.Draw(screen)
}

func (scene *Scene) setup() {
	world := donburi.NewWorld()
	scene.ecs = ecs.NewECS(world)

	debugEntry := scene.ecs.World.Entry(scene.ecs.World.Create(components.Debug))
	components.Debug.SetValue(debugEntry, components.DebugData{
		DrawResolvLines: config.C.Debug,
		DrawGrid:        config.C.Debug,
		DrawPositions:   config.C.Debug,
	})

	cameraEntry := scene.ecs.World.Entry(scene.ecs.World.Create(components.Camera))
	components.Camera.SetValue(cameraEntry, components.CameraData{
		Position:   r2.Vec{X: 0.0, Y: 0.0},
		ZoomFactor: 1.0,
	})

	// Update systems
	scene.ecs.AddSystem(systems.UpdateCamera)
	scene.ecs.AddSystem(systems.Debug.Update)
	scene.ecs.AddSystem(systems.Objects.Update)
	scene.ecs.AddSystem(systems.UpdateMovement)
	scene.ecs.AddSystem(systems.UpdatePlayer)
	scene.ecs.AddSystem(systems.Sprites.Update)
	scene.ecs.AddSystem(systems.UpdateAseprites)

	// Draw systems
	scene.ecs.AddRenderer(layers.Water, systems.DrawWater)
	scene.ecs.AddRenderer(layers.Water, systems.DrawAnimatedWater)
	scene.ecs.AddRenderer(layers.Foreground, systems.Sprites.Draw)
	scene.ecs.AddRenderer(layers.Hud, systems.Objects.Draw)
	scene.ecs.AddRenderer(layers.Hud, systems.Debug.Draw)

	_ = entities.CreateWater(scene.ecs, scene.resourceRegistry)
	_ = entities.CreateAnimatedWater(scene.ecs, scene.resourceRegistry)

	space := entities.CreateSpace(scene.ecs)

	resolv.Add(space,
		entities.CreatePlayer(scene.ecs, scene.resourceRegistry, utility.HScaler(0.1)),
		entities.CreateEnemy(scene.ecs, scene.resourceRegistry, utility.HScaler(0.1), r2.Vec{X: -0.7, Y: 0.05}, r2.Vec{X: 0.001}),
		entities.CreateEnemy(scene.ecs, scene.resourceRegistry, utility.HScaler(-0.1), r2.Vec{X: 0.8, Y: 0.2}, r2.Vec{X: -0.001}),
	)
}
