package fuel_entry

import (
	"context"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
)

type Field string

const (
	FieldID        Field = "id"
	FieldTenantID  Field = "tenant_id"
	FieldVehicleID Field = "vehicle_id"
	FieldDriverID  Field = "driver_id"
	FieldDate      Field = "date"
	FieldQuantity  Field = "quantity"
	FieldCost      Field = "cost"
	FieldOdometer  Field = "odometer"
	FieldFuelType  Field = "fuel_type"
	FieldLocation  Field = "location"
	FieldCreatedAt Field = "created_at"
	FieldUpdatedAt Field = "updated_at"
)

type FindParams struct {
	TenantID       uuid.UUID
	VehicleID      *uuid.UUID
	DriverID       *uuid.UUID
	FuelType       *enums.FuelType
	DateFrom       *string
	DateTo         *string
	Limit          int
	Offset         int
	SortBy         Field
	SortDesc       bool
	IncludeDeleted bool
}

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (FuelEntry, error)
	GetPaginated(ctx context.Context, params *FindParams) ([]FuelEntry, error)
	Count(ctx context.Context, params *FindParams) (int64, error)
	GetByVehicle(ctx context.Context, vehicleID uuid.UUID) ([]FuelEntry, error)
	GetByDriver(ctx context.Context, driverID uuid.UUID) ([]FuelEntry, error)
	GetLastEntry(ctx context.Context, vehicleID uuid.UUID) (FuelEntry, error)
	Create(ctx context.Context, entry FuelEntry) (FuelEntry, error)
	Update(ctx context.Context, entry FuelEntry) (FuelEntry, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
