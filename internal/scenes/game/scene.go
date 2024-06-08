package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/samber/lo"
	"github.com/ubootgame/ubootgame/internal/config"
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	gameSystemComponents "github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	actorEntities "github.com/ubootgame/ubootgame/internal/scenes/game/entities/actors"
	gameSystemEntities "github.com/ubootgame/ubootgame/internal/scenes/game/entities/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/scene_graph"
	"github.com/ubootgame/ubootgame/internal/scenes/game/events"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems"
	actorSystems "github.com/ubootgame/ubootgame/internal/scenes/game/systems/actors"
	environmentSystems "github.com/ubootgame/ubootgame/internal/scenes/game/systems/environment"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/weapons"
	"github.com/ubootgame/ubootgame/internal/scenes/game/tags"
	"github.com/ubootgame/ubootgame/internal/utility"
	ecsSystems "github.com/ubootgame/ubootgame/internal/utility/ecs/systems"
	"github.com/ubootgame/ubootgame/internal/utility/resources"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	devents "github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/transform"
	"gonum.org/v1/gonum/spatial/r2"
	"image/color"
	"sync"
)

type System interface {
	Update(e *ecs.ECS)
	Layers() []lo.Tuple2[ecs.LayerID, ecsSystems.Renderer]
}

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

	scene.ecs.DrawLayer(layers.Game, screen)
	if debug.Enabled {
		scene.ecs.DrawLayer(layers.Debug, screen)
	}
}

func (scene *Scene) setup() {
	// Game system components
	scene.ecs.World.Entry(scene.ecs.World.Create(gameSystemComponents.Camera))
	scene.ecs.World.Entry(scene.ecs.World.Create(gameSystemComponents.Display))
	scene.ecs.World.Entry(scene.ecs.World.Create(gameSystemComponents.Cursor))

	debugEntry := scene.ecs.World.Entry(scene.ecs.World.Create(gameSystemComponents.Debug))
	gameSystemComponents.Debug.SetValue(debugEntry, gameSystemComponents.DebugData{
		Enabled:        config.C.Debug,
		DrawGrid:       true,
		DrawCollisions: true,
		DrawPositions:  true,
		FontScale:      1.0,
	})

	// Systems
	debugSystem := game_system.NewDebugSystem()
	displaySystem := game_system.NewDisplaySystem()

	scene.registerSystem(debugSystem)
	scene.registerSystem(displaySystem)
	scene.registerSystem(game_system.NewCameraSystem())
	scene.registerSystem(game_system.NewCursorSystem())
	scene.registerSystem(actorSystems.NewPlayerSystem())
	scene.registerSystem(weapons.NewBulletSystem())
	scene.registerSystem(systems.NewMovementSystem())
	scene.registerSystem(systems.NewResolvSystem())
	scene.registerSystem(environmentSystems.NewWaterSystem())
	scene.registerSystem(visuals.NewSpriteSystem())
	scene.registerSystem(visuals.NewAnimatedSpriteSystem())

	// Events
	events.DisplayUpdatedEvent.Subscribe(scene.ecs.World, displaySystem.UpdateDisplay)
	events.DisplayUpdatedEvent.Subscribe(scene.ecs.World, debugSystem.UpdateFontFace)

	// Environment
	//water := environmentEntities.CreateWater(scene.ecs, scene.resourceRegistry, utility.VScale(0.1))
	//animatedWater := environmentEntities.CreateAnimatedWater(scene.ecs, scene.resourceRegistry, utility.HScale(0.2), r2.Vec{})

	environment := scene_graph.CreateSceneGroup(scene.ecs, tags.EnvironmentTag)
	//transform.AppendChild(environment, water, false)
	//transform.AppendChild(environment, animatedWater, false)

	// Objects
	player := actorEntities.CreatePlayer(scene.ecs, scene.resourceRegistry, utility.HScale(100))
	enemies := []*donburi.Entry{
		actorEntities.CreateEnemy(scene.ecs, scene.resourceRegistry, utility.HScale(100), r2.Vec{X: -700, Y: 300}, r2.Vec{X: 100}),
		actorEntities.CreateEnemy(scene.ecs, scene.resourceRegistry, utility.HScale(100), r2.Vec{X: 800, Y: 150}, r2.Vec{X: -50}),
	}

	objects := scene_graph.CreateSceneGroup(scene.ecs, tags.ObjectsTag)
	transform.AppendChild(objects, player, true)
	for _, enemy := range enemies {
		transform.AppendChild(objects, enemy, true)
	}

	// Projectiles
	projectiles := scene_graph.CreateSceneGroup(scene.ecs, tags.ProjectilesTag)

	// Scene graph
	sceneGraph := scene_graph.CreateSceneGraph(scene.ecs)
	transform.AppendChild(sceneGraph, environment, false)
	transform.AppendChild(sceneGraph, objects, false)
	transform.AppendChild(sceneGraph, projectiles, false)

	camera := gameSystemEntities.CreateCamera(scene.ecs)
	transform.AppendChild(sceneGraph, camera, false)
}

func (scene *Scene) registerSystem(system System) {
	scene.ecs.AddSystem(system.Update)
	for _, layer := range system.Layers() {
		scene.ecs.AddRenderer(layer.A, layer.B)
	}
}
