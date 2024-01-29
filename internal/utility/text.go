package utility

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

func FormatBytes(b uint64) string {
	if b >= 1073741824 {
		return fmt.Sprintf("%0.2f GiB", float64(b)/1073741824)
	} else if b >= 1048576 {
		return fmt.Sprintf("%0.2f MiB", float64(b)/1048576)
	} else if b >= 1024 {
		return fmt.Sprintf("%0.2f KiB", float64(b)/1024)
	} else {
		return fmt.Sprintf("%d B", b)
	}
}

func DrawTextAtWithOptions(screen *ebiten.Image, s string, f font.Face, opts *ebiten.DrawImageOptions) {
	y2 := f.Metrics().Height.Round()

	opts.GeoM.Translate(0, float64(y2))

	text.DrawWithOptions(screen, s, f, opts)
}
