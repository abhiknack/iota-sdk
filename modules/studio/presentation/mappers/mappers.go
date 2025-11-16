package mappers

import (
	"time"

	"github.com/iota-uz/iota-sdk/modules/studio/domain/aggregates/module_definition"
	"github.com/iota-uz/iota-sdk/modules/studio/presentation/controllers/dtos"
)

func ToEntityDefinition(dto *dtos.EntityDefinitionDTO) module_definition.EntityDefinition {
	fields := make([]module_definition.FieldDefinition, len(dto.Fields))
	for i, f := range dto.Fields {
		fields[i] = module_definition.FieldDefinition{
			ID:         f.ID,
			Name:       f.Name,
			Type:       f.Type,
			Required:   f.Required,
			Validation: f.Validation,
			Order:      f.Order,
		}
	}

	return module_definition.EntityDefinition{
		ID:          dto.ID,
		Name:        dto.Name,
		DisplayName: dto.DisplayName,
		Fields:      fields,
		CreatedAt:   time.Now(),
	}
}

func ToEntityDefinitionDTO(entity module_definition.EntityDefinition) *dtos.EntityDefinitionDTO {
	fields := make([]dtos.FieldDefinitionDTO, len(entity.Fields))
	for i, f := range entity.Fields {
		fields[i] = dtos.FieldDefinitionDTO{
			ID:         f.ID,
			Name:       f.Name,
			Type:       f.Type,
			Required:   f.Required,
			Validation: f.Validation,
			Order:      f.Order,
		}
	}

	return &dtos.EntityDefinitionDTO{
		ID:          entity.ID,
		Name:        entity.Name,
		DisplayName: entity.DisplayName,
		Fields:      fields,
	}
}
