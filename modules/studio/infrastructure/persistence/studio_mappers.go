package persistence

import (
	"github.com/iota-uz/iota-sdk/modules/studio/domain/aggregates/module_definition"
	"github.com/iota-uz/iota-sdk/modules/studio/infrastructure/persistence/models"
)

func ToDomainModuleDefinition(model *models.ModuleDefinition) module_definition.ModuleDefinition {
	entities := make([]module_definition.EntityDefinition, len(model.Entities))
	for i, e := range model.Entities {
		fields := make([]module_definition.FieldDefinition, len(e.Fields))
		for j, f := range e.Fields {
			fields[j] = module_definition.FieldDefinition{
				ID:         f.ID,
				Name:       f.Name,
				Type:       f.Type,
				Required:   f.Required,
				Validation: f.Validation,
				Order:      f.Order,
			}
		}

		entities[i] = module_definition.EntityDefinition{
			ID:          e.ID,
			Name:        e.Name,
			DisplayName: e.DisplayName,
			Fields:      fields,
			CreatedAt:   e.CreatedAt,
		}
	}

	mod := module_definition.New(
		model.ID,
		model.TenantID,
		model.Name,
		model.DisplayName,
		model.Description,
		model.Icon,
	)

	for _, entity := range entities {
		mod = mod.AddEntity(entity)
	}

	mod = mod.UpdateStatus(module_definition.Status(model.Status))

	return mod
}

func ToDBModuleDefinition(mod module_definition.ModuleDefinition) *models.ModuleDefinition {
	entities := make(models.EntitiesJSON, len(mod.Entities()))
	for i, e := range mod.Entities() {
		fields := make([]models.FieldDefinitionJSON, len(e.Fields))
		for j, f := range e.Fields {
			fields[j] = models.FieldDefinitionJSON{
				ID:         f.ID,
				Name:       f.Name,
				Type:       f.Type,
				Required:   f.Required,
				Validation: f.Validation,
				Order:      f.Order,
			}
		}

		entities[i] = models.EntityDefinitionJSON{
			ID:          e.ID,
			Name:        e.Name,
			DisplayName: e.DisplayName,
			Fields:      fields,
			CreatedAt:   e.CreatedAt,
		}
	}

	return &models.ModuleDefinition{
		ID:          mod.ID(),
		TenantID:    mod.TenantID(),
		Name:        mod.Name(),
		DisplayName: mod.DisplayName(),
		Description: mod.Description(),
		Icon:        mod.Icon(),
		Status:      int(mod.Status()),
		Entities:    entities,
		CreatedAt:   mod.CreatedAt(),
		UpdatedAt:   mod.UpdatedAt(),
	}
}
