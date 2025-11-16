package viewmodels

import (
	"github.com/iota-uz/iota-sdk/modules/studio/domain/aggregates/module_definition"
)

type ModuleListViewModel struct {
	PageCtx interface{}
	Modules []ModuleViewModel
}

type ModuleViewModel struct {
	ID          string
	Name        string
	DisplayName string
	Description string
	Icon        string
	Status      string
	EntityCount int
	CreatedAt   string
}

type ModuleDetailViewModel struct {
	PageCtx  interface{}
	Module   ModuleViewModel
	Entities []EntityViewModel
}

type EntityViewModel struct {
	ID          string
	Name        string
	DisplayName string
	FieldCount  int
	Fields      []FieldViewModel
}

type FieldViewModel struct {
	ID         string
	Name       string
	Type       string
	Required   bool
	Validation string
	Order      int
}

type ModuleFormViewModel struct {
	PageCtx          interface{}
	Module           *ModuleViewModel
	ValidationErrors map[string]string
}

func NewModuleListViewModel(pageCtx interface{}, modules []module_definition.ModuleDefinition) *ModuleListViewModel {
	vms := make([]ModuleViewModel, len(modules))
	for i, mod := range modules {
		vms[i] = ModuleViewModel{
			ID:          mod.ID().String(),
			Name:        mod.Name(),
			DisplayName: mod.DisplayName(),
			Description: mod.Description(),
			Icon:        mod.Icon(),
			Status:      mod.Status().String(),
			EntityCount: len(mod.Entities()),
			CreatedAt:   mod.CreatedAt().Format("2006-01-02"),
		}
	}

	return &ModuleListViewModel{
		PageCtx: pageCtx,
		Modules: vms,
	}
}

func NewModuleDetailViewModel(pageCtx interface{}, mod module_definition.ModuleDefinition) *ModuleDetailViewModel {
	entities := make([]EntityViewModel, len(mod.Entities()))
	for i, entity := range mod.Entities() {
		fields := make([]FieldViewModel, len(entity.Fields))
		for j, field := range entity.Fields {
			fields[j] = FieldViewModel{
				ID:         field.ID.String(),
				Name:       field.Name,
				Type:       field.Type,
				Required:   field.Required,
				Validation: field.Validation,
				Order:      field.Order,
			}
		}

		entities[i] = EntityViewModel{
			ID:          entity.ID.String(),
			Name:        entity.Name,
			DisplayName: entity.DisplayName,
			FieldCount:  len(entity.Fields),
			Fields:      fields,
		}
	}

	return &ModuleDetailViewModel{
		PageCtx: pageCtx,
		Module: ModuleViewModel{
			ID:          mod.ID().String(),
			Name:        mod.Name(),
			DisplayName: mod.DisplayName(),
			Description: mod.Description(),
			Icon:        mod.Icon(),
			Status:      mod.Status().String(),
			EntityCount: len(mod.Entities()),
			CreatedAt:   mod.CreatedAt().Format("2006-01-02"),
		},
		Entities: entities,
	}
}

func NewModuleFormViewModel(pageCtx interface{}, mod module_definition.ModuleDefinition, validationErrors map[string]string) *ModuleFormViewModel {
	var vm *ModuleViewModel
	if mod != nil {
		vm = &ModuleViewModel{
			ID:          mod.ID().String(),
			Name:        mod.Name(),
			DisplayName: mod.DisplayName(),
			Description: mod.Description(),
			Icon:        mod.Icon(),
			Status:      mod.Status().String(),
			EntityCount: len(mod.Entities()),
			CreatedAt:   mod.CreatedAt().Format("2006-01-02"),
		}
	}

	return &ModuleFormViewModel{
		PageCtx:          pageCtx,
		Module:           vm,
		ValidationErrors: validationErrors,
	}
}
