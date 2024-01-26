package config

type Config struct {
	DefaultWidth  int
	DefaultHeight int
	Ratio         float64
	TargetTPS     int
}

var C *Config

func init() {
	C = &Config{
		DefaultWidth:  1200,
		DefaultHeight: 800,
		Ratio:         1200.0 / 800.0,
		TargetTPS:     60,
	}
}
