package dtos

import (
	"time"

	"github.com/google/uuid"
)

type FuelEntryCreateDTO struct {
	VehicleID uuid.UUID  `form:"VehicleID" validate:"required"`
	DriverID  *uuid.UUID `form:"DriverID"`
	Date      time.Time  `form:"Date" validate:"required"`
	Quantity  float64    `form:"Quantity" validate:"required,min=0"`
	Cost      float64    `form:"Cost" validate:"required,min=0"`
	Odometer  int        `form:"Odometer" validate:"required,min=0"`
	FuelType  string     `form:"FuelType" validate:"required"`
	Location  string     `form:"Location" validate:"omitempty,max=255"`
}

type FuelEntryUpdateDTO struct {
	ID        uuid.UUID  `form:"ID" validate:"required"`
	VehicleID uuid.UUID  `form:"VehicleID" validate:"required"`
	DriverID  *uuid.UUID `form:"DriverID"`
	Date      time.Time  `form:"Date" validate:"required"`
	Quantity  float64    `form:"Quantity" validate:"required,min=0"`
	Cost      float64    `form:"Cost" validate:"required,min=0"`
	Odometer  int        `form:"Odometer" validate:"required,min=0"`
	FuelType  string     `form:"FuelType" validate:"required"`
	Location  string     `form:"Location" validate:"omitempty,max=255"`
}

type FuelEntryFilterDTO struct {
	VehicleID uuid.UUID `form:"VehicleID"`
	DriverID  uuid.UUID `form:"DriverID"`
	FuelType  string    `form:"FuelType"`
	StartDate time.Time `form:"StartDate"`
	EndDate   time.Time `form:"EndDate"`
	Page      int       `form:"Page"`
	PageSize  int       `form:"PageSize"`
}
