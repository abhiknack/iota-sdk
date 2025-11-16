package fuel_entry

import (
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
)

type FuelEntry interface {
	ID() uuid.UUID
	TenantID() uuid.UUID
	VehicleID() uuid.UUID
	DriverID() *uuid.UUID
	Date() time.Time
	Quantity() float64
	Cost() float64
	Odometer() int
	FuelType() enums.FuelType
	Location() string
	CreatedAt() time.Time
	UpdatedAt() time.Time

	CalculateEfficiency(previousOdometer int) float64
	UpdateCost(cost float64) FuelEntry
}

type FuelEntryOption func(*fuelEntry)

func NewFuelEntry(
	id uuid.UUID,
	tenantID uuid.UUID,
	vehicleID uuid.UUID,
	date time.Time,
	quantity float64,
	cost float64,
	odometer int,
	fuelType enums.FuelType,
	location string,
	opts ...FuelEntryOption,
) FuelEntry {
	f := &fuelEntry{
		id:        id,
		tenantID:  tenantID,
		vehicleID: vehicleID,
		date:      date,
		quantity:  quantity,
		cost:      cost,
		odometer:  odometer,
		fuelType:  fuelType,
		location:  location,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

func WithDriverID(driverID uuid.UUID) FuelEntryOption {
	return func(f *fuelEntry) {
		f.driverID = &driverID
	}
}

func WithTimestamps(createdAt, updatedAt time.Time) FuelEntryOption {
	return func(f *fuelEntry) {
		f.createdAt = createdAt
		f.updatedAt = updatedAt
	}
}

type fuelEntry struct {
	id        uuid.UUID
	tenantID  uuid.UUID
	vehicleID uuid.UUID
	driverID  *uuid.UUID
	date      time.Time
	quantity  float64
	cost      float64
	odometer  int
	fuelType  enums.FuelType
	location  string
	createdAt time.Time
	updatedAt time.Time
}

func (f *fuelEntry) ID() uuid.UUID {
	return f.id
}

func (f *fuelEntry) TenantID() uuid.UUID {
	return f.tenantID
}

func (f *fuelEntry) VehicleID() uuid.UUID {
	return f.vehicleID
}

func (f *fuelEntry) DriverID() *uuid.UUID {
	return f.driverID
}

func (f *fuelEntry) Date() time.Time {
	return f.date
}

func (f *fuelEntry) Quantity() float64 {
	return f.quantity
}

func (f *fuelEntry) Cost() float64 {
	return f.cost
}

func (f *fuelEntry) Odometer() int {
	return f.odometer
}

func (f *fuelEntry) FuelType() enums.FuelType {
	return f.fuelType
}

func (f *fuelEntry) Location() string {
	return f.location
}

func (f *fuelEntry) CreatedAt() time.Time {
	return f.createdAt
}

func (f *fuelEntry) UpdatedAt() time.Time {
	return f.updatedAt
}

func (f *fuelEntry) CalculateEfficiency(previousOdometer int) float64 {
	if previousOdometer >= f.odometer || f.quantity == 0 {
		return 0
	}
	distance := float64(f.odometer - previousOdometer)
	return distance / f.quantity
}

func (f *fuelEntry) UpdateCost(cost float64) FuelEntry {
	return &fuelEntry{
		id:        f.id,
		tenantID:  f.tenantID,
		vehicleID: f.vehicleID,
		driverID:  f.driverID,
		date:      f.date,
		quantity:  f.quantity,
		cost:      cost,
		odometer:  f.odometer,
		fuelType:  f.fuelType,
		location:  f.location,
		createdAt: f.createdAt,
		updatedAt: time.Now(),
	}
}
