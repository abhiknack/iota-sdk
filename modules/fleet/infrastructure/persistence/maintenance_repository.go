package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/maintenance"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/repo"
)

var (
	ErrMaintenanceNotFound = errors.New("maintenance not found")
)

const (
	selectMaintenanceQuery = `
		SELECT
			id,
			tenant_id,
			vehicle_id,
			service_type,
			service_date,
			odometer,
			cost,
			service_provider,
			description,
			next_service_due,
			next_service_odometer,
			created_at,
			updated_at
		FROM fleet_maintenance
	`
	countMaintenanceQuery  = `SELECT COUNT(*) FROM fleet_maintenance`
	deleteMaintenanceQuery = `UPDATE fleet_maintenance SET deleted_at = $1 WHERE id = $2 AND tenant_id = $3`
)

type MaintenanceRepository struct {
	fieldMap map[maintenance.Field]string
}

func NewMaintenanceRepository() maintenance.Repository {
	return &MaintenanceRepository{
		fieldMap: map[maintenance.Field]string{
			maintenance.FieldID:                  "id",
			maintenance.FieldTenantID:            "tenant_id",
			maintenance.FieldVehicleID:           "vehicle_id",
			maintenance.FieldServiceType:         "service_type",
			maintenance.FieldServiceDate:         "service_date",
			maintenance.FieldOdometer:            "odometer",
			maintenance.FieldCost:                "cost",
			maintenance.FieldServiceProvider:     "service_provider",
			maintenance.FieldDescription:         "description",
			maintenance.FieldNextServiceDue:      "next_service_due",
			maintenance.FieldNextServiceOdometer: "next_service_odometer",
			maintenance.FieldCreatedAt:           "created_at",
			maintenance.FieldUpdatedAt:           "updated_at",
		},
	}
}

