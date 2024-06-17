package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/samber/do"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/framework/game"
	"github.com/ubootgame/ubootgame/framework/graphics/display"
	"github.com/ubootgame/ubootgame/framework/input"
	"github.com/ubootgame/ubootgame/framework/resources"
	"github.com/ubootgame/ubootgame/framework/settings"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/entities"
	"github.com/ubootgame/ubootgame/internal/entities/actors"
	"github.com/ubootgame/ubootgame/internal/layers"
	"github.com/ubootgame/ubootgame/internal/scenes/game/assets"
	"github.com/ubootgame/ubootgame/internal/systems"
	"github.com/ubootgame/ubootgame/internal/systems/actors/enemy"
	"github.com/ubootgame/ubootgame/internal/systems/actors/player"
	"github.com/ubootgame/ubootgame/internal/systems/environment"
	"github.com/ubootgame/ubootgame/internal/systems/game_systems"
	"github.com/ubootgame/ubootgame/internal/systems/game_systems/camera"
	"github.com/ubootgame/ubootgame/internal/systems/game_systems/debug"
	"github.com/ubootgame/ubootgame/internal/systems/graphics"
	"github.com/ubootgame/ubootgame/internal/systems/weapons"
	devents "github.com/yohamta/donburi/features/events"
	"image/color"
)

type gameScene struct {
	injector *do.Injector

	settingsProvider settings.Provider[internal.Settings]
	resourceRegistry resources.Registry
	input            input.Input
	ecs              ecsFramework.Service
}

func NewGameScene(i *do.Injector) game.Scene {
	scene := &gameScene{
		injector:         i,
		settingsProvider: do.MustInvoke[settings.Provider[internal.Settings]](i),
		resourceRegistry: do.MustInvoke[resources.Registry](i),
		input:            do.MustInvoke[input.Input](i),
		ecs:              do.MustInvoke[ecsFramework.Service](i),
	}

	scene.ecs.RegisterSystem(debug.NewDebugSystem)
	scene.ecs.RegisterSystem(game_systems.NewInputSystem)
	scene.ecs.RegisterSystem(camera.NewCameraSystem)
	scene.ecs.RegisterSystem(player.NewPlayerSystem)
	scene.ecs.RegisterSystem(enemy.NewEnemySystem)
	scene.ecs.RegisterSystem(weapons.NewBulletSystem)
	scene.ecs.RegisterSystem(systems.NewPhysicsSystem)
	scene.ecs.RegisterSystem(environment.NewWaterSystem)
	scene.ecs.RegisterSystem(graphics.NewSpriteSystem)
	scene.ecs.RegisterSystem(graphics.NewAnimatedSpriteSystem)

	return scene
}

func (scene *gameScene) Load() error {
	if err := scene.resourceRegistry.RegisterResources(assets.Assets); err != nil {
		return err
	}

	// Camera
	entities.CameraFactory.Spawn(scene.injector, entities.NewCameraParams{
		MoveSpeed:     500,
		RotationSpeed: 2,
		ZoomSpeed:     10,
		MinZoom:       -100,
		MaxZoom:       100,
	})

	// Physics
	entities.SpaceFactory.Spawn(scene.injector, entities.NewSpaceParams{})

	// Objects
	actors.PlayerFactory.Spawn(scene.injector, actors.NewPlayerParams{
		ImageID: assets.Battleship,
		Scale:   display.HScale(0.1),
	})

	actors.EnemyFactory.Spawn(scene.injector, actors.NewEnemyParams{
		ImageID:  assets.Submarine,
		Scale:    display.HScale(0.1),
		Position: cp.Vector{X: -0.5, Y: 0.2},
		Velocity: cp.Vector{X: 0.1},
	})
	actors.EnemyFactory.Spawn(scene.injector, actors.NewEnemyParams{
		ImageID:  assets.Submarine,
		Scale:    display.HScale(0.1),
		Position: cp.Vector{X: 0.4, Y: 0.1},
		Velocity: cp.Vector{X: -0.05},
	})

	return nil
}

func (scene *gameScene) Update() error {
	devents.ProcessAllEvents(scene.ecs.World())

	scene.ecs.ECS().Update()

	return nil
}

func (scene *gameScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 4, G: 0, B: 43, A: 255})

	scene.ecs.ECS().DrawLayer(layers.Game, screen)
	if scene.settingsProvider.Settings().Debug.Enabled {
		scene.ecs.ECS().DrawLayer(layers.Debug, screen)
	}
}
