package viewmodels

import (
	"github.com/iota-uz/iota-sdk/components/charts"
)

type DashboardViewModel struct {
	TotalVehicles        int
	ActiveVehicles       int
	VehiclesInUse        int
	VehiclesMaintenance  int
	ActiveDrivers        int
	UpcomingMaintenance  int
	FuelCostMonth        string
	MaintenanceCostMonth string
	TotalCostMonth       string
	UtilizationChart     UtilizationChartViewModel
	CostTrendChart       CostTrendViewModel
}

type UtilizationChartViewModel struct {
	Options charts.ChartOptions
}

type CostTrendViewModel struct {
	Options charts.ChartOptions
}
