package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/setanarut/kamera/v2"
	"github.com/ubootgame/ubootgame/framework"
	ecsFramework "github.com/ubootgame/ubootgame/framework/ecs"
	"github.com/ubootgame/ubootgame/framework/input"
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
	"github.com/ubootgame/ubootgame/internal/systems/game_systems/camera"
	"github.com/ubootgame/ubootgame/internal/systems/game_systems/debug"
	"github.com/ubootgame/ubootgame/internal/systems/graphics"
	"github.com/ubootgame/ubootgame/internal/systems/weapons"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	devents "github.com/yohamta/donburi/features/events"
	"image/color"
)

type Scene struct {
	settings  framework.SettingsService[internal.Settings]
	resources framework.ResourceService
	display   framework.DisplayService

	ecs *ecs.ECS

	cursor *input.Cursor
}

func NewScene(settings framework.SettingsService[internal.Settings], resources framework.ResourceService, display framework.DisplayService) *Scene {
	scene := &Scene{
		settings:  settings,
		resources: resources,
		display:   display,
		ecs:       ecs.NewECS(donburi.NewWorld()),
		cursor:    input.NewCursor(),
	}

	ecsFramework.RegisterSystem(scene.ecs, debug.NewSystem(scene.settings, scene.ecs, scene.cursor))
	ecsFramework.RegisterSystem(scene.ecs, game_systems.NewInputSystem(scene.display, scene.cursor))
	ecsFramework.RegisterSystem(scene.ecs, camera.NewSystem(scene.settings, scene.ecs, scene.display))
	ecsFramework.RegisterSystem(scene.ecs, player.NewSystem(scene.settings, scene.ecs, scene.cursor))
	ecsFramework.RegisterSystem(scene.ecs, enemy.NewSystem())
	ecsFramework.RegisterSystem(scene.ecs, weapons.NewBulletSystem(scene.display))
	ecsFramework.RegisterSystem(scene.ecs, systems.NewPhysicsSystem())
	ecsFramework.RegisterSystem(scene.ecs, environment.NewWaterSystem(scene.settings))
	ecsFramework.RegisterSystem(scene.ecs, graphics.NewSpriteSystem(scene.settings, scene.display))
	ecsFramework.RegisterSystem(scene.ecs, graphics.NewAnimatedSpriteSystem(scene.settings))

	return scene
}

func (scene *Scene) Load() error {
	if err := scene.resources.RegisterResources(assets.Assets); err != nil {
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
	_ = actors.CreatePlayer(scene.resources, scene.ecs, assets.Battleship, display.HScale(0.1), space)
	_ = []*donburi.Entry{
		actors.CreateEnemy(scene.resources, scene.ecs, assets.Submarine, display.HScale(0.1), cp.Vector{X: -0.5, Y: 0.2}, cp.Vector{X: 0.1}, space),
		actors.CreateEnemy(scene.resources, scene.ecs, assets.Submarine, display.HScale(0.1), cp.Vector{X: 0.4, Y: 0.1}, cp.Vector{X: -0.05}, space),
	}

	return nil
}

func (scene *Scene) Update() error {
	devents.ProcessAllEvents(scene.ecs.World)

	scene.ecs.Update()

	return nil
}

func (scene *Scene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 4, G: 0, B: 43, A: 255})

	scene.ecs.DrawLayer(layers.Game, screen)
	if scene.settings.Settings().Debug.Enabled {
		scene.ecs.DrawLayer(layers.Debug, screen)
	}
}
