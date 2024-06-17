package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/samber/do"
	"github.com/setanarut/kamera/v2"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/framework/game"
	"github.com/ubootgame/ubootgame/framework/graphics/display"
	"github.com/ubootgame/ubootgame/framework/input"
	"github.com/ubootgame/ubootgame/framework/resources"
	"github.com/ubootgame/ubootgame/framework/settings"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/components"
	"github.com/ubootgame/ubootgame/internal/components/physics"
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
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	devents "github.com/yohamta/donburi/features/events"
	"image/color"
)

type gameScene struct {
	injector *do.Injector

	settingsProvider settings.Provider[internal.Settings]
	resourceRegistry resources.Registry
	input            input.Input
	display          display.Display

	ecs *ecs.ECS
}

func NewGameScene(i *do.Injector) game.Scene {
	e := ecs.NewECS(donburi.NewWorld())
	scene := &gameScene{
		injector:         i,
		settingsProvider: do.MustInvoke[settings.Provider[internal.Settings]](i),
		resourceRegistry: do.MustInvoke[resources.Registry](i),
		input:            do.MustInvoke[input.Input](i),
		display:          do.MustInvoke[display.Display](i),
		ecs:              e,
	}

	ecsFramework.RegisterSystem(scene.ecs, debug.NewDebugSystem(i, scene.ecs))
	ecsFramework.RegisterSystem(scene.ecs, game_systems.NewInputSystem(i))
	ecsFramework.RegisterSystem(scene.ecs, camera.NewCameraSystem(i, scene.ecs))
	ecsFramework.RegisterSystem(scene.ecs, player.NewPlayerSystem(i, scene.ecs))
	ecsFramework.RegisterSystem(scene.ecs, enemy.NewSystem())
	ecsFramework.RegisterSystem(scene.ecs, weapons.NewBulletSystem(i))
	ecsFramework.RegisterSystem(scene.ecs, systems.NewPhysicsSystem())
	ecsFramework.RegisterSystem(scene.ecs, environment.NewWaterSystem(i))
	ecsFramework.RegisterSystem(scene.ecs, graphics.NewSpriteSystem(i))
	ecsFramework.RegisterSystem(scene.ecs, graphics.NewAnimatedSpriteSystem(i))

	return scene
}

func (scene *gameScene) Load() error {
	if err := scene.resourceRegistry.RegisterResources(assets.Assets); err != nil {
		return err
	}

	// Camera
	virtualResolution := scene.display.VirtualResolution()

	cameraEntry := scene.ecs.World.Entry(scene.ecs.World.Create(components.Camera))
	components.Camera.SetValue(cameraEntry, components.CameraData{
		Camera:        kamera.NewCamera(0, 0, virtualResolution.X, virtualResolution.Y),
		MoveSpeed:     500,
		RotationSpeed: 2,
		ZoomSpeed:     10,
		MinZoom:       -100,
		MaxZoom:       100,
	})

	// Physics
	spaceEntry := scene.ecs.World.Entry(scene.ecs.World.Create(physics.Space))
	space := cp.NewSpace()
	physics.Space.Set(spaceEntry, space)

	// Objects
	ecsFramework.Spawn(scene.injector, scene.ecs, actors.CreatePlayer, actors.NewPlayerParams{
		ImageID: assets.Battleship,
		Scale:   display.HScale(0.1),
		Space:   space,
	})

	ecsFramework.Spawn(scene.injector, scene.ecs, actors.CreateEnemy, actors.NewEnemyParams{
		ImageID:  assets.Submarine,
		Scale:    display.HScale(0.1),
		Position: cp.Vector{X: -0.5, Y: 0.2},
		Velocity: cp.Vector{X: 0.1},
		Space:    space,
	})
	ecsFramework.Spawn(scene.injector, scene.ecs, actors.CreateEnemy, actors.NewEnemyParams{
		ImageID:  assets.Submarine,
		Scale:    display.HScale(0.1),
		Position: cp.Vector{X: 0.4, Y: 0.1},
		Velocity: cp.Vector{X: -0.05},
		Space:    space,
	})

	return nil
}

func (scene *gameScene) Update() error {
	devents.ProcessAllEvents(scene.ecs.World)

	scene.ecs.Update()

	return nil
}

func (scene *gameScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 4, G: 0, B: 43, A: 255})

	scene.ecs.DrawLayer(layers.Game, screen)
	if scene.settingsProvider.Settings().Debug.Enabled {
		scene.ecs.DrawLayer(layers.Debug, screen)
	}
}
