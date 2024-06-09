package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/samber/lo"
	"github.com/solarlune/resolv"
	"github.com/ubootgame/ubootgame/internal"
	"github.com/ubootgame/ubootgame/internal/framework"
	"github.com/ubootgame/ubootgame/internal/framework/draw"
	ecsFramework "github.com/ubootgame/ubootgame/internal/framework/ecs"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components/geometry"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities/actors"
	"github.com/ubootgame/ubootgame/internal/scenes/game/layers"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
	"github.com/yohamta/donburi/filter"
	"gonum.org/v1/gonum/spatial/r2"
	"image/color"
)

type CollisionSystem struct {
	ecsFramework.System

	settings *internal.Settings

	camera          *framework.Camera
	cursor          *framework.Cursor
	playerTransform *transform.TransformData

	query *donburi.Query
}

func NewCollisionSystem(settings *internal.Settings, cursor *framework.Cursor, camera *framework.Camera) *CollisionSystem {
	system := &CollisionSystem{
		settings: settings,
		cursor:   cursor,
		camera:   camera,
		query:    donburi.NewQuery(filter.Contains(transform.Transform, geometry.Bounds, geometry.Scale)),
	}
	system.Injector = ecsFramework.NewInjector([]ecsFramework.Injection{
		ecsFramework.WithTag(actors.PlayerTag, []ecsFramework.Injection{
			ecsFramework.Component(&system.playerTransform, transform.Transform),
		}),
	})
	return system
}

func (system *CollisionSystem) Layers() []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer] {
	return []lo.Tuple2[ecs.LayerID, ecsFramework.Renderer]{
		{A: layers.Debug, B: system.DrawDebug},
	}
}

func (system *CollisionSystem) Update(e *ecs.ECS) {
	system.System.Update(e)

	system.query.Each(e.World, func(entry *donburi.Entry) {
		bounds := geometry.Bounds.Get(entry)
		scale := geometry.Scale.Get(entry)

		worldScale := transform.WorldScale(entry)
		worldPosition := transform.WorldPosition(entry)
		worldRotation := transform.WorldRotation(entry)

		bounds.SetRotation(resolv.ToRadians(worldRotation))
		bounds.SetScale(worldScale.X, worldScale.Y)
		bounds.SetPosition(worldPosition.X-(scale.NormalizedSize.X*worldScale.X)/2, worldPosition.Y-(scale.NormalizedSize.Y*worldScale.Y)/2)
	})
}

func (system *CollisionSystem) DrawDebug(e *ecs.ECS, screen *ebiten.Image) {
	if !system.settings.Debug.DrawCollisions {
		return
	}

	player, _ := actors.PlayerTag.First(e.World)

	playerWorld := transform.WorldPosition(player)
	playerScreen := system.camera.WorldToScreenPosition(r2.Vec(playerWorld))

	line := resolv.NewLine(playerWorld.X, playerWorld.Y, system.cursor.WorldPosition.X, system.cursor.WorldPosition.Y)

	intersectionPoints := make([]resolv.Vector, 0)
	lineColor := color.RGBA{R: 255, G: 255, A: 255}

	geometry.Bounds.Each(e.World, func(shapeEntry *donburi.Entry) {
		bounds := geometry.Bounds.Get(shapeEntry)

		if intersection := line.Intersection(0, 0, bounds); intersection != nil {
			intersectionPoints = append(intersectionPoints, intersection.Points...)
			lineColor = color.RGBA{R: 255, A: 255}
		}

		drawPolygon(screen, system.camera, bounds, color.White)
	})

	l := line.Lines()[0]

	lineStart := system.camera.WorldToScreenPosition(r2.Vec(l.Start))
	lineEnd := system.camera.WorldToScreenPosition(r2.Vec(l.End))

	vector.StrokeLine(screen, float32(lineStart.X), float32(lineStart.Y), float32(lineEnd.X), float32(lineEnd.Y), 2, lineColor, true)

	draw.Dot(screen, playerScreen, lineColor)

	for _, point := range intersectionPoints {
		pointScreen := system.camera.WorldToScreenPosition(r2.Vec{X: point.X, Y: point.Y})
		draw.Dot(screen, pointScreen, color.RGBA{G: 255, A: 255})
	}
}

func drawPolygon(screen *ebiten.Image, camera *framework.Camera, shape *resolv.ConvexPolygon, color color.Color) {
	vertices := shape.Transformed()
	for i := 0; i < len(vertices); i++ {
		vert := vertices[i]
		vertScreen := camera.WorldToScreenPosition(r2.Vec(vert))

		next := vertices[0]
		if i < len(vertices)-1 {
			next = vertices[i+1]
		}
		nextScreen := camera.WorldToScreenPosition(r2.Vec(next))

		vector.StrokeLine(screen, float32(vertScreen.X), float32(vertScreen.Y), float32(nextScreen.X), float32(nextScreen.Y), 1, color, true)
	}
}
