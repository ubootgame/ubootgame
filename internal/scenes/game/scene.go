package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	gameSystemComponents "github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	actorEntities "github.com/ubootgame/ubootgame/internal/scenes/game/entities/actors"
	environmentEntities "github.com/ubootgame/ubootgame/internal/scenes/game/entities/environment"
	"github.com/ubootgame/ubootgame/internal/scenes/game/events"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems"
	actorSystems "github.com/ubootgame/ubootgame/internal/scenes/game/systems/actors"
	environmentSystems "github.com/ubootgame/ubootgame/internal/scenes/game/systems/environment"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/weapons"
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

	if scene.debugEntry == nil {
		var ok bool
		if scene.debugEntry, ok = gameSystemComponents.Debug.First(scene.ecs.World); !ok {
			panic("no debug found")
		}
	}

	debug := gameSystemComponents.Debug.Get(scene.debugEntry)

	scene.ecs.DrawLayer(layers.Background, screen)
	scene.ecs.DrawLayer(layers.Foreground, screen)
	if debug.Enabled {
		scene.ecs.DrawLayer(layers.Debug, screen)
	}
}

func (scene *Scene) setup() {
	debugEntry := scene.ecs.World.Entry(scene.ecs.World.Create(gameSystemComponents.Debug))
	gameSystemComponents.Debug.SetValue(debugEntry, gameSystemComponents.DebugData{
		Enabled:        config.C.Debug,
		DrawGrid:       true,
		DrawCollisions: true,
		DrawPositions:  true,
	})

	_ = scene.ecs.World.Entry(scene.ecs.World.Create(gameSystemComponents.Camera))
	_ = scene.ecs.World.Entry(scene.ecs.World.Create(gameSystemComponents.Display))
	_ = scene.ecs.World.Entry(scene.ecs.World.Create(gameSystemComponents.Cursor))

	// Update systems
	scene.ecs.AddSystem(game_system.Cursor.Update)
	scene.ecs.AddSystem(game_system.Camera.Update)
	scene.ecs.AddSystem(game_system.Debug.Update)
	scene.ecs.AddSystem(systems.Movement.Update)
	scene.ecs.AddSystem(systems.Resolv.Update)
	scene.ecs.AddSystem(actorSystems.Player.Update)
	scene.ecs.AddSystem(visuals.Sprite.Update)
	scene.ecs.AddSystem(visuals.Aseprite.Update)
	scene.ecs.AddSystem(weapons.Bullet.Update)

	// Draw systems
	scene.ecs.AddRenderer(layers.Background, environmentSystems.Water.Draw)
	scene.ecs.AddRenderer(layers.Background, environmentSystems.AnimatedWater.Draw)
	scene.ecs.AddRenderer(layers.Foreground, visuals.Sprite.Draw)
	scene.ecs.AddRenderer(layers.Foreground, weapons.Bullet.Draw)
	scene.ecs.AddRenderer(layers.Debug, systems.Resolv.Draw)
	scene.ecs.AddRenderer(layers.Debug, game_system.Debug.Draw)

	// Events
	events.DisplayUpdatedEvent.Subscribe(scene.ecs.World, game_system.Display.UpdateDisplay)
	events.DisplayUpdatedEvent.Subscribe(scene.ecs.World, utility.Debug.UpdateFontFace)

	_ = environmentEntities.CreateWater(scene.ecs, scene.resourceRegistry)
	_ = environmentEntities.CreateAnimatedWater(scene.ecs, scene.resourceRegistry)

	actorEntities.CreatePlayer(scene.ecs, scene.resourceRegistry, utility.HScaler(0.1))
	actorEntities.CreateEnemy(scene.ecs, scene.resourceRegistry, utility.HScaler(0.1), r2.Vec{X: -0.7, Y: 0.05}, r2.Vec{X: 0.1})
	actorEntities.CreateEnemy(scene.ecs, scene.resourceRegistry, utility.HScaler(0.1), r2.Vec{X: 0.8, Y: 0.2}, r2.Vec{X: -0.05})
}
