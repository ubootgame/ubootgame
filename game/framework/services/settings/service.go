package settings

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/gobold"
	"log"
)

var defaultDebugFontSize = 12.0

type Service[S any] struct {
	settings *Settings[S]
}

func NewService[S any](settings *Settings[S]) *Service[S] {
	return &Service[S]{settings: settings}
}

func (service *Service[S]) Settings() *Settings[S] {
	return service.settings
}

func (service *Service[S]) UpdateDebugFontFace(scalingFactor float64) {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(gobold.TTF))
	if err != nil {
		log.Fatal(err)
	}

	service.settings.Debug.FontScale = scalingFactor

	service.settings.Debug.FontFace = &text.GoTextFace{
		Source: source,
		Size:   defaultDebugFontSize * scalingFactor,
	}
}
