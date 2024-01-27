package config

import "gonum.org/v1/gonum/spatial/r2"

type Config struct {
	DefaultOuterSize, ActualOuterSize, VirtualResolution r2.Vec
	TargetTPS                                            int
	Debug                                                bool
}

var C *Config

func init() {
	C = &Config{
		DefaultOuterSize:  r2.Vec{X: 1280, Y: 720},
		ActualOuterSize:   r2.Vec{X: 1280, Y: 720},
		VirtualResolution: r2.Vec{X: 1920, Y: 1080},
		TargetTPS:         60,
		Debug:             true,
	}
}
