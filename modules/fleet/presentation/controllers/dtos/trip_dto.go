package dtos

import (
	"time"

	"github.com/google/uuid"
)

type TripCreateDTO struct {
	VehicleID     uuid.UUID `form:"VehicleID" validate:"required"`
	DriverID      uuid.UUID `form:"DriverID" validate:"required"`
	Origin        string    `form:"Origin" validate:"required,max=255"`
	Destination   string    `form:"Destination" validate:"required,max=255"`
	Purpose       string    `form:"Purpose" validate:"omitempty"`
	StartTime     time.Time `form:"StartTime" validate:"required"`
	StartOdometer int       `form:"StartOdometer" validate:"required,min=0"`
}

type TripCompleteDTO struct {
	ID          uuid.UUID `form:"ID" validate:"required"`
	EndTime     time.Time `form:"EndTime" validate:"required"`
	EndOdometer int       `form:"EndOdometer" validate:"required,min=0"`
}

type TripFilterDTO struct {
	VehicleID uuid.UUID `form:"VehicleID"`
	DriverID  uuid.UUID `form:"DriverID"`
	Status    string    `form:"Status"`
	StartDate time.Time `form:"StartDate"`
	EndDate   time.Time `form:"EndDate"`
	Page      int       `form:"Page"`
	PageSize  int       `form:"PageSize"`
}
