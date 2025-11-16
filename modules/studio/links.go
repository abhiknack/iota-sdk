package studio

import (
	icons "github.com/iota-uz/icons/phosphor"
	"github.com/iota-uz/iota-sdk/pkg/types"
)

var (
	ModulesItem = types.NavigationItem{
		Name:        "NavigationLinks.StudioModules",
		Href:        "/studio/modules",
		Permissions: nil,
		Children:    nil,
	}
)

var StudioItem = types.NavigationItem{
	Name: "NavigationLinks.Studio",
	Href: "/studio",
	Icon: icons.Wrench(icons.Props{Size: "20"}),
	Children: []types.NavigationItem{
		ModulesItem,
	},
}

var NavItems = []types.NavigationItem{StudioItem}
