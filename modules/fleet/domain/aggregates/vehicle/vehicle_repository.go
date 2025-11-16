package vehicle

import (
	"context"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
)

type Field string

const (
	FieldID                 Field = "id"
	FieldTenantID           Field = "tenant_id"
	FieldMake               Field = "make"
	FieldModel              Field = "model"
	FieldYear               Field = "year"
	FieldVIN                Field = "vin"
	FieldLicensePlate       Field = "license_plate"
	FieldStatus             Field = "status"
	FieldCurrentOdometer    Field = "current_odometer"
	FieldRegistrationExpiry Field = "registration_expiry"
	FieldInsuranceExpiry    Field = "insurance_expiry"
	FieldCreatedAt          Field = "created_at"
	FieldUpdatedAt          Field = "updated_at"
)

type FindParams struct {
	TenantID       uuid.UUID
	Status         *enums.VehicleStatus
	Make           *string
	Model          *string
	Year           *int
	Search         *string
	Limit          int
	Offset         int
	SortBy         Field
	SortDesc       bool
	IncludeDeleted bool
}

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (Vehicle, error)
	GetPaginated(ctx context.Context, params *FindParams) ([]Vehicle, error)
	Count(ctx context.Context, params *FindParams) (int64, error)
	GetByStatus(ctx context.Context, tenantID uuid.UUID, status enums.VehicleStatus) ([]Vehicle, error)
	GetExpiringRegistrations(ctx context.Context, tenantID uuid.UUID, days int) ([]Vehicle, error)
	Create(ctx context.Context, vehicle Vehicle) (Vehicle, error)
	Update(ctx context.Context, vehicle Vehicle) (Vehicle, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
