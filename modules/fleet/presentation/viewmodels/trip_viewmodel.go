package viewmodels

import "github.com/iota-uz/iota-sdk/components/base/badge"

type TripListViewModel struct {
	ID                 string
	VehicleID          string
	VehicleName        string
	DriverID           string
	DriverName         string
	Origin             string
	Destination        string
	StartTime          string
	EndTime            string
	Status             string
	StatusBadgeVariant badge.Variant
	Distance           int
	Duration           string
}

type TripDetailViewModel struct {
	ID                 string
	TenantID           string
	VehicleID          string
	VehicleName        string
	DriverID           string
	DriverName         string
	Origin             string
	Destination        string
	Purpose            string
	StartTime          string
	EndTime            string
	StartOdometer      int
	EndOdometer        int
	Status             string
	StatusBadgeVariant badge.Variant
	Distance           int
	Duration           string
	AverageSpeed       float64
	CreatedAt          string
	UpdatedAt          string
	CanComplete        bool
	CanCancel          bool
}

type TripFormViewModel struct {
	ID            string
	VehicleID     string
	DriverID      string
	Origin        string
	Destination   string
	Purpose       string
	StartTime     string
	StartOdometer int
	IsEdit        bool
	Errors        map[string]string
	Vehicles      []VehicleOption
	Drivers       []DriverOption
	HasConflict   bool
	ConflictMsg   string
}

type VehicleOption struct {
	ID    string
	Label string
}

type DriverOption struct {
	ID    string
	Label string
}
