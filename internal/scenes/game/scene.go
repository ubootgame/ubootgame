package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities"
	"github.com/ubootgame/ubootgame/internal/scenes/game/events"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems"
	"github.com/ubootgame/ubootgame/internal/utility"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	devents "github.com/yohamta/donburi/features/events"
	"gonum.org/v1/gonum/spatial/r2"
	"image/color"
	"sync"
)

type Scene struct {
	ecs              *ecs.ECS
	once             sync.Once
	resourceRegistry *resources.Registry
	debugEntry       *donburi.Entry
	menuEntry        *donburi.Entry
}

func NewGameScene(resourceRegistry *resources.Registry) *Scene {
	world := donburi.NewWorld()
	return &Scene{
		ecs:              ecs.NewECS(world),
		resourceRegistry: resourceRegistry,
	}
}

func (scene *Scene) Assets() *resources.Library {
	return assets.Assets
}

func (scene *Scene) AdjustScreen(windowSize r2.Vec, virtualResolution r2.Vec, scalingFactor float64) {
	events.DisplayUpdatedEvent.Publish(scene.ecs.World, events.DisplayUpdatedEventData{
		WindowSize:        windowSize,
		VirtualResolution: virtualResolution,
		ScalingFactor:     scalingFactor,
	})
}

func (scene *Scene) Update() {
	scene.once.Do(scene.setup)

	devents.ProcessAllEvents(scene.ecs.World)

	scene.ecs.Update()
}

func (scene *Scene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 4, G: 0, B: 43, A: 255})

	debug := components.Debug.Get(scene.debugEntry)
	menu := components.Menu.Get(scene.menuEntry)

	scene.ecs.DrawLayer(layers.Background, screen)
	scene.ecs.DrawLayer(layers.Foreground, screen)

	if menu.IsOpen {
		scene.ecs.DrawLayer(layers.Menu, screen)
	}

	if debug.Enabled {
		scene.ecs.DrawLayer(layers.Debug, screen)
	}
}

func (scene *Scene) setup() {
	debugEntry := scene.ecs.World.Entry(scene.ecs.World.Create(components.Debug))
	components.Debug.SetValue(debugEntry, components.DebugData{
		Enabled:        config.C.Debug,
		DrawGrid:       true,
		DrawCollisions: true,
		DrawPositions:  true,
	})
	scene.debugEntry = debugEntry

	menuEntry := scene.ecs.World.Entry(scene.ecs.World.Create(components.Menu))
	components.Menu.SetValue(menuEntry, components.MenuData{})
	scene.menuEntry = menuEntry

	_ = scene.ecs.World.Entry(scene.ecs.World.Create(components.Camera))
	_ = scene.ecs.World.Entry(scene.ecs.World.Create(components.Display))

	// Update systems
	scene.ecs.AddSystem(systems.Camera.Update)
	scene.ecs.AddSystem(systems.Debug.Update)
	scene.ecs.AddSystem(systems.Movement.Update)
	scene.ecs.AddSystem(systems.Resolv.Update)
	scene.ecs.AddSystem(systems.Player.Update)
	scene.ecs.AddSystem(systems.Sprite.Update)
	scene.ecs.AddSystem(systems.Aseprite.Update)
	scene.ecs.AddSystem(systems.Bullet.Update)
	scene.ecs.AddSystem(systems.Menu.Update)

	// Draw systems
	scene.ecs.AddRenderer(layers.Background, systems.Water.Draw)
	scene.ecs.AddRenderer(layers.Background, systems.AnimatedWater.Draw)
	scene.ecs.AddRenderer(layers.Foreground, systems.Sprite.Draw)
	scene.ecs.AddRenderer(layers.Foreground, systems.Bullet.Draw)
	scene.ecs.AddRenderer(layers.Debug, systems.Resolv.Draw)
	scene.ecs.AddRenderer(layers.Debug, systems.Debug.Draw)
	scene.ecs.AddRenderer(layers.Menu, systems.Menu.Draw)

	// Events
	events.DisplayUpdatedEvent.Subscribe(scene.ecs.World, systems.Display.UpdateDisplay)
	events.DisplayUpdatedEvent.Subscribe(scene.ecs.World, utility.Debug.UpdateFontFace)

	_ = entities.CreateWater(scene.ecs, scene.resourceRegistry)
	_ = entities.CreateAnimatedWater(scene.ecs, scene.resourceRegistry)

	entities.CreatePlayer(scene.ecs, scene.resourceRegistry, utility.HScaler(0.1))
	entities.CreateEnemy(scene.ecs, scene.resourceRegistry, utility.HScaler(0.1), r2.Vec{X: -0.7, Y: 0.05}, r2.Vec{X: 0.1})
	entities.CreateEnemy(scene.ecs, scene.resourceRegistry, utility.HScaler(0.1), r2.Vec{X: 0.8, Y: 0.2}, r2.Vec{X: -0.05})
}
