package config

type Config struct {
	Width  int
	Height int
}

var C *Config

func init() {
	C = &Config{
		Width:  1200,
		Height: 800,
	}
}
