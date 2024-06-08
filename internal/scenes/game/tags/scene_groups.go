package tags

import "github.com/yohamta/donburi"

var (
	EnvironmentTag = donburi.NewTag().SetName("Environment")
	ObjectsTag     = donburi.NewTag().SetName("Objects")
	ProjectilesTag = donburi.NewTag().SetName("Projectiles")
)
