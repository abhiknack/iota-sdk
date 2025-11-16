package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const dtoTemplate = `package dtos

import (
	"github.com/google/uuid"
)

type {{.EntityName}}CreateDTO struct {
{{- range .Fields}}
	{{.Name}} {{.Type}} ` + "`form:\"{{.Name}}\" validate:\"{{.Validation}}\"`" + `
{{- end}}
}

type {{.EntityName}}UpdateDTO struct {
	ID uuid.UUID ` + "`form:\"ID\" validate:\"required\"`" + `
{{- range .Fields}}
	{{.Name}} {{.Type}} ` + "`form:\"{{.Name}}\" validate:\"{{.Validation}}\"`" + `
{{- end}}
}

type {{.EntityName}}FilterDTO struct {
{{- range .Fields}}
	{{.Name}} {{.Type}} ` + "`form:\"{{.Name}}\"`" + `
{{- end}}
	Page     int ` + "`form:\"Page\"`" + `
	PageSize int ` + "`form:\"PageSize\"`" + `
}
`

type dtoTemplateData struct {
	EntityName string
	Fields     []dtoFieldTemplateData
}

type dtoFieldTemplateData struct {
	Name       string
	Type       string
	Validation string
}

func GenerateDTO(moduleName, entityName string, fields []Field) error {
	data := dtoTemplateData{
		EntityName: entityName,
		Fields:     make([]dtoFieldTemplateData, len(fields)),
	}

	for i, f := range fields {
		validation := f.Validation
		if validation == "" {
			validation = "required"
		}
		data.Fields[i] = dtoFieldTemplateData{
			Name:       f.Name,
			Type:       f.Type,
			Validation: validation,
		}
	}

	basePath := filepath.Join("modules", moduleName, "presentation", "controllers", "dtos")
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	entityLower := strings.ToLower(entityName[:1]) + entityName[1:]
	outputPath := filepath.Join(basePath, entityLower+"_dto.go")
	return generateFromTemplate(dtoTemplate, outputPath, data)
}
