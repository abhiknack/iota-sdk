package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const repositoryTemplate = `package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/{{.ModuleName}}/domain/aggregates/{{.EntityLower}}"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/repo"
)

var (
	Err{{.EntityName}}NotFound = errors.New("{{.EntityLower}} not found")
)

const (
	select{{.EntityName}}Query = ` + "`" + `
		SELECT
			id,
			tenant_id,
{{- range .Fields}}
			{{.NameSnake}},
{{- end}}
			created_at,
			updated_at
		FROM {{.TableName}}
	` + "`" + `
	count{{.EntityName}}Query  = ` + "`SELECT COUNT(*) FROM {{.TableName}}`" + `
	delete{{.EntityName}}Query = ` + "`UPDATE {{.TableName}} SET deleted_at = $1 WHERE id = $2 AND tenant_id = $3`" + `
)

type {{.EntityName}}Repository struct {
	fieldMap map[{{.EntityLower}}.Field]string
}

func New{{.EntityName}}Repository() {{.EntityLower}}.Repository {
	return &{{.EntityName}}Repository{
		fieldMap: map[{{.EntityLower}}.Field]string{
			{{.EntityLower}}.FieldID:        "id",
			{{.EntityLower}}.FieldTenantID:  "tenant_id",
{{- range .Fields}}
			{{$.EntityLower}}.Field{{.Name}}: "{{.NameSnake}}",
{{- end}}
			{{.EntityLower}}.FieldCreatedAt: "created_at",
			{{.EntityLower}}.FieldUpdatedAt: "updated_at",
		},
	}
}

func (r *{{.EntityName}}Repository) query{{.EntityName}}s(
	ctx context.Context,
	query string,
	args ...interface{},
) ([]{{.EntityLower}}.{{.EntityName}}, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make([]{{.EntityLower}}.{{.EntityName}}, 0)
	for rows.Next() {
		var id, tenantID uuid.UUID
{{- range .Fields}}
		var {{.NameLower}} {{.Type}}
{{- end}}
		var createdAt, updatedAt time.Time

		if err := rows.Scan(
			&id,
			&tenantID,
{{- range .Fields}}
			&{{.NameLower}},
{{- end}}
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		e := {{.EntityLower}}.New{{$.EntityName}}(
			id,
			tenantID,
{{- range .Fields}}
			{{.NameLower}},
{{- end}}
			{{.EntityLower}}.WithTimestamps(createdAt, updatedAt),
		)

		entities = append(entities, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *{{.EntityName}}Repository) GetByID(ctx context.Context, id uuid.UUID) ({{.EntityLower}}.{{.EntityName}}, error) {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Join(select{{.EntityName}}Query, "WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL")
	entities, err := r.query{{.EntityName}}s(ctx, query, id, tenantID)
	if err != nil {
		return nil, err
	}
	if len(entities) == 0 {
		return nil, Err{{.EntityName}}NotFound
	}
	return entities[0], nil
}

func (r *{{.EntityName}}Repository) GetPaginated(ctx context.Context, params *{{.EntityLower}}.FindParams) ([]{{.EntityLower}}.{{.EntityName}}, error) {
	where := make([]string, 0)
	args := make([]interface{}, 0)

	where = append(where, fmt.Sprintf("tenant_id = $%d", len(args)+1))
	args = append(args, params.TenantID)

	if !params.IncludeDeleted {
		where = append(where, "deleted_at IS NULL")
	}

	if params.Search != nil && *params.Search != "" {
		searchPlaceholder := fmt.Sprintf("$%d", len(args)+1)
		where = append(where, fmt.Sprintf("(CAST(id AS TEXT) ILIKE %s)", searchPlaceholder))
		args = append(args, "%"+*params.Search+"%")
	}

	sortColumn := r.fieldMap[params.SortBy]
	if sortColumn == "" {
		sortColumn = "created_at"
	}
	sortDir := "ASC"
	if params.SortDesc {
		sortDir = "DESC"
	}

	query := repo.Join(
		select{{.EntityName}}Query,
		repo.JoinWhere(where...),
		fmt.Sprintf("ORDER BY %s %s", sortColumn, sortDir),
		repo.FormatLimitOffset(params.Limit, params.Offset),
	)

	return r.query{{.EntityName}}s(ctx, query, args...)
}

func (r *{{.EntityName}}Repository) Count(ctx context.Context, params *{{.EntityLower}}.FindParams) (int64, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return 0, err
	}

	where := make([]string, 0)
	args := make([]interface{}, 0)

	where = append(where, fmt.Sprintf("tenant_id = $%d", len(args)+1))
	args = append(args, params.TenantID)

	if !params.IncludeDeleted {
		where = append(where, "deleted_at IS NULL")
	}

	if params.Search != nil && *params.Search != "" {
		searchPlaceholder := fmt.Sprintf("$%d", len(args)+1)
		where = append(where, fmt.Sprintf("(CAST(id AS TEXT) ILIKE %s)", searchPlaceholder))
		args = append(args, "%"+*params.Search+"%")
	}

	query := repo.Join(count{{.EntityName}}Query, repo.JoinWhere(where...))

	var count int64
	if err := pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *{{.EntityName}}Repository) Create(ctx context.Context, e {{.EntityLower}}.{{.EntityName}}) ({{.EntityLower}}.{{.EntityName}}, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Insert(
		"{{.TableName}}",
		[]string{
			"id",
			"tenant_id",
{{- range .Fields}}
			"{{.NameSnake}}",
{{- end}}
			"created_at",
			"updated_at",
		},
		"",
	)

	_, err = pool.Exec(
		ctx,
		query,
		e.ID(),
		e.TenantID(),
{{- range .Fields}}
		e.{{.Name}}(),
{{- end}}
		e.CreatedAt(),
		e.UpdatedAt(),
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, e.ID())
}

func (r *{{.EntityName}}Repository) Update(ctx context.Context, e {{.EntityLower}}.{{.EntityName}}) ({{.EntityLower}}.{{.EntityName}}, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Update(
		"{{.TableName}}",
		[]string{
{{- range .Fields}}
			"{{.NameSnake}}",
{{- end}}
			"updated_at",
		},
		fmt.Sprintf("id = $%d AND tenant_id = $%d", {{.UpdateParamCount}}, {{.UpdateParamCountPlus1}}),
	)

	_, err = pool.Exec(
		ctx,
		query,
{{- range .Fields}}
		e.{{.Name}}(),
{{- end}}
		time.Now(),
		e.ID(),
		e.TenantID(),
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, e.ID())
}

func (r *{{.EntityName}}Repository) Delete(ctx context.Context, id uuid.UUID) error {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return err
	}

	pool, err := composables.UseTx(ctx)
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, delete{{.EntityName}}Query, time.Now(), id, tenantID)
	return err
}
`

