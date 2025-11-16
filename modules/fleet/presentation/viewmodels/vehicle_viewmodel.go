package viewmodels

import "github.com/iota-uz/iota-sdk/components/base/badge"

type VehicleListViewModel struct {
	ID                 string
	Make               string
	Model              string
	Year               int
	LicensePlate       string
	Status             string
	StatusBadgeVariant badge.Variant
	CurrentOdometer    int
	RegistrationExpiry string
	InsuranceExpiry    string
}

type VehicleDetailViewModel struct {
	ID                 string
	TenantID           string
	Make               string
	Model              string
	Year               int
	VIN                string
	LicensePlate       string
	Status             string
	StatusBadgeVariant badge.Variant
	CurrentOdometer    int
	RegistrationExpiry string
	InsuranceExpiry    string
	CreatedAt          string
	UpdatedAt          string
	FullName           string
}

type VehicleFormViewModel struct {
	ID                 string
	Make               string
	Model              string
	Year               int
	VIN                string
	LicensePlate       string
	CurrentOdometer    int
	RegistrationExpiry string
	InsuranceExpiry    string
	IsEdit             bool
	Errors             map[string]string
}
