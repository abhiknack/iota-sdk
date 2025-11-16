package dtos

import (
	"time"

	"github.com/google/uuid"
)

type MaintenanceCreateDTO struct {
	VehicleID           uuid.UUID  `form:"VehicleID" validate:"required"`
	ServiceType         string     `form:"ServiceType" validate:"required"`
	ServiceDate         time.Time  `form:"ServiceDate" validate:"required"`
	Odometer            int        `form:"Odometer" validate:"required,min=0"`
	Cost                float64    `form:"Cost" validate:"required,min=0"`
	ServiceProvider     string     `form:"ServiceProvider" validate:"omitempty,max=255"`
	Description         string     `form:"Description" validate:"omitempty"`
	NextServiceDue      *time.Time `form:"NextServiceDue"`
	NextServiceOdometer *int       `form:"NextServiceOdometer" validate:"omitempty,min=0"`
}

type MaintenanceUpdateDTO struct {
	ID                  uuid.UUID  `form:"ID" validate:"required"`
	VehicleID           uuid.UUID  `form:"VehicleID" validate:"required"`
	ServiceType         string     `form:"ServiceType" validate:"required"`
	ServiceDate         time.Time  `form:"ServiceDate" validate:"required"`
	Odometer            int        `form:"Odometer" validate:"required,min=0"`
	Cost                float64    `form:"Cost" validate:"required,min=0"`
	ServiceProvider     string     `form:"ServiceProvider" validate:"omitempty,max=255"`
	Description         string     `form:"Description" validate:"omitempty"`
	NextServiceDue      *time.Time `form:"NextServiceDue"`
	NextServiceOdometer *int       `form:"NextServiceOdometer" validate:"omitempty,min=0"`
}

type MaintenanceFilterDTO struct {
	VehicleID   uuid.UUID `form:"VehicleID"`
	ServiceType string    `form:"ServiceType"`
	StartDate   time.Time `form:"StartDate"`
	EndDate     time.Time `form:"EndDate"`
	Page        int       `form:"Page"`
	PageSize    int       `form:"PageSize"`
}
