package settings

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/samber/do"
	"golang.org/x/image/font/gofont/gobold"
	"log"
)

type Provider[S interface{}] interface {
	Settings() *Settings[S]
	UpdateDebugFontFace(scalingFactor float64)
}

var defaultDebugFontSize = 12.0

type settingsProvider[S interface{}] struct {
	settings *Settings[S]
}

func NewProvider[S interface{}](_ *do.Injector, settings Settings[S]) (Provider[S], error) {
	return &settingsProvider[S]{settings: &settings}, nil
}

func (provider *settingsProvider[S]) Set(settings Settings[S]) {
	provider.settings = &settings
}

func (provider *settingsProvider[S]) Settings() *Settings[S] {
	return provider.settings
}

func (provider *settingsProvider[S]) UpdateDebugFontFace(scalingFactor float64) {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(gobold.TTF))
	if err != nil {
		log.Fatal(err)
	}

	provider.settings.Debug.FontScale = scalingFactor

	provider.settings.Debug.FontFace = &text.GoTextFace{
		Source: source,
		Size:   defaultDebugFontSize * scalingFactor,
	}
}
