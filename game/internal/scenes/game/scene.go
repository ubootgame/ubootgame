package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/setanarut/kamera/v2"
	"github.com/ubootgame/ubootgame/framework"
	"github.com/ubootgame/ubootgame/framework/game/ecs"
	"github.com/ubootgame/ubootgame/framework/game/input"
	"github.com/ubootgame/ubootgame/framework/services/display"
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
	cameraSystem "github.com/ubootgame/ubootgame/internal/systems/game_systems/camera"
	"github.com/ubootgame/ubootgame/internal/systems/game_systems/debug"
	"github.com/ubootgame/ubootgame/internal/systems/graphics"
	"github.com/ubootgame/ubootgame/internal/systems/weapons"
	"github.com/yohamta/donburi"
	devents "github.com/yohamta/donburi/features/events"
	"image/color"
)

type Scene struct {
	*ecs.Scene[internal.Settings]

	display framework.DisplayService

	cursor *input.Cursor
}

func NewScene(settings framework.SettingsService[internal.Settings], display framework.DisplayService, resources framework.ResourceService) *Scene {
	return &Scene{
		Scene:   ecs.NewECSScene(settings, resources, assets.Assets),
		display: display,
		cursor:  input.NewCursor(),
	}
}

func (scene *Scene) Load() error {
	if err := scene.Scene.Load(); err != nil {
		return err
	}

	virtualResolution := scene.display.VirtualResolution()

	cameraEntry := scene.ECS.World.Entry(scene.ECS.World.Create(components.Camera))
	components.Camera.SetValue(cameraEntry, components.CameraData{
		Camera:        kamera.NewCamera(0, 0, virtualResolution.X, virtualResolution.Y),
		MoveSpeed:     500,
		RotationSpeed: 2,
		ZoomSpeed:     10,
		MinZoom:       -100,
		MaxZoom:       100,
	})

	spaceEntry := scene.ECS.World.Entry(scene.ECS.World.Create(physics.Space))
	space := cp.NewSpace()
	physics.Space.Set(spaceEntry, space)

	// Systems
	scene.RegisterSystem(debug.NewSystem(scene.Settings, scene.ECS, scene.cursor))
	scene.RegisterSystem(game_systems.NewInputSystem(scene.display, scene.cursor))
	scene.RegisterSystem(cameraSystem.NewSystem(scene.Settings, scene.ECS, scene.display))
	scene.RegisterSystem(player.NewSystem(scene.Settings, scene.ECS, scene.cursor))
	scene.RegisterSystem(enemy.NewSystem())
	scene.RegisterSystem(weapons.NewBulletSystem(scene.display))
	scene.RegisterSystem(systems.NewPhysicsSystem())
	scene.RegisterSystem(environment.NewWaterSystem(scene.Settings))
	scene.RegisterSystem(graphics.NewSpriteSystem(scene.Settings, scene.display))
	scene.RegisterSystem(graphics.NewAnimatedSpriteSystem(scene.Settings))

	// Objects
	_ = actors.CreatePlayer(scene.Resources, scene.ECS, assets.Battleship, display.HScale(0.1), space)
	_ = []*donburi.Entry{
		actors.CreateEnemy(scene.Resources, scene.ECS, assets.Submarine, display.HScale(0.1), cp.Vector{X: -0.5, Y: 0.2}, cp.Vector{X: 0.1}, space),
		actors.CreateEnemy(scene.Resources, scene.ECS, assets.Submarine, display.HScale(0.1), cp.Vector{X: 0.4, Y: 0.1}, cp.Vector{X: -0.05}, space),
	}

	return nil
}

func (scene *Scene) Update() error {
	if err := scene.Scene.Update(); err != nil {
		return err
	}

	devents.ProcessAllEvents(scene.ECS.World)

	scene.ECS.Update()

	return nil
}

func (scene *Scene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 4, G: 0, B: 43, A: 255})

	scene.ECS.DrawLayer(layers.Game, screen)
	if scene.Settings.Settings().Debug.Enabled {
		scene.ECS.DrawLayer(layers.Debug, screen)
	}
}