type repositoryTemplateData struct {
	ModuleName            string
	EntityName            string
	EntityLower           string
	TableName             string
	Fields                []repoFieldTemplateData
	UpdateParamCount      int
	UpdateParamCountPlus1 int
}

type repoFieldTemplateData struct {
	Name      string
	NameLower string
	NameSnake string
	Type      string
}

func GenerateRepository(moduleName, entityName string, fields []Field) error {
	entityLower := strings.ToLower(entityName[:1]) + entityName[1:]
	tableName := moduleName + "_" + toSnakeCase(entityName) + "s"

	data := repositoryTemplateData{
		ModuleName:            moduleName,
		EntityName:            entityName,
		EntityLower:           entityLower,
		TableName:             tableName,
		Fields:                make([]repoFieldTemplateData, len(fields)),
		UpdateParamCount:      len(fields) + 1,
		UpdateParamCountPlus1: len(fields) + 2,
	}

	for i, f := range fields {
		data.Fields[i] = repoFieldTemplateData{
			Name:      f.Name,
			NameLower: strings.ToLower(f.Name[:1]) + f.Name[1:],
			NameSnake: toSnakeCase(f.Name),
			Type:      f.Type,
		}
	}

	basePath := filepath.Join("modules", moduleName, "infrastructure", "persistence")
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	outputPath := filepath.Join(basePath, entityLower+"_repository.go")
	return generateFromTemplate(repositoryTemplate, outputPath, data)
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
