package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/studio/domain/aggregates/module_definition"
	"github.com/iota-uz/iota-sdk/modules/studio/infrastructure/persistence/models"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/repo"
)

const (
	moduleDefinitionSelectQuery = `
		SELECT id, tenant_id, name, display_name, description, icon, status, entities, created_at, updated_at
		FROM studio_module_definitions
	`

	moduleDefinitionCreateQuery = `
		INSERT INTO studio_module_definitions (id, tenant_id, name, display_name, description, icon, status, entities, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, tenant_id, name, display_name, description, icon, status, entities, created_at, updated_at
	`

	moduleDefinitionUpdateQuery = `
		UPDATE studio_module_definitions
		SET display_name = $1, description = $2, icon = $3, status = $4, entities = $5, updated_at = $6
		WHERE id = $7 AND tenant_id = $8 AND deleted_at IS NULL
		RETURNING id, tenant_id, name, display_name, description, icon, status, entities, created_at, updated_at
	`

	moduleDefinitionDeleteQuery = `
		UPDATE studio_module_definitions
		SET deleted_at = NOW()
		WHERE id = $1 AND tenant_id = $2
	`
)

type PgModuleDefinitionRepository struct{}

func NewModuleDefinitionRepository() module_definition.Repository {
	return &PgModuleDefinitionRepository{}
}

func (r *PgModuleDefinitionRepository) queryModuleDefinitions(
	ctx context.Context,
	query string,
	args ...interface{},
) ([]module_definition.ModuleDefinition, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	modules := make([]module_definition.ModuleDefinition, 0)
	for rows.Next() {
		var model models.ModuleDefinition
		if err := rows.Scan(
			&model.ID,
			&model.TenantID,
			&model.Name,
			&model.DisplayName,
			&model.Description,
			&model.Icon,
			&model.Status,
			&model.Entities,
			&model.CreatedAt,
			&model.UpdatedAt,
		); err != nil {
			return nil, err
		}
		modules = append(modules, ToDomainModuleDefinition(&model))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return modules, nil
}

func (r *PgModuleDefinitionRepository) GetByID(ctx context.Context, id uuid.UUID) (module_definition.ModuleDefinition, error) {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Join(moduleDefinitionSelectQuery, "WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL")
	modules, err := r.queryModuleDefinitions(ctx, query, id, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get module definition: %w", err)
	}

	if len(modules) == 0 {
		return nil, fmt.Errorf("module definition not found")
	}

	return modules[0], nil
}

func (r *PgModuleDefinitionRepository) GetByName(ctx context.Context, name string) (module_definition.ModuleDefinition, error) {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Join(moduleDefinitionSelectQuery, "WHERE name = $1 AND tenant_id = $2 AND deleted_at IS NULL")
	modules, err := r.queryModuleDefinitions(ctx, query, name, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get module definition: %w", err)
	}

	if len(modules) == 0 {
		return nil, fmt.Errorf("module definition not found")
	}

	return modules[0], nil
}

func (r *PgModuleDefinitionRepository) GetAll(ctx context.Context) ([]module_definition.ModuleDefinition, error) {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Join(
		moduleDefinitionSelectQuery,
		"WHERE tenant_id = $1 AND deleted_at IS NULL",
		"ORDER BY created_at DESC",
	)

	return r.queryModuleDefinitions(ctx, query, tenantID)
}

func (r *PgModuleDefinitionRepository) GetByStatus(ctx context.Context, status module_definition.Status) ([]module_definition.ModuleDefinition, error) {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Join(
		moduleDefinitionSelectQuery,
		"WHERE tenant_id = $1 AND status = $2 AND deleted_at IS NULL",
		"ORDER BY created_at DESC",
	)

	return r.queryModuleDefinitions(ctx, query, tenantID, int(status))
}

func (r *PgModuleDefinitionRepository) Create(ctx context.Context, mod module_definition.ModuleDefinition) (module_definition.ModuleDefinition, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	dbModel := ToDBModuleDefinition(mod)

	var result models.ModuleDefinition
	err = pool.QueryRow(
		ctx,
		moduleDefinitionCreateQuery,
		dbModel.ID,
		dbModel.TenantID,
		dbModel.Name,
		dbModel.DisplayName,
		dbModel.Description,
		dbModel.Icon,
		dbModel.Status,
		dbModel.Entities,
		dbModel.CreatedAt,
		dbModel.UpdatedAt,
	).Scan(
		&result.ID,
		&result.TenantID,
		&result.Name,
		&result.DisplayName,
		&result.Description,
		&result.Icon,
		&result.Status,
		&result.Entities,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create module definition: %w", err)
	}

	return ToDomainModuleDefinition(&result), nil
}

func (r *PgModuleDefinitionRepository) Update(ctx context.Context, mod module_definition.ModuleDefinition) (module_definition.ModuleDefinition, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	dbModel := ToDBModuleDefinition(mod)

	var result models.ModuleDefinition
	err = pool.QueryRow(
		ctx,
		moduleDefinitionUpdateQuery,
		dbModel.DisplayName,
		dbModel.Description,
		dbModel.Icon,
		dbModel.Status,
		dbModel.Entities,
		dbModel.UpdatedAt,
		dbModel.ID,
		tenantID,
	).Scan(
		&result.ID,
		&result.TenantID,
		&result.Name,
		&result.DisplayName,
		&result.Description,
		&result.Icon,
		&result.Status,
		&result.Entities,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("module definition not found")
		}
		return nil, fmt.Errorf("failed to update module definition: %w", err)
	}

	return ToDomainModuleDefinition(&result), nil
}

func (r *PgModuleDefinitionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return err
	}

	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, moduleDefinitionDeleteQuery, id, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete module definition: %w", err)
	}

	return nil
}
