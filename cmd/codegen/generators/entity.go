package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const entityTemplate = `package {{.EntityLower}}

import (
	"time"

	"github.com/google/uuid"
)

type {{.EntityName}} interface {
	ID() uuid.UUID
	TenantID() uuid.UUID
{{- range .Fields}}
	{{.Name}}() {{.Type}}
{{- end}}
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type {{.EntityName}}Option func(*{{.EntityLower}})

func New{{.EntityName}}(
	id uuid.UUID,
	tenantID uuid.UUID,
{{- range .Fields}}
	{{.NameLower}} {{.Type}},
{{- end}}
	opts ...{{.EntityName}}Option,
) {{.EntityName}} {
	e := &{{.EntityLower}}{
		id:        id,
		tenantID:  tenantID,
{{- range .Fields}}
		{{.NameLower}}: {{.NameLower}},
{{- end}}
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func WithTimestamps(createdAt, updatedAt time.Time) {{.EntityName}}Option {
	return func(e *{{.EntityLower}}) {
		e.createdAt = createdAt
		e.updatedAt = updatedAt
	}
}

type {{.EntityLower}} struct {
	id        uuid.UUID
	tenantID  uuid.UUID
{{- range .Fields}}
	{{.NameLower}} {{.Type}}
{{- end}}
	createdAt time.Time
	updatedAt time.Time
}

func (e *{{.EntityLower}}) ID() uuid.UUID {
	return e.id
}

func (e *{{.EntityLower}}) TenantID() uuid.UUID {
	return e.tenantID
}

{{- range .Fields}}

func (e *{{$.EntityLower}}) {{.Name}}() {{.Type}} {
	return e.{{.NameLower}}
}
{{- end}}

func (e *{{.EntityLower}}) CreatedAt() time.Time {
	return e.createdAt
}

func (e *{{.EntityLower}}) UpdatedAt() time.Time {
	return e.updatedAt
}
`

const repositoryInterfaceTemplate = `package {{.EntityLower}}

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) ({{.EntityName}}, error)
	GetPaginated(ctx context.Context, params *FindParams) ([]{{.EntityName}}, error)
	Count(ctx context.Context, params *FindParams) (int64, error)
	Create(ctx context.Context, entity {{.EntityName}}) ({{.EntityName}}, error)
	Update(ctx context.Context, entity {{.EntityName}}) ({{.EntityName}}, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Field int

const (
	FieldID Field = iota
	FieldTenantID
{{- range $i, $field := .Fields}}
	Field{{$field.Name}}
{{- end}}
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
`

const eventsTemplate = `package {{.EntityLower}}

import (
	"context"

	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/event"
)

const (
	{{.EntityName}}CreatedEvent = "{{.EntityLower}}.created"
	{{.EntityName}}UpdatedEvent = "{{.EntityLower}}.updated"
	{{.EntityName}}DeletedEvent = "{{.EntityLower}}.deleted"
)

func New{{.EntityName}}CreatedEvent(ctx context.Context, entity {{.EntityName}}) (*event.Event, error) {
	userID, _ := composables.UseUserID(ctx)
	tenantID, _ := composables.UseTenantID(ctx)

	return &event.Event{
		Type:     {{.EntityName}}CreatedEvent,
		TenantID: tenantID,
		UserID:   userID,
		Payload:  entity,
	}, nil
}

func New{{.EntityName}}UpdatedEvent(ctx context.Context, entity {{.EntityName}}) (*event.Event, error) {
	userID, _ := composables.UseUserID(ctx)
	tenantID, _ := composables.UseTenantID(ctx)

	return &event.Event{
		Type:     {{.EntityName}}UpdatedEvent,
		TenantID: tenantID,
		UserID:   userID,
		Payload:  entity,
	}, nil
}

func New{{.EntityName}}DeletedEvent(ctx context.Context) (*event.Event, error) {
	userID, _ := composables.UseUserID(ctx)
	tenantID, _ := composables.UseTenantID(ctx)

	return &event.Event{
		Type:     {{.EntityName}}DeletedEvent,
		TenantID: tenantID,
		UserID:   userID,
	}, nil
}
`

type entityTemplateData struct {
	ModuleName  string
	EntityName  string
	EntityLower string
	Fields      []fieldTemplateData
}

type fieldTemplateData struct {
	Name      string
	NameLower string
	Type      string
}

func GenerateEntity(moduleName, entityName string, fields []Field) error {
	entityLower := strings.ToLower(entityName[:1]) + entityName[1:]

	data := entityTemplateData{
		ModuleName:  moduleName,
		EntityName:  entityName,
		EntityLower: entityLower,
		Fields:      make([]fieldTemplateData, len(fields)),
	}

	for i, f := range fields {
		data.Fields[i] = fieldTemplateData{
			Name:      f.Name,
			NameLower: strings.ToLower(f.Name[:1]) + f.Name[1:],
			Type:      f.Type,
		}
	}

	basePath := filepath.Join("modules", moduleName, "domain", "aggregates", entityLower)
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := generateFromTemplate(entityTemplate, filepath.Join(basePath, entityLower+".go"), data); err != nil {
		return err
	}

	if err := generateFromTemplate(repositoryInterfaceTemplate, filepath.Join(basePath, entityLower+"_repository.go"), data); err != nil {
		return err
	}

	if err := generateFromTemplate(eventsTemplate, filepath.Join(basePath, entityLower+"_events.go"), data); err != nil {
		return err
	}

	return nil
}

func generateFromTemplate(tmplStr, outputPath string, data interface{}) error {
	tmpl, err := template.New("gen").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
