package trip

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
)

type Field string

const (
	FieldID            Field = "id"
	FieldTenantID      Field = "tenant_id"
	FieldVehicleID     Field = "vehicle_id"
	FieldDriverID      Field = "driver_id"
	FieldOrigin        Field = "origin"
	FieldDestination   Field = "destination"
	FieldPurpose       Field = "purpose"
	FieldStartTime     Field = "start_time"
	FieldEndTime       Field = "end_time"
	FieldStartOdometer Field = "start_odometer"
	FieldEndOdometer   Field = "end_odometer"
	FieldStatus        Field = "status"
	FieldCreatedAt     Field = "created_at"
	FieldUpdatedAt     Field = "updated_at"
)

type FindParams struct {
	TenantID       uuid.UUID
	VehicleID      *uuid.UUID
	DriverID       *uuid.UUID
	Status         *enums.TripStatus
	StartTimeFrom  *time.Time
	StartTimeTo    *time.Time
	Search         *string
	Limit          int
	Offset         int
	SortBy         Field
	SortDesc       bool
	IncludeDeleted bool
}

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (Trip, error)
	GetPaginated(ctx context.Context, params *FindParams) ([]Trip, error)
	Count(ctx context.Context, params *FindParams) (int64, error)
	GetByVehicle(ctx context.Context, vehicleID uuid.UUID) ([]Trip, error)
	GetByDriver(ctx context.Context, driverID uuid.UUID) ([]Trip, error)
	GetActiveTrips(ctx context.Context, tenantID uuid.UUID) ([]Trip, error)
	CheckConflict(ctx context.Context, vehicleID uuid.UUID, startTime, endTime time.Time, excludeTripID *uuid.UUID) (bool, error)
	Create(ctx context.Context, trip Trip) (Trip, error)
	Update(ctx context.Context, trip Trip) (Trip, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
