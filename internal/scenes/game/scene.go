package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/actors"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/scene_graph"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/actors/player"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/environment"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/game_systems"
	cameraSystem "github.com/ubootgame/ubootgame/internal/scenes/game/systems/game_systems/camera"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/game_systems/debug"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/visuals"
	"github.com/ubootgame/ubootgame/internal/scenes/game/systems/weapons"
	"github.com/ubootgame/ubootgame/internal/scenes/game/tags"
	"github.com/ubootgame/ubootgame/pkg"
	"github.com/ubootgame/ubootgame/pkg/camera"
	"github.com/ubootgame/ubootgame/pkg/input"
	"github.com/ubootgame/ubootgame/pkg/services/scenes"
	"github.com/ubootgame/ubootgame/pkg/world"
	"github.com/yohamta/donburi"
	devents "github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/features/transform"
	"gonum.org/v1/gonum/spatial/r2"
	"image/color"
)

type Scene struct {
	*scenes.ECSScene[internal.Settings]

	display pkg.DisplayService

	camera *camera.Camera
	cursor *input.Cursor
}

func NewScene(settings pkg.SettingsService[internal.Settings], display pkg.DisplayService, resources pkg.ResourceService) *Scene {
	return &Scene{
		ECSScene: scenes.NewECSScene(settings, resources, assets.Assets),
		display:  display,
		camera:   camera.NewCamera(display),
		cursor:   input.NewCursor(),
	}
}

func (scene *Scene) Load() error {
	if err := scene.ECSScene.Load(); err != nil {
		return err
	}

	// Systems
	scene.RegisterSystem(debug.NewDebugSystem(scene.Settings, scene.ECS, scene.cursor, scene.camera))
	scene.RegisterSystem(game_systems.NewInputSystem(scene.cursor, scene.camera))
	scene.RegisterSystem(cameraSystem.NewCameraSystem(scene.Settings, scene.ECS, scene.camera))
	scene.RegisterSystem(player.NewPlayerSystem(scene.Settings, scene.ECS, scene.cursor))
	scene.RegisterSystem(weapons.NewBulletSystem(scene.camera))
	scene.RegisterSystem(systems.NewMovementSystem(scene.Settings))
	scene.RegisterSystem(systems.NewCollisionSystem(scene.Settings, scene.cursor, scene.camera))
	scene.RegisterSystem(environment.NewWaterSystem(scene.Settings))
	scene.RegisterSystem(visuals.NewSpriteSystem(scene.Settings, scene.camera))
	scene.RegisterSystem(visuals.NewAnimatedSpriteSystem(scene.Settings))

	// Environment
	//water := environmentEntities.CreateWater(scene.ecs, scene.resourceRegistry, utility.VScale(0.1))
	//animatedWater := environmentEntities.CreateAnimatedWater(scene.ecs, scene.resourceRegistry, utility.HScale(0.2), r2.Vec{})

	environmentGroup := scene_graph.CreateSceneGroup(scene.ECS, tags.EnvironmentTag)
	//transform.AppendChild(environmentGroup, water, false)
	//transform.AppendChild(environmentGroup, animatedWater, false)

	// Objects
	playerEntry := actors.CreatePlayer(scene.Resources, scene.ECS, world.HScale(100))
	enemyEntries := []*donburi.Entry{
		actors.CreateEnemy(scene.Resources, scene.ECS, world.HScale(100), r2.Vec{X: -700, Y: 300}, r2.Vec{X: 100}),
		actors.CreateEnemy(scene.Resources, scene.ECS, world.HScale(100), r2.Vec{X: 800, Y: 150}, r2.Vec{X: -50}),
	}

	objectsGroup := scene_graph.CreateSceneGroup(scene.ECS, tags.ObjectsTag)
	transform.AppendChild(objectsGroup, playerEntry, true)
	for _, enemy := range enemyEntries {
		transform.AppendChild(objectsGroup, enemy, true)
	}

	// Projectiles
	projectilesGroup := scene_graph.CreateSceneGroup(scene.ECS, tags.ProjectilesTag)

	// Scene graph
	sceneGraph := scene_graph.CreateSceneGraph(scene.ECS)
	transform.AppendChild(sceneGraph, environmentGroup, false)
	transform.AppendChild(sceneGraph, objectsGroup, false)
	transform.AppendChild(sceneGraph, projectilesGroup, false)

	return nil
}

func (scene *Scene) Update() error {
	if err := scene.ECSScene.Update(); err != nil {
		return err
	}

	devents.ProcessAllEvents(scene.ECS.World)

	scene.ECS.Update()

	scene.camera.UpdateCameraMatrix()

	return nil
}

func (scene *Scene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 4, G: 0, B: 43, A: 255})

	scene.ECS.DrawLayer(layers.Game, screen)
	if scene.Settings.Settings().Debug.Enabled {
		scene.ECS.DrawLayer(layers.Debug, screen)
	}
}
