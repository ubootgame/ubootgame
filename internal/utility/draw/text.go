package draw

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"math"
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

func TextWithOptions(s string, f text.Face, opts *text.DrawOptions) *ebiten.Image {
	width, height := text.Measure(s, f, opts.LineSpacing)

	image := ebiten.NewImage(int(math.Round(width)), int(math.Round(height)))

	text.Draw(image, s, f, opts)

	return image
}