func (r *MaintenanceRepository) queryMaintenance(
	ctx context.Context,
	query string,
	args ...interface{},
) ([]maintenance.Maintenance, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	maintenanceRecords := make([]maintenance.Maintenance, 0)
	for rows.Next() {
		var id, tenantID, vehicleID uuid.UUID
		var serviceType, odometer int
		var cost float64
		var serviceDate, createdAt, updatedAt time.Time
		var serviceProvider, description sql.NullString
		var nextServiceDue sql.NullTime
		var nextServiceOdometer sql.NullInt32

		if err := rows.Scan(
			&id,
			&tenantID,
			&vehicleID,
			&serviceType,
			&serviceDate,
			&odometer,
			&cost,
			&serviceProvider,
			&description,
			&nextServiceDue,
			&nextServiceOdometer,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		serviceProviderStr := ""
		if serviceProvider.Valid {
			serviceProviderStr = serviceProvider.String
		}

		descriptionStr := ""
		if description.Valid {
			descriptionStr = description.String
		}

		opts := []maintenance.MaintenanceOption{
			maintenance.WithTimestamps(createdAt, updatedAt),
		}

		if nextServiceDue.Valid {
			opts = append(opts, maintenance.WithNextServiceDue(nextServiceDue.Time))
		}

		if nextServiceOdometer.Valid {
			odo := int(nextServiceOdometer.Int32)
			opts = append(opts, maintenance.WithNextServiceOdometer(odo))
		}

		m := maintenance.NewMaintenance(
			id,
			tenantID,
			vehicleID,
			enums.ServiceType(serviceType),
			serviceDate,
			odometer,
			cost,
			serviceProviderStr,
			descriptionStr,
			opts...,
		)

		maintenanceRecords = append(maintenanceRecords, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return maintenanceRecords, nil
}

func (r *MaintenanceRepository) GetByID(ctx context.Context, id uuid.UUID) (maintenance.Maintenance, error) {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Join(selectMaintenanceQuery, "WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL")
	records, err := r.queryMaintenance(ctx, query, id, tenantID)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, ErrMaintenanceNotFound
	}
	return records[0], nil
}

func (r *MaintenanceRepository) GetPaginated(ctx context.Context, params *maintenance.FindParams) ([]maintenance.Maintenance, error) {
	where := make([]string, 0)
	args := make([]interface{}, 0)

	where = append(where, fmt.Sprintf("tenant_id = $%d", len(args)+1))
	args = append(args, params.TenantID)

	if !params.IncludeDeleted {
		where = append(where, "deleted_at IS NULL")
	}

	if params.VehicleID != nil {
		where = append(where, fmt.Sprintf("vehicle_id = $%d", len(args)+1))
		args = append(args, *params.VehicleID)
	}

	if params.ServiceType != nil {
		where = append(where, fmt.Sprintf("service_type = $%d", len(args)+1))
		args = append(args, int(*params.ServiceType))
	}

	if params.ServiceDateFrom != nil {
		where = append(where, fmt.Sprintf("service_date >= $%d", len(args)+1))
		args = append(args, *params.ServiceDateFrom)
	}

	if params.ServiceDateTo != nil {
		where = append(where, fmt.Sprintf("service_date <= $%d", len(args)+1))
		args = append(args, *params.ServiceDateTo)
	}

	if params.Search != nil && *params.Search != "" {
		searchPlaceholder := fmt.Sprintf("$%d", len(args)+1)
		where = append(where, fmt.Sprintf("(service_provider ILIKE %s OR description ILIKE %s)", searchPlaceholder, searchPlaceholder))
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
		selectMaintenanceQuery,
		repo.JoinWhere(where...),
		fmt.Sprintf("ORDER BY %s %s", sortColumn, sortDir),
		repo.FormatLimitOffset(params.Limit, params.Offset),
	)

	return r.queryMaintenance(ctx, query, args...)
}

func (r *MaintenanceRepository) Count(ctx context.Context, params *maintenance.FindParams) (int64, error) {
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

	if params.VehicleID != nil {
		where = append(where, fmt.Sprintf("vehicle_id = $%d", len(args)+1))
		args = append(args, *params.VehicleID)
	}

	if params.ServiceType != nil {
		where = append(where, fmt.Sprintf("service_type = $%d", len(args)+1))
		args = append(args, int(*params.ServiceType))
	}

	if params.ServiceDateFrom != nil {
		where = append(where, fmt.Sprintf("service_date >= $%d", len(args)+1))
		args = append(args, *params.ServiceDateFrom)
	}

	if params.ServiceDateTo != nil {
		where = append(where, fmt.Sprintf("service_date <= $%d", len(args)+1))
		args = append(args, *params.ServiceDateTo)
	}

	if params.Search != nil && *params.Search != "" {
		searchPlaceholder := fmt.Sprintf("$%d", len(args)+1)
		where = append(where, fmt.Sprintf("(service_provider ILIKE %s OR description ILIKE %s)", searchPlaceholder, searchPlaceholder))
		args = append(args, "%"+*params.Search+"%")
	}

	query := repo.Join(countMaintenanceQuery, repo.JoinWhere(where...))

	var count int64
	if err := pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *MaintenanceRepository) GetByVehicle(ctx context.Context, vehicleID uuid.UUID) ([]maintenance.Maintenance, error) {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Join(
		selectMaintenanceQuery,
		"WHERE vehicle_id = $1 AND tenant_id = $2 AND deleted_at IS NULL",
		"ORDER BY service_date DESC",
	)
	return r.queryMaintenance(ctx, query, vehicleID, tenantID)
}

func (r *MaintenanceRepository) GetDueMaintenance(ctx context.Context, tenantID uuid.UUID) ([]maintenance.Maintenance, error) {
	query := repo.Join(
		selectMaintenanceQuery,
		"WHERE tenant_id = $1 AND next_service_due IS NOT NULL AND next_service_due <= $2 AND deleted_at IS NULL",
		"ORDER BY next_service_due ASC",
	)
	return r.queryMaintenance(ctx, query, tenantID, time.Now())
}

func (r *MaintenanceRepository) Create(ctx context.Context, m maintenance.Maintenance) (maintenance.Maintenance, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	var serviceProvider, description interface{}
	if m.ServiceProvider() != "" {
		serviceProvider = m.ServiceProvider()
	}
	if m.Description() != "" {
		description = m.Description()
	}

	var nextServiceDue interface{}
	if m.NextServiceDue() != nil {
		nextServiceDue = *m.NextServiceDue()
	}

	var nextServiceOdometer interface{}
	if m.NextServiceOdometer() != nil {
		nextServiceOdometer = *m.NextServiceOdometer()
	}

	query := repo.Insert(
		"fleet_maintenance",
		[]string{
			"id",
			"tenant_id",
			"vehicle_id",
			"service_type",
			"service_date",
			"odometer",
			"cost",
			"service_provider",
			"description",
			"next_service_due",
			"next_service_odometer",
			"created_at",
			"updated_at",
		},
		"",
	)

	_, err = pool.Exec(
		ctx,
		query,
		m.ID(),
		m.TenantID(),
		m.VehicleID(),
		int(m.ServiceType()),
		m.ServiceDate(),
		m.Odometer(),
		m.Cost(),
		serviceProvider,
		description,
		nextServiceDue,
		nextServiceOdometer,
		m.CreatedAt(),
		m.UpdatedAt(),
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, m.ID())
}

func (r *MaintenanceRepository) Update(ctx context.Context, m maintenance.Maintenance) (maintenance.Maintenance, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	var serviceProvider, description interface{}
	if m.ServiceProvider() != "" {
		serviceProvider = m.ServiceProvider()
	}
	if m.Description() != "" {
		description = m.Description()
	}

	var nextServiceDue interface{}
	if m.NextServiceDue() != nil {
		nextServiceDue = *m.NextServiceDue()
	}

	var nextServiceOdometer interface{}
	if m.NextServiceOdometer() != nil {
		nextServiceOdometer = *m.NextServiceOdometer()
	}

	query := repo.Update(
		"fleet_maintenance",
		[]string{
			"vehicle_id",
			"service_type",
			"service_date",
			"odometer",
			"cost",
			"service_provider",
			"description",
			"next_service_due",
			"next_service_odometer",
			"updated_at",
		},
		fmt.Sprintf("id = $%d AND tenant_id = $%d", 11, 12),
	)

	_, err = pool.Exec(
		ctx,
		query,
		m.VehicleID(),
		int(m.ServiceType()),
		m.ServiceDate(),
		m.Odometer(),
		m.Cost(),
		serviceProvider,
		description,
		nextServiceDue,
		nextServiceOdometer,
		time.Now(),
		m.ID(),
		m.TenantID(),
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, m.ID())
}

func (r *MaintenanceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return err
	}

	pool, err := composables.UseTx(ctx)
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, deleteMaintenanceQuery, time.Now(), id, tenantID)
	return err
}
