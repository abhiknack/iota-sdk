package dtos

import (
	"time"

	"github.com/google/uuid"
)

type DriverCreateDTO struct {
	UserID        *int64    `form:"UserID"`
	FirstName     string    `form:"FirstName" validate:"required,max=100"`
	LastName      string    `form:"LastName" validate:"required,max=100"`
	LicenseNumber string    `form:"LicenseNumber" validate:"required,max=50"`
	LicenseExpiry time.Time `form:"LicenseExpiry" validate:"required"`
	Phone         string    `form:"Phone" validate:"omitempty,max=20"`
	Email         string    `form:"Email" validate:"omitempty,email,max=255"`
}

type DriverUpdateDTO struct {
	ID            uuid.UUID `form:"ID" validate:"required"`
	UserID        *int64    `form:"UserID"`
	FirstName     string    `form:"FirstName" validate:"required,max=100"`
	LastName      string    `form:"LastName" validate:"required,max=100"`
	LicenseNumber string    `form:"LicenseNumber" validate:"required,max=50"`
	LicenseExpiry time.Time `form:"LicenseExpiry" validate:"required"`
	Phone         string    `form:"Phone" validate:"omitempty,max=20"`
	Email         string    `form:"Email" validate:"omitempty,email,max=255"`
	Status        string    `form:"Status" validate:"required"`
}

type DriverFilterDTO struct {
	FirstName     string `form:"FirstName"`
	LastName      string `form:"LastName"`
	LicenseNumber string `form:"LicenseNumber"`
	Status        string `form:"Status"`
	Page          int    `form:"Page"`
	PageSize      int    `form:"PageSize"`
}
