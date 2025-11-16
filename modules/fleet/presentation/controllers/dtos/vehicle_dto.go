package dtos

import (
	"time"

	"github.com/google/uuid"
)

type VehicleCreateDTO struct {
	Make               string    `form:"Make" validate:"required,max=100"`
	Model              string    `form:"Model" validate:"required,max=100"`
	Year               int       `form:"Year" validate:"required,min=1900,max=2100"`
	VIN                string    `form:"VIN" validate:"required,len=17"`
	LicensePlate       string    `form:"LicensePlate" validate:"required,max=20"`
	CurrentOdometer    int       `form:"CurrentOdometer" validate:"required,min=0"`
	RegistrationExpiry time.Time `form:"RegistrationExpiry" validate:"required"`
	InsuranceExpiry    time.Time `form:"InsuranceExpiry" validate:"required"`
}

type VehicleUpdateDTO struct {
	ID                 uuid.UUID `form:"ID" validate:"required"`
	Make               string    `form:"Make" validate:"required,max=100"`
	Model              string    `form:"Model" validate:"required,max=100"`
	Year               int       `form:"Year" validate:"required,min=1900,max=2100"`
	VIN                string    `form:"VIN" validate:"required,len=17"`
	LicensePlate       string    `form:"LicensePlate" validate:"required,max=20"`
	CurrentOdometer    int       `form:"CurrentOdometer" validate:"required,min=0"`
	RegistrationExpiry time.Time `form:"RegistrationExpiry" validate:"required"`
	InsuranceExpiry    time.Time `form:"InsuranceExpiry" validate:"required"`
}

type VehicleFilterDTO struct {
	Status       string `form:"Status"`
	Make         string `form:"Make"`
	Model        string `form:"Model"`
	Year         int    `form:"Year"`
	LicensePlate string `form:"LicensePlate"`
	Page         int    `form:"Page"`
	PageSize     int    `form:"PageSize"`
}
