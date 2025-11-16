package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/cmd/codegen/generators"
	"github.com/iota-uz/iota-sdk/modules/studio/domain/aggregates/module_definition"
)

type CodeGeneratorService struct {
	moduleDefService *ModuleDefinitionService
}

func NewCodeGeneratorService(moduleDefService *ModuleDefinitionService) *CodeGeneratorService {
	return &CodeGeneratorService{
		moduleDefService: moduleDefService,
	}
}

func (s *CodeGeneratorService) GenerateModule(ctx context.Context, moduleID uuid.UUID) error {
	mod, err := s.moduleDefService.GetByID(ctx, moduleID)
	if err != nil {
		return fmt.Errorf("failed to get module definition: %w", err)
	}

	moduleName := mod.Name()

	moduleDir := filepath.Join("modules", moduleName)
	if _, err := os.Stat(moduleDir); !os.IsNotExist(err) {
		return fmt.Errorf("module directory already exists: %s", moduleDir)
	}

	if err := os.MkdirAll(moduleDir, 0755); err != nil {
		return fmt.Errorf("failed to create module directory: %w", err)
	}

	for _, entity := range mod.Entities() {
		fields := s.convertFields(entity.Fields)

		if err := generators.GenerateCRUD(moduleName, entity.Name, fields); err != nil {
			return fmt.Errorf("failed to generate CRUD for entity %s: %w", entity.Name, err)
		}
	}

	if err := s.generateModuleFile(mod); err != nil {
		return fmt.Errorf("failed to generate module file: %w", err)
	}

	if err := s.generateLinksFile(mod); err != nil {
		return fmt.Errorf("failed to generate links file: %w", err)
	}

	if err := s.generatePermissionsFile(mod); err != nil {
		return fmt.Errorf("failed to generate permissions file: %w", err)
	}

	if err := s.generateLocaleFiles(mod); err != nil {
		return fmt.Errorf("failed to generate locale files: %w", err)
	}

	if _, err := s.moduleDefService.MarkAsGenerated(ctx, moduleID); err != nil {
		return fmt.Errorf("failed to mark module as generated: %w", err)
	}

	return nil
}

func (s *CodeGeneratorService) convertFields(fields []module_definition.FieldDefinition) []generators.Field {
	result := make([]generators.Field, len(fields))
	for i, f := range fields {
		result[i] = generators.Field{
			Name:       f.Name,
			Type:       f.Type,
			Validation: f.Validation,
		}
	}
	return result
}

func (s *CodeGeneratorService) generateModuleFile(mod module_definition.ModuleDefinition) error {
	moduleName := mod.Name()
	moduleDir := filepath.Join("modules", moduleName)

	content := fmt.Sprintf(`package %s

import (
	"github.com/iota-uz/iota-sdk/pkg/application"
	"github.com/iota-uz/iota-sdk/pkg/event"
)

type Module struct{}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) Register(app application.Application) error {
	eventPublisher := event.NewPublisher()
	
	// TODO: Register repositories, services, and controllers
	
	return nil
}

func (m *Module) Name() string {
	return "%s"
}
`, moduleName, mod.DisplayName())

	filePath := filepath.Join(moduleDir, "module.go")
	return os.WriteFile(filePath, []byte(content), 0644)
}

func (s *CodeGeneratorService) generateLinksFile(mod module_definition.ModuleDefinition) error {
	moduleName := mod.Name()
	moduleDir := filepath.Join("modules", moduleName)

	content := fmt.Sprintf(`package %s

import (
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/types"
)

func NavItems(ctx *composables.PageCtx) []types.NavigationLink {
	return []types.NavigationLink{
		{
			Title: ctx.T("%s.NavigationLinks.Dashboard"),
			Href:  "/%s/dashboard",
			Icon:  "%s",
		},
	}
}
`, moduleName, mod.DisplayName(), moduleName, mod.Icon())

	filePath := filepath.Join(moduleDir, "links.go")
	return os.WriteFile(filePath, []byte(content), 0644)
}

func (s *CodeGeneratorService) generatePermissionsFile(mod module_definition.ModuleDefinition) error {
	moduleName := mod.Name()
	moduleDir := filepath.Join("modules", moduleName)
	permissionsDir := filepath.Join(moduleDir, "permissions")

	if err := os.MkdirAll(permissionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create permissions directory: %w", err)
	}

	var permissions []string
	permissions = append(permissions, fmt.Sprintf("\tView = \"%s.view\"", moduleName))
	permissions = append(permissions, fmt.Sprintf("\tManage = \"%s.manage\"", moduleName))

	for _, entity := range mod.Entities() {
		entityLower := strings.ToLower(entity.Name)
		permissions = append(permissions, fmt.Sprintf("\t%sCreate = \"%s.%s.create\"", entity.Name, moduleName, entityLower))
		permissions = append(permissions, fmt.Sprintf("\t%sUpdate = \"%s.%s.update\"", entity.Name, moduleName, entityLower))
		permissions = append(permissions, fmt.Sprintf("\t%sDelete = \"%s.%s.delete\"", entity.Name, moduleName, entityLower))
	}

	content := fmt.Sprintf(`package permissions

const (
%s
)
`, strings.Join(permissions, "\n"))

	filePath := filepath.Join(permissionsDir, "constants.go")
	return os.WriteFile(filePath, []byte(content), 0644)
}

func (s *CodeGeneratorService) generateLocaleFiles(mod module_definition.ModuleDefinition) error {
	moduleName := mod.Name()
	moduleDir := filepath.Join("modules", moduleName)
	localesDir := filepath.Join(moduleDir, "presentation", "locales")

	if err := os.MkdirAll(localesDir, 0755); err != nil {
		return fmt.Errorf("failed to create locales directory: %w", err)
	}

	enContent := fmt.Sprintf(`{
  "%s": {
    "NavigationLinks": {
      "Dashboard": "%s Dashboard"
    },
    "Meta": {
      "Title": "%s"
    }
  }
}
`, mod.DisplayName(), mod.DisplayName(), mod.DisplayName())

	enPath := filepath.Join(localesDir, "en.json")
	if err := os.WriteFile(enPath, []byte(enContent), 0644); err != nil {
		return err
	}

	ruPath := filepath.Join(localesDir, "ru.json")
	if err := os.WriteFile(ruPath, []byte(enContent), 0644); err != nil {
		return err
	}

	uzPath := filepath.Join(localesDir, "uz.json")
	if err := os.WriteFile(uzPath, []byte(enContent), 0644); err != nil {
		return err
	}

	return nil
}

func (s *CodeGeneratorService) PreviewCode(ctx context.Context, moduleID uuid.UUID) (map[string]string, error) {
	mod, err := s.moduleDefService.GetByID(ctx, moduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get module definition: %w", err)
	}

	preview := make(map[string]string)

	preview["module.go"] = fmt.Sprintf("// Module: %s\n// Description: %s\n// Entities: %d",
		mod.Name(), mod.Description(), len(mod.Entities()))

	for _, entity := range mod.Entities() {
		preview[fmt.Sprintf("%s_entity.go", strings.ToLower(entity.Name))] =
			fmt.Sprintf("// Entity: %s\n// Fields: %d", entity.Name, len(entity.Fields))
	}

	return preview, nil
}
