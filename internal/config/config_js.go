//go:build js || wasm

package config

import (
	"gonum.org/v1/gonum/spatial/r2"
	"strconv"
	"syscall/js"
)

type Config struct {
	DefaultWindowSize r2.Vec
	Ratio             float64
	TargetTPS         int
	Debug             bool
}

var (
	C *Config
)

func init() {
	defaultConfiguration := getDefaultConfiguration()
	mergedConfiguration := getLocalStorageConfiguration(defaultConfiguration)

	C = mergedConfiguration
}

func getDefaultConfiguration() *Config {
	return &Config{
		DefaultWindowSize: r2.Vec{X: 1280, Y: 720},
		Ratio:             1280.0 / 720.0,
		TargetTPS:         60,
		Debug:             true,
	}
}

func getLocalStorageConfiguration(config *Config) *Config {
	localStorage := js.Global().Get("localStorage")

	DefaultWindowSizeX := localStorage.Call("getItem", "DefaultWindowSize.X")
	DefaultWindowSizeY := localStorage.Call("getItem", "DefaultWindowSize.Y")
	TargetTPS := localStorage.Call("getItem", "TargetTPS")
	Debug := localStorage.Call("getItem", "Debug")

	if js.Value.IsUndefined(DefaultWindowSizeX) != true && js.Value.IsUndefined(DefaultWindowSizeY) != true && js.Value.IsNull(DefaultWindowSizeX) != true && js.Value.IsNull(DefaultWindowSizeY) != true {
		x, _ := strconv.ParseFloat(DefaultWindowSizeX.String(), 32)
		y, _ := strconv.ParseFloat(DefaultWindowSizeY.String(), 32)
		println(y)
		println(x)
		config.DefaultWindowSize = r2.Vec{X: x, Y: y}
		config.Ratio = x / y
	}

	if js.Value.IsUndefined(TargetTPS) != true && js.Value.IsNull(TargetTPS) != true {
		converted, _ := strconv.ParseInt(TargetTPS.String(), 10, 32)
		config.TargetTPS = int(converted)
	}

	if js.Value.IsUndefined(Debug) != true && js.Value.IsNull(Debug) != true {
		config.Debug, _ = strconv.ParseBool(Debug.String())
	}

	return config
}
