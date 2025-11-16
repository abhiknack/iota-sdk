package module_definition

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (ModuleDefinition, error)
	GetByName(ctx context.Context, name string) (ModuleDefinition, error)
	GetAll(ctx context.Context) ([]ModuleDefinition, error)
	GetByStatus(ctx context.Context, status Status) ([]ModuleDefinition, error)
	Create(ctx context.Context, module ModuleDefinition) (ModuleDefinition, error)
	Update(ctx context.Context, module ModuleDefinition) (ModuleDefinition, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
