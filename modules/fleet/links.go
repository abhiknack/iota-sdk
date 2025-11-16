package fleet

import (
	icons "github.com/iota-uz/icons/phosphor"
	"github.com/iota-uz/iota-sdk/pkg/types"
)

var (
	DashboardItem = types.NavigationItem{
		Name:        "NavigationLinks.FleetDashboard",
		Href:        "/fleet/dashboard",
		Permissions: nil,
		Children:    nil,
	}
	VehiclesItem = types.NavigationItem{
		Name:        "NavigationLinks.Vehicles",
		Href:        "/fleet/vehicles",
		Permissions: nil,
		Children:    nil,
	}
	DriversItem = types.NavigationItem{
		Name:        "NavigationLinks.Drivers",
		Href:        "/fleet/drivers",
		Permissions: nil,
		Children:    nil,
	}
	TripsItem = types.NavigationItem{
		Name:        "NavigationLinks.Trips",
		Href:        "/fleet/trips",
		Permissions: nil,
		Children:    nil,
	}
	MaintenanceItem = types.NavigationItem{
		Name:        "NavigationLinks.Maintenance",
		Href:        "/fleet/maintenance",
		Permissions: nil,
		Children:    nil,
	}
	FuelItem = types.NavigationItem{
		Name:        "NavigationLinks.Fuel",
		Href:        "/fleet/fuel",
		Permissions: nil,
		Children:    nil,
	}
)

var FleetItem = types.NavigationItem{
	Name: "NavigationLinks.Fleet",
	Href: "/fleet",
	Icon: icons.Truck(icons.Props{Size: "20"}),
	Children: []types.NavigationItem{
		DashboardItem,
		VehiclesItem,
		DriversItem,
		TripsItem,
		MaintenanceItem,
		FuelItem,
	},
}

var NavItems = []types.NavigationItem{FleetItem}
