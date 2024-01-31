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
	"github.com/yohamta/donburi/features/events"
	"gonum.org/v1/gonum/spatial/r2"
	"image/color"
	"sync"
)

type Scene struct {
	ecs              *ecs.ECS
	once             sync.Once
	resourceRegistry *resources.Registry
	debugEntry       *donburi.Entry
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

func (scene *Scene) AdjustScreen(windowSize r2.Vec, virtualResolution r2.Vec) {
	systems.ScreenUpdatedEvent.Publish(scene.ecs.World, systems.ScreenUpdatedEventData{
		WindowSize:        windowSize,
		VirtualResolution: virtualResolution,
	})
}

func (scene *Scene) Update() {
	scene.once.Do(scene.setup)

	events.ProcessAllEvents(scene.ecs.World)

	scene.ecs.Update()
}

func (scene *Scene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 4, G: 0, B: 43, A: 255})

	if scene.debugEntry == nil {
		var ok bool
		if scene.debugEntry, ok = components.Debug.First(scene.ecs.World); !ok {
			panic("no debug found")
		}
	}

	debug := components.Debug.Get(scene.debugEntry)

	scene.ecs.DrawLayer(layers.Background, screen)
	scene.ecs.DrawLayer(layers.Foreground, screen)
	if debug.Enabled {
		scene.ecs.DrawLayer(layers.Debug, screen)
	}
}

func (scene *Scene) setup() {
	debugEntry := scene.ecs.World.Entry(scene.ecs.World.Create(components.Debug))
	components.Debug.SetValue(debugEntry, components.DebugData{
		Enabled:         config.C.Debug,
		DrawResolvLines: true,
		DrawGrid:        true,
		DrawPositions:   true,
	})

	_ = scene.ecs.World.Entry(scene.ecs.World.Create(components.Camera))

	displayEntry := scene.ecs.World.Entry(scene.ecs.World.Create(components.Display))
	display := components.Display.Get(displayEntry)

	// Update systems
	scene.ecs.AddSystem(systems.Camera.Update)
	scene.ecs.AddSystem(systems.Debug.Update)
	scene.ecs.AddSystem(systems.Movement.Update)
	scene.ecs.AddSystem(systems.Resolv.Update)
	scene.ecs.AddSystem(systems.Player.Update)
	scene.ecs.AddSystem(systems.Sprite.Update)
	scene.ecs.AddSystem(systems.Aseprite.Update)

	// Draw systems
	scene.ecs.AddRenderer(layers.Background, systems.Water.Draw)
	scene.ecs.AddRenderer(layers.Background, systems.AnimatedWater.Draw)
	scene.ecs.AddRenderer(layers.Foreground, systems.Sprite.Draw)
	scene.ecs.AddRenderer(layers.Debug, systems.Resolv.Draw)
	scene.ecs.AddRenderer(layers.Debug, systems.Debug.Draw)

	// Events
	systems.ScreenUpdatedEvent.Subscribe(scene.ecs.World, systems.Display.UpdateScreen)

	_ = entities.CreateWater(scene.ecs, scene.resourceRegistry)
	_ = entities.CreateAnimatedWater(scene.ecs, scene.resourceRegistry)

	space := entities.CreateSpace(scene.ecs, r2.Scale(2, display.VirtualResolution))

	resolv.Add(space,
		entities.CreatePlayer(scene.ecs, scene.resourceRegistry, utility.HScaler(0.1)),
		entities.CreateEnemy(scene.ecs, scene.resourceRegistry, utility.HScaler(0.1), r2.Vec{X: -0.7, Y: 0.05}, r2.Vec{X: 0.1}),
		entities.CreateEnemy(scene.ecs, scene.resourceRegistry, utility.HScaler(0.1), r2.Vec{X: 0.8, Y: 0.2}, r2.Vec{X: -0.05}),
	)
}
