package driver

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
	FieldUserID        Field = "user_id"
	FieldFirstName     Field = "first_name"
	FieldLastName      Field = "last_name"
	FieldLicenseNumber Field = "license_number"
	FieldLicenseExpiry Field = "license_expiry"
	FieldPhone         Field = "phone"
	FieldEmail         Field = "email"
	FieldStatus        Field = "status"
	FieldCreatedAt     Field = "created_at"
	FieldUpdatedAt     Field = "updated_at"
)

type FindParams struct {
	TenantID       uuid.UUID
	Status         *enums.DriverStatus
	Search         *string
	Limit          int
	Offset         int
	SortBy         Field
	SortDesc       bool
	IncludeDeleted bool
}

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (Driver, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (Driver, error)
	GetPaginated(ctx context.Context, params *FindParams) ([]Driver, error)
	Count(ctx context.Context, params *FindParams) (int64, error)
	GetExpiringLicenses(ctx context.Context, tenantID uuid.UUID, days int) ([]Driver, error)
	GetAvailable(ctx context.Context, tenantID uuid.UUID, startTime, endTime time.Time) ([]Driver, error)
	Create(ctx context.Context, driver Driver) (Driver, error)
	Update(ctx context.Context, driver Driver) (Driver, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
