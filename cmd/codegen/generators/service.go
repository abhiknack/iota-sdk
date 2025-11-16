package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const serviceTemplate = `package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/{{.ModuleName}}/domain/aggregates/{{.EntityLower}}"
	"github.com/iota-uz/iota-sdk/modules/{{.ModuleName}}/permissions"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/eventbus"
)

type {{.EntityName}}Service struct {
	repo      {{.EntityLower}}.Repository
	publisher eventbus.EventBus
}

func New{{.EntityName}}Service(
	repo {{.EntityLower}}.Repository,
	publisher eventbus.EventBus,
) *{{.EntityName}}Service {
	return &{{.EntityName}}Service{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *{{.EntityName}}Service) GetByID(ctx context.Context, id uuid.UUID) ({{.EntityLower}}.{{.EntityName}}, error) {
	if err := composables.CanUser(ctx, permissions.{{.EntityName}}Read); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *{{.EntityName}}Service) GetPaginated(ctx context.Context, params *{{.EntityLower}}.FindParams) ([]{{.EntityLower}}.{{.EntityName}}, error) {
	if err := composables.CanUser(ctx, permissions.{{.EntityName}}Read); err != nil {
		return nil, err
	}
	return s.repo.GetPaginated(ctx, params)
}

func (s *{{.EntityName}}Service) Count(ctx context.Context, params *{{.EntityLower}}.FindParams) (int64, error) {
	if err := composables.CanUser(ctx, permissions.{{.EntityName}}Read); err != nil {
		return 0, err
	}
	return s.repo.Count(ctx, params)
}

func (s *{{.EntityName}}Service) Create(ctx context.Context, e {{.EntityLower}}.{{.EntityName}}) ({{.EntityLower}}.{{.EntityName}}, error) {
	if err := composables.CanUser(ctx, permissions.{{.EntityName}}Create); err != nil {
		return nil, err
	}

	created, err := s.repo.Create(ctx, e)
	if err != nil {
		return nil, fmt.Errorf("failed to create {{.EntityLower}}: %w", err)
	}

	event, err := {{.EntityLower}}.New{{.EntityName}}CreatedEvent(ctx, created)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	event.Result = created
	s.publisher.Publish(event)

	return created, nil
}

func (s *{{.EntityName}}Service) Update(ctx context.Context, e {{.EntityLower}}.{{.EntityName}}) ({{.EntityLower}}.{{.EntityName}}, error) {
	if err := composables.CanUser(ctx, permissions.{{.EntityName}}Update); err != nil {
		return nil, err
	}

	updated, err := s.repo.Update(ctx, e)
	if err != nil {
		return nil, fmt.Errorf("failed to update {{.EntityLower}}: %w", err)
	}

	event, err := {{.EntityLower}}.New{{.EntityName}}UpdatedEvent(ctx, updated)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	event.Result = updated
	s.publisher.Publish(event)

	return updated, nil
}

func (s *{{.EntityName}}Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := composables.CanUser(ctx, permissions.{{.EntityName}}Delete); err != nil {
		return err
	}

	e, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get {{.EntityLower}}: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete {{.EntityLower}}: %w", err)
	}

	event, err := {{.EntityLower}}.New{{.EntityName}}DeletedEvent(ctx)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	event.Result = e
	s.publisher.Publish(event)

	return nil
}
`

type serviceTemplateData struct {
	ModuleName  string
	EntityName  string
	EntityLower string
}

func GenerateService(moduleName, entityName string) error {
	entityLower := strings.ToLower(entityName[:1]) + entityName[1:]

	data := serviceTemplateData{
		ModuleName:  moduleName,
		EntityName:  entityName,
		EntityLower: entityLower,
	}

	basePath := filepath.Join("modules", moduleName, "services")
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	outputPath := filepath.Join(basePath, entityLower+"_service.go")
	return generateFromTemplate(serviceTemplate, outputPath, data)
}
