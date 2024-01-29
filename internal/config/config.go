package config

import "gonum.org/v1/gonum/spatial/r2"

type Config struct {
	DefaultOuterSize, ActualOuterSize, VirtualResolution r2.Vec
	Ratio                                                float64
	TargetTPS                                            int
	Debug                                                bool
}

var C *Config

func init() {
	C = &Config{
		DefaultOuterSize:  r2.Vec{X: 1280, Y: 720},
		ActualOuterSize:   r2.Vec{X: 1280, Y: 720},
		VirtualResolution: r2.Vec{X: 1280, Y: 720},
		Ratio:             1280.0 / 720.0,
		TargetTPS:         60,
		Debug:             true,
	}
}
