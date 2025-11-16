package viewmodels

import "github.com/iota-uz/iota-sdk/components/base/badge"

type DriverListViewModel struct {
	ID                   string
	FirstName            string
	LastName             string
	FullName             string
	LicenseNumber        string
	LicenseExpiry        string
	LicenseExpiryWarning bool
	Phone                string
	Email                string
	Status               string
	StatusBadgeVariant   badge.Variant
}

type DriverDetailViewModel struct {
	ID                   string
	TenantID             string
	UserID               string
	FirstName            string
	LastName             string
	FullName             string
	LicenseNumber        string
	LicenseExpiry        string
	LicenseExpiryWarning bool
	Phone                string
	Email                string
	Status               string
	StatusBadgeVariant   badge.Variant
	CreatedAt            string
	UpdatedAt            string
	TotalTrips           int
	ActiveTrips          int
}

type DriverFormViewModel struct {
	ID            string
	UserID        string
	FirstName     string
	LastName      string
	LicenseNumber string
	LicenseExpiry string
	Phone         string
	Email         string
	Status        string
	IsEdit        bool
	Errors        map[string]string
}
