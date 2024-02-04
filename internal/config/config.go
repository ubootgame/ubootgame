//go:build darwin || linxux || windows

package config

import (
	"github.com/ovh/configstore"
	"gonum.org/v1/gonum/spatial/r2"
)

type Config struct {
	DefaultWindowSize r2.Vec
	Ratio             float64
	TargetTPS         int
	Debug             bool
}

var C *Config

func init() {
	configstore.RegisterProvider("defaults", getDefaultConfiguration)
	configstore.File("configuration.yml")

	items, err := configstore.Filter().Squash().GetItemList()
	if err != nil {
		panic(err)
	}

	x, _ := items.GetItemValueFloat("DefaultWindowSize.X")
	y, _ := items.GetItemValueFloat("DefaultWindowSize.Y")
	targetTPS, _ := items.GetItemValueInt("TargetTPS")
	debug, _ := items.GetItemValueBool("Debug")

	C = createConfiguration(x, y, int(targetTPS), debug)
}

func getDefaultConfiguration() (configstore.ItemList, error) {
	ret := configstore.ItemList{
		Items: []configstore.Item{
			configstore.NewItem("DefaultWindowSize.X", "1280", 0),
			configstore.NewItem("DefaultWindowSize.Y", "720", 0),
			configstore.NewItem("TargetTPS", "60", 0),
			configstore.NewItem("Debug", "true", 0),
		},
	}
	return ret, nil
}

func createConfiguration(x float64, y float64, targetTps int, debug bool) *Config {
	return &Config{
		DefaultWindowSize: r2.Vec{X: x, Y: y},
		Ratio:             x / y,
		TargetTPS:         targetTps,
		Debug:             debug,
	}
}
