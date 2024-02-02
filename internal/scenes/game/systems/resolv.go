package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	evector "github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/quartercastle/vector"
	"github.com/solarlune/resolv"
	"github.com/ubootgame/ubootgame/internal/scenes/game/components"
	"github.com/ubootgame/ubootgame/internal/scenes/game/entities"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"gonum.org/v1/gonum/spatial/r2"
	"image/color"
)

type resolvSystem struct {
	query                                              *donburi.Query
	debugEntry, cameraEntry, displayEntry, playerEntry *donburi.Entry
}

var Resolv = &resolvSystem{
	query: donburi.NewQuery(filter.Contains(components.Transform, components.Shape)),
}

func (system *resolvSystem) Update(e *ecs.ECS) {
	var ok bool
	if system.debugEntry == nil {
		if system.debugEntry, ok = components.Debug.First(e.World); !ok {
			panic("no debug found")
		}
	}
	if system.displayEntry == nil {
		if system.displayEntry, ok = components.Display.First(e.World); !ok {
			panic("no display found")
		}
	}
	if system.cameraEntry == nil {
		if system.cameraEntry, ok = components.Camera.First(e.World); !ok {
			panic("no camera found")
		}
	}
	if system.playerEntry == nil {
		if system.playerEntry, ok = entities.PlayerTag.First(e.World); !ok {
			panic("no player found")
		}
	}

	camera := components.Camera.Get(system.cameraEntry)
	display := components.Display.Get(system.displayEntry)

	system.query.Each(e.World, func(entry *donburi.Entry) {
		transform := components.Transform.Get(entry)
		shape := components.Shape.Get(entry)

		position := camera.WorldToScreenPosition(r2.Vec{X: transform.Center.X - transform.Size.X/2, Y: transform.Center.Y - transform.Size.Y/2})
		shape.SetPosition(position.X, position.Y)
		shape.SetScale(display.VirtualResolution.X*camera.ZoomFactor, display.VirtualResolution.X*camera.ZoomFactor)
		shape.SetRotation(resolv.ToRadians(transform.Rotate))
	})
}

func (system *resolvSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	debug := components.Debug.Get(system.debugEntry)

	if debug.Enabled && debug.DrawCollisions {
		system.drawDebug(e, screen)
	}
}

func (system *resolvSystem) drawDebug(e *ecs.ECS, screen *ebiten.Image) {
	transform := components.Transform.Get(system.playerEntry)
	camera := components.Camera.Get(system.cameraEntry)

	mx, my := ebiten.CursorPosition()
	mousePosition := camera.ScreenToWorldPosition(r2.Vec{X: float64(mx), Y: float64(my)})
	mouseScreen := camera.WorldToScreenPosition(mousePosition)

	playerScreen := camera.WorldToScreenPosition(transform.Center)

	line := resolv.NewLine(playerScreen.X, playerScreen.Y, mouseScreen.X, mouseScreen.Y)

	intersectionPoints := make([]vector.Vector, 0)
	lineColor := color.RGBA{R: 255, G: 255, A: 255}

	components.Shape.Each(e.World, func(shapeEntry *donburi.Entry) {
		shape := components.Shape.Get(shapeEntry)

		if intersection := line.Intersection(0, 0, shape); intersection != nil {
			intersectionPoints = append(intersectionPoints, intersection.Points...)
			lineColor = color.RGBA{R: 255, A: 255}
		}

		drawPolygon(screen, shape, color.White)
	})

	l := line.Lines()[0]

	evector.StrokeLine(screen, float32(l.Start.X()), float32(l.Start.Y()), float32(l.End.X()), float32(l.End.Y()), 2, lineColor, true)

	drawBigDot(screen, playerScreen, lineColor)

	for _, point := range intersectionPoints {
		pointVec := r2.Vec{X: point.X(), Y: point.Y()}
		drawBigDot(screen, pointVec, color.RGBA{G: 255, A: 255})
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
		evector.StrokeLine(screen, float32(vert.X()), float32(vert.Y()), float32(next.X()), float32(next.Y()), 1, color, false)
	}
}

var bigDotImg *ebiten.Image

func drawBigDot(screen *ebiten.Image, position r2.Vec, drawColor color.Color) {
	if bigDotImg == nil {
		bigDotImg = ebiten.NewImage(4, 4)
		bigDotImg.Fill(color.White)
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(position.X-2, position.Y-2)
	opt.ColorScale.ScaleWithColor(drawColor)
	screen.DrawImage(bigDotImg, opt)
}
