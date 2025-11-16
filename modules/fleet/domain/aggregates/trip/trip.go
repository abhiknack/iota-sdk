package trip

import (
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
)

type Trip interface {
	ID() uuid.UUID
	TenantID() uuid.UUID
	VehicleID() uuid.UUID
	DriverID() uuid.UUID
	Origin() string
	Destination() string
	Purpose() string
	StartTime() time.Time
	EndTime() *time.Time
	StartOdometer() int
	EndOdometer() *int
	Status() enums.TripStatus
	CreatedAt() time.Time
	UpdatedAt() time.Time

	Complete(endTime time.Time, endOdometer int) Trip
	Cancel(reason string) Trip
	UpdateRoute(origin, destination string) Trip
}

type TripOption func(*trip)

func NewTrip(
	id uuid.UUID,
	tenantID uuid.UUID,
	vehicleID uuid.UUID,
	driverID uuid.UUID,
	origin string,
	destination string,
	purpose string,
	startTime time.Time,
	startOdometer int,
	opts ...TripOption,
) Trip {
	t := &trip{
		id:            id,
		tenantID:      tenantID,
		vehicleID:     vehicleID,
		driverID:      driverID,
		origin:        origin,
		destination:   destination,
		purpose:       purpose,
		startTime:     startTime,
		startOdometer: startOdometer,
		status:        enums.TripStatusScheduled,
		createdAt:     time.Now(),
		updatedAt:     time.Now(),
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

func WithStatus(status enums.TripStatus) TripOption {
	return func(t *trip) {
		t.status = status
	}
}

func WithEndTime(endTime time.Time) TripOption {
	return func(t *trip) {
		t.endTime = &endTime
	}
}

func WithEndOdometer(endOdometer int) TripOption {
	return func(t *trip) {
		t.endOdometer = &endOdometer
	}
}

func WithTimestamps(createdAt, updatedAt time.Time) TripOption {
	return func(t *trip) {
		t.createdAt = createdAt
		t.updatedAt = updatedAt
	}
}

type trip struct {
	id            uuid.UUID
	tenantID      uuid.UUID
	vehicleID     uuid.UUID
	driverID      uuid.UUID
	origin        string
	destination   string
	purpose       string
	startTime     time.Time
	endTime       *time.Time
	startOdometer int
	endOdometer   *int
	status        enums.TripStatus
	createdAt     time.Time
	updatedAt     time.Time
}

func (t *trip) ID() uuid.UUID {
	return t.id
}

func (t *trip) TenantID() uuid.UUID {
	return t.tenantID
}

func (t *trip) VehicleID() uuid.UUID {
	return t.vehicleID
}

func (t *trip) DriverID() uuid.UUID {
	return t.driverID
}

func (t *trip) Origin() string {
	return t.origin
}

func (t *trip) Destination() string {
	return t.destination
}

func (t *trip) Purpose() string {
	return t.purpose
}

func (t *trip) StartTime() time.Time {
	return t.startTime
}

func (t *trip) EndTime() *time.Time {
	return t.endTime
}

func (t *trip) StartOdometer() int {
	return t.startOdometer
}

func (t *trip) EndOdometer() *int {
	return t.endOdometer
}

func (t *trip) Status() enums.TripStatus {
	return t.status
}

func (t *trip) CreatedAt() time.Time {
	return t.createdAt
}

func (t *trip) UpdatedAt() time.Time {
	return t.updatedAt
}

func (t *trip) Complete(endTime time.Time, endOdometer int) Trip {
	return &trip{
		id:            t.id,
		tenantID:      t.tenantID,
		vehicleID:     t.vehicleID,
		driverID:      t.driverID,
		origin:        t.origin,
		destination:   t.destination,
		purpose:       t.purpose,
		startTime:     t.startTime,
		endTime:       &endTime,
		startOdometer: t.startOdometer,
		endOdometer:   &endOdometer,
		status:        enums.TripStatusCompleted,
		createdAt:     t.createdAt,
		updatedAt:     time.Now(),
	}
}

func (t *trip) Cancel(reason string) Trip {
	return &trip{
		id:            t.id,
		tenantID:      t.tenantID,
		vehicleID:     t.vehicleID,
		driverID:      t.driverID,
		origin:        t.origin,
		destination:   t.destination,
		purpose:       reason,
		startTime:     t.startTime,
		endTime:       t.endTime,
		startOdometer: t.startOdometer,
		endOdometer:   t.endOdometer,
		status:        enums.TripStatusCancelled,
		createdAt:     t.createdAt,
		updatedAt:     time.Now(),
	}
}

func (t *trip) UpdateRoute(origin, destination string) Trip {
	return &trip{
		id:            t.id,
		tenantID:      t.tenantID,
		vehicleID:     t.vehicleID,
		driverID:      t.driverID,
		origin:        origin,
		destination:   destination,
		purpose:       t.purpose,
		startTime:     t.startTime,
		endTime:       t.endTime,
		startOdometer: t.startOdometer,
		endOdometer:   t.endOdometer,
		status:        t.status,
		createdAt:     t.createdAt,
		updatedAt:     time.Now(),
	}
}
