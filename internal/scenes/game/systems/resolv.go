package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/samber/lo"
	"github.com/solarlune/resolv"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/game_system"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/actors"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/ubootgame/ubootgame/internal/utility/draw"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/injector"
	"github.com/ubootgame/ubootgame/internal/utility/ecs/systems"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"gonum.org/v1/gonum/spatial/r2"
	"image/color"
)

type ResolvSystem struct {
	systems.BaseSystem

	debug           *game_system.DebugData
	camera          *game_system.CameraData
	display         *game_system.DisplayData
	cursor          *game_system.CursorData
	playerTransform *transform.TransformData

	query *donburi.Query
}

func NewResolvSystem() *ResolvSystem {
	system := &ResolvSystem{
		query: donburi.NewQuery(filter.Contains(transform.Transform, geometry.Shape)),
	}
	system.Injector = injector.NewInjector([]injector.Injection{
		injector.Once([]injector.Injection{
			injector.Component(&system.debug, game_system.Debug),
			injector.Component(&system.camera, game_system.Camera),
			injector.Component(&system.display, game_system.Display),
			injector.Component(&system.cursor, game_system.Cursor),
		}),
		injector.WithTag(actors.PlayerTag, []injector.Injection{
			injector.Component(&system.playerTransform, transform.Transform),
		}),
	})
	return system
}

func (system *ResolvSystem) Layers() []lo.Tuple2[ecs.LayerID, systems.Renderer] {
	return []lo.Tuple2[ecs.LayerID, systems.Renderer]{
		{A: layers.Debug, B: system.DrawDebug},
	}
}

func (system *ResolvSystem) Update(e *ecs.ECS) {
	system.BaseSystem.Update(e)

	system.query.Each(e.World, func(entry *donburi.Entry) {
		t := transform.Transform.Get(entry)
		shape := geometry.Shape.Get(entry)

		worldPosition := transform.WorldPosition(entry)
		shape.SetPosition(worldPosition.X, worldPosition.Y)
		shape.SetScale(t.LocalScale.X, t.LocalScale.Y)
		shape.SetRotation(resolv.ToRadians(t.LocalRotation))
	})
}

func (system *ResolvSystem) DrawDebug(e *ecs.ECS, screen *ebiten.Image) {
	if !system.debug.DrawCollisions {
		return
	}

	player, _ := actors.PlayerTag.First(e.World)

	playerWorld := transform.WorldPosition(player)

	playerScreen := system.camera.WorldToScreenPosition(r2.Vec(playerWorld))

	line := resolv.NewLine(playerScreen.X, playerScreen.Y, system.cursor.ScreenPosition.X, system.cursor.ScreenPosition.Y)

	intersectionPoints := make([]resolv.Vector, 0)
	lineColor := color.RGBA{R: 255, G: 255, A: 255}

	geometry.Shape.Each(e.World, func(shapeEntry *donburi.Entry) {
		shape := geometry.Shape.Get(shapeEntry)

		if intersection := line.Intersection(0, 0, shape); intersection != nil {
			intersectionPoints = append(intersectionPoints, intersection.Points...)
			lineColor = color.RGBA{R: 255, A: 255}
		}

		drawPolygon(screen, shape, color.White)
	})

	l := line.Lines()[0]

	vector.StrokeLine(screen, float32(l.Start.X), float32(l.Start.Y), float32(l.End.X), float32(l.End.Y), 2, lineColor, true)

	draw.BigDot(screen, playerScreen, lineColor)

	for _, point := range intersectionPoints {
		pointVec := r2.Vec{X: point.X, Y: point.Y}
		draw.BigDot(screen, pointVec, color.RGBA{G: 255, A: 255})
	}
}

func drawPolygon(screen *ebiten.Image, shape *resolv.ConvexPolygon, color color.Color) {
	vertices := shape.Transformed()
	for i := 0; i < len(vertices); i++ {
		vert := vertices[i]
		next := vertices[0]

		if i < len(vertices)-1 {
			next = vertices[i+1]
		}
		vector.StrokeLine(screen, float32(vert.X), float32(vert.Y), float32(next.X), float32(next.Y), 1, color, true)
	}
}
