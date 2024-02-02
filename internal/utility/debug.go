package utility

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ubootgame/ubootgame/internal/scenes/game/events"
	"github.com/yohamta/donburi"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/opentype"
	"log"
)

const dpi = 72

var debugFont *opentype.Font
var defaultFontSize = 14.0

func init() {
	var err error
	debugFont, err = opentype.Parse(gobold.TTF)
	if err != nil {
		log.Fatal(err)
	}
}

type debug struct {
	fontScale float64
	fontFace  font.Face
}

func (d *debug) UpdateFontFace(_ donburi.World, event events.DisplayUpdatedEventData) {
	if d.fontScale != event.ScalingFactor || d.fontFace == nil {
		var err error

		d.fontScale = event.ScalingFactor
		d.fontFace, err = opentype.NewFace(debugFont, &opentype.FaceOptions{
			Size:    defaultFontSize * event.ScalingFactor,
			DPI:     dpi,
			Hinting: font.HintingVertical,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (d *debug) PrintDebugTextAt(screen *ebiten.Image, debugText string, opts *ebiten.DrawImageOptions) {
	DrawTextAtWithOptions(screen, debugText, d.fontFace, opts)
}

var Debug = &debug{
	fontScale: 1.0,
}
