package maintenance

import (
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
)

type Maintenance interface {
	ID() uuid.UUID
	TenantID() uuid.UUID
	VehicleID() uuid.UUID
	ServiceType() enums.ServiceType
	ServiceDate() time.Time
	Odometer() int
	Cost() float64
	ServiceProvider() string
	Description() string
	NextServiceDue() *time.Time
	NextServiceOdometer() *int
	CreatedAt() time.Time
	UpdatedAt() time.Time

	UpdateCost(cost float64) Maintenance
	UpdateNextService(date *time.Time, odometer *int) Maintenance
}

type MaintenanceOption func(*maintenance)

func NewMaintenance(
	id uuid.UUID,
	tenantID uuid.UUID,
	vehicleID uuid.UUID,
	serviceType enums.ServiceType,
	serviceDate time.Time,
	odometer int,
	cost float64,
	serviceProvider string,
	description string,
	opts ...MaintenanceOption,
) Maintenance {
	m := &maintenance{
		id:              id,
		tenantID:        tenantID,
		vehicleID:       vehicleID,
		serviceType:     serviceType,
		serviceDate:     serviceDate,
		odometer:        odometer,
		cost:            cost,
		serviceProvider: serviceProvider,
		description:     description,
		createdAt:       time.Now(),
		updatedAt:       time.Now(),
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func WithNextServiceDue(date time.Time) MaintenanceOption {
	return func(m *maintenance) {
		m.nextServiceDue = &date
	}
}

func WithNextServiceOdometer(odometer int) MaintenanceOption {
	return func(m *maintenance) {
		m.nextServiceOdometer = &odometer
	}
}

func WithTimestamps(createdAt, updatedAt time.Time) MaintenanceOption {
	return func(m *maintenance) {
		m.createdAt = createdAt
		m.updatedAt = updatedAt
	}
}

type maintenance struct {
	id                  uuid.UUID
	tenantID            uuid.UUID
	vehicleID           uuid.UUID
	serviceType         enums.ServiceType
	serviceDate         time.Time
	odometer            int
	cost                float64
	serviceProvider     string
	description         string
	nextServiceDue      *time.Time
	nextServiceOdometer *int
	createdAt           time.Time
	updatedAt           time.Time
}

func (m *maintenance) ID() uuid.UUID {
	return m.id
}

func (m *maintenance) TenantID() uuid.UUID {
	return m.tenantID
}

func (m *maintenance) VehicleID() uuid.UUID {
	return m.vehicleID
}

func (m *maintenance) ServiceType() enums.ServiceType {
	return m.serviceType
}

func (m *maintenance) ServiceDate() time.Time {
	return m.serviceDate
}

func (m *maintenance) Odometer() int {
	return m.odometer
}

func (m *maintenance) Cost() float64 {
	return m.cost
}

func (m *maintenance) ServiceProvider() string {
	return m.serviceProvider
}

func (m *maintenance) Description() string {
	return m.description
}

func (m *maintenance) NextServiceDue() *time.Time {
	return m.nextServiceDue
}

func (m *maintenance) NextServiceOdometer() *int {
	return m.nextServiceOdometer
}

func (m *maintenance) CreatedAt() time.Time {
	return m.createdAt
}

func (m *maintenance) UpdatedAt() time.Time {
	return m.updatedAt
}

func (m *maintenance) UpdateCost(cost float64) Maintenance {
	return &maintenance{
		id:                  m.id,
		tenantID:            m.tenantID,
		vehicleID:           m.vehicleID,
		serviceType:         m.serviceType,
		serviceDate:         m.serviceDate,
		odometer:            m.odometer,
		cost:                cost,
		serviceProvider:     m.serviceProvider,
		description:         m.description,
		nextServiceDue:      m.nextServiceDue,
		nextServiceOdometer: m.nextServiceOdometer,
		createdAt:           m.createdAt,
		updatedAt:           time.Now(),
	}
}

func (m *maintenance) UpdateNextService(date *time.Time, odometer *int) Maintenance {
	return &maintenance{
		id:                  m.id,
		tenantID:            m.tenantID,
		vehicleID:           m.vehicleID,
		serviceType:         m.serviceType,
		serviceDate:         m.serviceDate,
		odometer:            m.odometer,
		cost:                m.cost,
		serviceProvider:     m.serviceProvider,
		description:         m.description,
		nextServiceDue:      date,
		nextServiceOdometer: odometer,
		createdAt:           m.createdAt,
		updatedAt:           time.Now(),
	}
}
