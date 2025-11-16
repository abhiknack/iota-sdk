package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Vehicle struct {
	ID                 uuid.UUID
	TenantID           uuid.UUID
	Make               string
	Model              string
	Year               int
	VIN                string
	LicensePlate       string
	Status             int
	CurrentOdometer    int
	RegistrationExpiry time.Time
	InsuranceExpiry    time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          sql.NullTime
}

type Driver struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	UserID        sql.NullInt64
	FirstName     string
	LastName      string
	LicenseNumber string
	LicenseExpiry time.Time
	Phone         sql.NullString
	Email         sql.NullString
	Status        int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     sql.NullTime
}

type Trip struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	VehicleID     uuid.UUID
	DriverID      uuid.UUID
	Origin        string
	Destination   string
	Purpose       sql.NullString
	StartTime     time.Time
	EndTime       sql.NullTime
	StartOdometer int
	EndOdometer   sql.NullInt32
	Status        int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     sql.NullTime
}

type Maintenance struct {
	ID                  uuid.UUID
	TenantID            uuid.UUID
	VehicleID           uuid.UUID
	ServiceType         int
	ServiceDate         time.Time
	Odometer            int
	Cost                float64
	ServiceProvider     sql.NullString
	Description         sql.NullString
	NextServiceDue      sql.NullTime
	NextServiceOdometer sql.NullInt32
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           sql.NullTime
}

type FuelEntry struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	VehicleID uuid.UUID
	DriverID  uuid.NullUUID
	Date      time.Time
	Quantity  float64
	Cost      float64
	Odometer  int
	FuelType  int
	Location  sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}
