package permissions

import (
	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/core/domain/entities/permission"
)

const (
	ResourceModuleDefinition permission.Resource = "module_definition"
)

var (
	ModuleDefinitionCreate = &permission.Permission{
		ID:       uuid.MustParse("00000000-0000-0000-0000-000000000a01"),
		Name:     "ModuleDefinition.Create",
		Resource: ResourceModuleDefinition,
		Action:   permission.ActionCreate,
		Modifier: permission.ModifierAll,
	}
	ModuleDefinitionRead = &permission.Permission{
		ID:       uuid.MustParse("00000000-0000-0000-0000-000000000a02"),
		Name:     "ModuleDefinition.Read",
		Resource: ResourceModuleDefinition,
		Action:   permission.ActionRead,
		Modifier: permission.ModifierAll,
	}
	ModuleDefinitionUpdate = &permission.Permission{
		ID:       uuid.MustParse("00000000-0000-0000-0000-000000000a03"),
		Name:     "ModuleDefinition.Update",
		Resource: ResourceModuleDefinition,
		Action:   permission.ActionUpdate,
		Modifier: permission.ModifierAll,
	}
	ModuleDefinitionDelete = &permission.Permission{
		ID:       uuid.MustParse("00000000-0000-0000-0000-000000000a04"),
		Name:     "ModuleDefinition.Delete",
		Resource: ResourceModuleDefinition,
		Action:   permission.ActionDelete,
		Modifier: permission.ModifierAll,
	}
)

var Permissions = []*permission.Permission{
	ModuleDefinitionCreate,
	ModuleDefinitionRead,
	ModuleDefinitionUpdate,
	ModuleDefinitionDelete,
}
