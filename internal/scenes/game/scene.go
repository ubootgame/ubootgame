package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/framework"
	"github.com/ubootgame/ubootgame/internal/framework/coordinate_system"
	"github.com/ubootgame/ubootgame/internal/framework/resources"
	"github.com/ubootgame/ubootgame/internal/framework/scenes"
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	actorEntities "github.com/ubootgame/ubootgame/internal/scenes/game/entities/actors"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/scene_graph"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/actors/player"
	environmentSystems "github.com/ubootgame/ubootgame/internal/scenes/game/systems/environment"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/game_system/camera"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/game_system/debug"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/weapons"
	"github.com/ubootgame/ubootgame/internal/scenes/game/tags"
	"github.com/yohamta/donburi"
	devents "github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/transform"
	"gonum.org/v1/gonum/spatial/r2"
	"image/color"
)

type Scene struct {
	*scenes.ECSScene

	camera *framework.Camera
	cursor *framework.Cursor
}

func NewScene(settings *internal.Settings) *Scene {
	return &Scene{
		ECSScene: scenes.NewECSScene(settings, assets.Assets),
		camera:   framework.NewCamera(settings),
		cursor:   framework.NewCursor(),
	}
}

func (scene *Scene) Load(resourceRegistry *resources.Registry) error {
	if err := scene.ECSScene.Load(resourceRegistry); err != nil {
		return err
	}

	// Systems
	scene.RegisterSystem(debug.NewDebugSystem(scene.ECS, scene.Settings, scene.cursor, scene.camera))
	scene.RegisterSystem(game_system.NewInputSystem(scene.cursor, scene.camera))
	scene.RegisterSystem(camera.NewCameraSystem(scene.ECS, scene.Settings, scene.camera))
	scene.RegisterSystem(player.NewPlayerSystem(scene.ECS, scene.Settings, scene.cursor))
	scene.RegisterSystem(weapons.NewBulletSystem(scene.camera))
	scene.RegisterSystem(systems.NewMovementSystem(scene.Settings))
	scene.RegisterSystem(systems.NewCollisionSystem(scene.Settings, scene.cursor, scene.camera))
	scene.RegisterSystem(environmentSystems.NewWaterSystem(scene.Settings))
	scene.RegisterSystem(visuals.NewSpriteSystem(scene.Settings, scene.camera))
	scene.RegisterSystem(visuals.NewAnimatedSpriteSystem(scene.Settings))

	// Environment
	//water := environmentEntities.CreateWater(scene.ecs, scene.resourceRegistry, utility.VScale(0.1))
	//animatedWater := environmentEntities.CreateAnimatedWater(scene.ecs, scene.resourceRegistry, utility.HScale(0.2), r2.Vec{})

	environment := scene_graph.CreateSceneGroup(scene.ECS, tags.EnvironmentTag)
	//transform.AppendChild(environment, water, false)
	//transform.AppendChild(environment, animatedWater, false)

	// Objects
	p := actorEntities.CreatePlayer(scene.ECS, scene.ResourceRegistry, coordinate_system.HScale(100))
	enemies := []*donburi.Entry{
		actorEntities.CreateEnemy(scene.ECS, scene.ResourceRegistry, coordinate_system.HScale(100), r2.Vec{X: -700, Y: 300}, r2.Vec{X: 100}),
		actorEntities.CreateEnemy(scene.ECS, scene.ResourceRegistry, coordinate_system.HScale(100), r2.Vec{X: 800, Y: 150}, r2.Vec{X: -50}),
	}

	objects := scene_graph.CreateSceneGroup(scene.ECS, tags.ObjectsTag)
	transform.AppendChild(objects, p, true)
	for _, enemy := range enemies {
		transform.AppendChild(objects, enemy, true)
	}

	// Projectiles
	projectiles := scene_graph.CreateSceneGroup(scene.ECS, tags.ProjectilesTag)

	// Scene graph
	sceneGraph := scene_graph.CreateSceneGraph(scene.ECS)
	transform.AppendChild(sceneGraph, environment, false)
	transform.AppendChild(sceneGraph, objects, false)
	transform.AppendChild(sceneGraph, projectiles, false)

	return nil
}

func (scene *Scene) Update() error {
	scene.ECSScene.Update()

	devents.ProcessAllEvents(scene.ECS.World)

	scene.ECS.Update()

	scene.camera.UpdateCameraMatrix()

	return nil
}

func (scene *Scene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 4, G: 0, B: 43, A: 255})

	scene.ECS.DrawLayer(layers.Game, screen)
	if scene.Settings.Debug.Enabled {
		scene.ECS.DrawLayer(layers.Debug, screen)
	}
}
