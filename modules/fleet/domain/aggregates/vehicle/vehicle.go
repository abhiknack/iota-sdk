package vehicle

import (
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
)

type Vehicle interface {
	ID() uuid.UUID
	TenantID() uuid.UUID
	Make() string
	Model() string
	Year() int
	VIN() string
	LicensePlate() string
	Status() enums.VehicleStatus
	CurrentOdometer() int
	RegistrationExpiry() time.Time
	InsuranceExpiry() time.Time
	CreatedAt() time.Time
	UpdatedAt() time.Time

	UpdateStatus(status enums.VehicleStatus) Vehicle
	UpdateOdometer(reading int) Vehicle
	UpdateDetails(make, model string, year int) Vehicle
}

type VehicleOption func(*vehicle)

func NewVehicle(
	id uuid.UUID,
	tenantID uuid.UUID,
	make string,
	model string,
	year int,
	vin string,
	licensePlate string,
	opts ...VehicleOption,
) Vehicle {
	v := &vehicle{
		id:                 id,
		tenantID:           tenantID,
		make:               make,
		model:              model,
		year:               year,
		vin:                vin,
		licensePlate:       licensePlate,
		status:             enums.VehicleStatusAvailable,
		currentOdometer:    0,
		registrationExpiry: time.Now(),
		insuranceExpiry:    time.Now(),
		createdAt:          time.Now(),
		updatedAt:          time.Now(),
	}

	for _, opt := range opts {
		opt(v)
	}

	return v
}

func WithStatus(status enums.VehicleStatus) VehicleOption {
	return func(v *vehicle) {
		v.status = status
	}
}

func WithOdometer(odometer int) VehicleOption {
	return func(v *vehicle) {
		v.currentOdometer = odometer
	}
}

func WithRegistrationExpiry(expiry time.Time) VehicleOption {
	return func(v *vehicle) {
		v.registrationExpiry = expiry
	}
}

func WithInsuranceExpiry(expiry time.Time) VehicleOption {
	return func(v *vehicle) {
		v.insuranceExpiry = expiry
	}
}

func WithTimestamps(createdAt, updatedAt time.Time) VehicleOption {
	return func(v *vehicle) {
		v.createdAt = createdAt
		v.updatedAt = updatedAt
	}
}

type vehicle struct {
	id                 uuid.UUID
	tenantID           uuid.UUID
	make               string
	model              string
	year               int
	vin                string
	licensePlate       string
	status             enums.VehicleStatus
	currentOdometer    int
	registrationExpiry time.Time
	insuranceExpiry    time.Time
	createdAt          time.Time
	updatedAt          time.Time
}

func (v *vehicle) ID() uuid.UUID {
	return v.id
}

func (v *vehicle) TenantID() uuid.UUID {
	return v.tenantID
}

func (v *vehicle) Make() string {
	return v.make
}

func (v *vehicle) Model() string {
	return v.model
}

func (v *vehicle) Year() int {
	return v.year
}

func (v *vehicle) VIN() string {
	return v.vin
}

func (v *vehicle) LicensePlate() string {
	return v.licensePlate
}

func (v *vehicle) Status() enums.VehicleStatus {
	return v.status
}

func (v *vehicle) CurrentOdometer() int {
	return v.currentOdometer
}

func (v *vehicle) RegistrationExpiry() time.Time {
	return v.registrationExpiry
}

func (v *vehicle) InsuranceExpiry() time.Time {
	return v.insuranceExpiry
}

func (v *vehicle) CreatedAt() time.Time {
	return v.createdAt
}

func (v *vehicle) UpdatedAt() time.Time {
	return v.updatedAt
}

func (v *vehicle) UpdateStatus(status enums.VehicleStatus) Vehicle {
	return &vehicle{
		id:                 v.id,
		tenantID:           v.tenantID,
		make:               v.make,
		model:              v.model,
		year:               v.year,
		vin:                v.vin,
		licensePlate:       v.licensePlate,
		status:             status,
		currentOdometer:    v.currentOdometer,
		registrationExpiry: v.registrationExpiry,
		insuranceExpiry:    v.insuranceExpiry,
		createdAt:          v.createdAt,
		updatedAt:          time.Now(),
	}
}

func (v *vehicle) UpdateOdometer(reading int) Vehicle {
	return &vehicle{
		id:                 v.id,
		tenantID:           v.tenantID,
		make:               v.make,
		model:              v.model,
		year:               v.year,
		vin:                v.vin,
		licensePlate:       v.licensePlate,
		status:             v.status,
		currentOdometer:    reading,
		registrationExpiry: v.registrationExpiry,
		insuranceExpiry:    v.insuranceExpiry,
		createdAt:          v.createdAt,
		updatedAt:          time.Now(),
	}
}

func (v *vehicle) UpdateDetails(make, model string, year int) Vehicle {
	return &vehicle{
		id:                 v.id,
		tenantID:           v.tenantID,
		make:               make,
		model:              model,
		year:               year,
		vin:                v.vin,
		licensePlate:       v.licensePlate,
		status:             v.status,
		currentOdometer:    v.currentOdometer,
		registrationExpiry: v.registrationExpiry,
		insuranceExpiry:    v.insuranceExpiry,
		createdAt:          v.createdAt,
		updatedAt:          time.Now(),
	}
}
