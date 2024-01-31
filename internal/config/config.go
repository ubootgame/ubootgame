package config

import "gonum.org/v1/gonum/spatial/r2"

type Config struct {
	DefaultWindowSize r2.Vec
	Ratio             float64
	TargetTPS         int
	Debug             bool
}

var C *Config

func init() {
	C = &Config{
		DefaultWindowSize: r2.Vec{X: 1280, Y: 720},
		Ratio:             1280.0 / 720.0,
		TargetTPS:         60,
		Debug:             true,
	}
}
