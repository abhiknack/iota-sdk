package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/studio/domain/aggregates/module_definition"
	"github.com/iota-uz/iota-sdk/pkg/eventbus"
)

type ModuleDefinitionService struct {
	repo      module_definition.Repository
	publisher eventbus.EventBus
}

func NewModuleDefinitionService(repo module_definition.Repository, publisher eventbus.EventBus) *ModuleDefinitionService {
	return &ModuleDefinitionService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *ModuleDefinitionService) Create(ctx context.Context, name, displayName, description, icon string, tenantID uuid.UUID) (module_definition.ModuleDefinition, error) {
	existing, err := s.repo.GetByName(ctx, name)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("module with name %s already exists", name)
	}

	mod := module_definition.New(
		uuid.New(),
		tenantID,
		name,
		displayName,
		description,
		icon,
	)

	created, err := s.repo.Create(ctx, mod)
	if err != nil {
		return nil, fmt.Errorf("failed to create module definition: %w", err)
	}

	s.publisher.Publish(module_definition.ModuleDefinitionCreatedEvent{
		ModuleID:    created.ID(),
		TenantID:    created.TenantID(),
		Name:        created.Name(),
		DisplayName: created.DisplayName(),
		OccurredAt:  created.CreatedAt(),
	})

	return created, nil
}

func (s *ModuleDefinitionService) Update(ctx context.Context, id uuid.UUID, displayName, description, icon string) (module_definition.ModuleDefinition, error) {
	mod, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get module definition: %w", err)
	}

	mod = mod.UpdateDetails(displayName, description, icon)

	updated, err := s.repo.Update(ctx, mod)
	if err != nil {
		return nil, fmt.Errorf("failed to update module definition: %w", err)
	}

	s.publisher.Publish(module_definition.ModuleDefinitionUpdatedEvent{
		ModuleID:    updated.ID(),
		TenantID:    updated.TenantID(),
		Name:        updated.Name(),
		DisplayName: updated.DisplayName(),
		OccurredAt:  updated.UpdatedAt(),
	})

	return updated, nil
}

func (s *ModuleDefinitionService) AddEntity(ctx context.Context, moduleID uuid.UUID, entity module_definition.EntityDefinition) (module_definition.ModuleDefinition, error) {
	mod, err := s.repo.GetByID(ctx, moduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get module definition: %w", err)
	}

	mod = mod.AddEntity(entity)

	updated, err := s.repo.Update(ctx, mod)
	if err != nil {
		return nil, fmt.Errorf("failed to add entity: %w", err)
	}

	return updated, nil
}

func (s *ModuleDefinitionService) RemoveEntity(ctx context.Context, moduleID, entityID uuid.UUID) (module_definition.ModuleDefinition, error) {
	mod, err := s.repo.GetByID(ctx, moduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get module definition: %w", err)
	}

	mod = mod.RemoveEntity(entityID)

	updated, err := s.repo.Update(ctx, mod)
	if err != nil {
		return nil, fmt.Errorf("failed to remove entity: %w", err)
	}

	return updated, nil
}

func (s *ModuleDefinitionService) GetByID(ctx context.Context, id uuid.UUID) (module_definition.ModuleDefinition, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ModuleDefinitionService) GetAll(ctx context.Context) ([]module_definition.ModuleDefinition, error) {
	return s.repo.GetAll(ctx)
}

func (s *ModuleDefinitionService) GetByStatus(ctx context.Context, status module_definition.Status) ([]module_definition.ModuleDefinition, error) {
	return s.repo.GetByStatus(ctx, status)
}

func (s *ModuleDefinitionService) Delete(ctx context.Context, id uuid.UUID) error {
	mod, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get module definition: %w", err)
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete module definition: %w", err)
	}

	s.publisher.Publish(module_definition.ModuleDefinitionDeletedEvent{
		ModuleID:   mod.ID(),
		TenantID:   mod.TenantID(),
		OccurredAt: mod.UpdatedAt(),
	})

	return nil
}

func (s *ModuleDefinitionService) MarkAsGenerated(ctx context.Context, id uuid.UUID) (module_definition.ModuleDefinition, error) {
	mod, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get module definition: %w", err)
	}

	mod = mod.UpdateStatus(module_definition.StatusGenerated)

	updated, err := s.repo.Update(ctx, mod)
	if err != nil {
		return nil, fmt.Errorf("failed to mark as generated: %w", err)
	}

	s.publisher.Publish(module_definition.ModuleDefinitionGeneratedEvent{
		ModuleID:   updated.ID(),
		TenantID:   updated.TenantID(),
		Name:       updated.Name(),
		OccurredAt: updated.UpdatedAt(),
	})

	return updated, nil
}
