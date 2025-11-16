package maintenance

import (
	"context"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
)

type Field string

const (
	FieldID                  Field = "id"
	FieldTenantID            Field = "tenant_id"
	FieldVehicleID           Field = "vehicle_id"
	FieldServiceType         Field = "service_type"
	FieldServiceDate         Field = "service_date"
	FieldOdometer            Field = "odometer"
	FieldCost                Field = "cost"
	FieldServiceProvider     Field = "service_provider"
	FieldDescription         Field = "description"
	FieldNextServiceDue      Field = "next_service_due"
	FieldNextServiceOdometer Field = "next_service_odometer"
	FieldCreatedAt           Field = "created_at"
	FieldUpdatedAt           Field = "updated_at"
)

type FindParams struct {
	TenantID        uuid.UUID
	VehicleID       *uuid.UUID
	ServiceType     *enums.ServiceType
	ServiceDateFrom *string
	ServiceDateTo   *string
	Search          *string
	Limit           int
	Offset          int
	SortBy          Field
	SortDesc        bool
	IncludeDeleted  bool
}

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (Maintenance, error)
	GetPaginated(ctx context.Context, params *FindParams) ([]Maintenance, error)
	Count(ctx context.Context, params *FindParams) (int64, error)
	GetByVehicle(ctx context.Context, vehicleID uuid.UUID) ([]Maintenance, error)
	GetDueMaintenance(ctx context.Context, tenantID uuid.UUID) ([]Maintenance, error)
	Create(ctx context.Context, maintenance Maintenance) (Maintenance, error)
	Update(ctx context.Context, maintenance Maintenance) (Maintenance, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
