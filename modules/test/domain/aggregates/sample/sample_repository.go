package sample

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (Sample, error)
	GetPaginated(ctx context.Context, params *FindParams) ([]Sample, error)
	Count(ctx context.Context, params *FindParams) (int64, error)
	Create(ctx context.Context, entity Sample) (Sample, error)
	Update(ctx context.Context, entity Sample) (Sample, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Field int

const (
	FieldID Field = iota
	FieldTenantID
	FieldName
	FieldAge
	FieldCreatedAt
	FieldUpdatedAt
)

type FindParams struct {
	TenantID       uuid.UUID
	Limit          int
	Offset         int
	SortBy         Field
	SortDesc       bool
	IncludeDeleted bool
	Search         *string
}
