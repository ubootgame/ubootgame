package components

import "github.com/yohamta/donburi"

type MenuData struct {
	IsOpen bool
}

var Menu = donburi.NewComponentType[MenuData](MenuData{IsOpen: false})

func (menu *MenuData) ToggleMenu() bool {
	menu.IsOpen = !menu.IsOpen
	return menu.IsOpen
}
