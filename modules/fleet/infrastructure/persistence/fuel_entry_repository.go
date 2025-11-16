package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/fuel_entry"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/repo"
)

var (
	ErrFuelEntryNotFound = errors.New("fuel entry not found")
)

const (
	selectFuelEntryQuery = `
		SELECT
			id,
			tenant_id,
			vehicle_id,
			driver_id,
			date,
			quantity,
			cost,
			odometer,
			fuel_type,
			location,
			created_at,
			updated_at
		FROM fleet_fuel_entries
	`
	countFuelEntryQuery  = `SELECT COUNT(*) FROM fleet_fuel_entries`
	deleteFuelEntryQuery = `UPDATE fleet_fuel_entries SET deleted_at = $1 WHERE id = $2 AND tenant_id = $3`
)

type FuelEntryRepository struct {
	fieldMap map[fuel_entry.Field]string
}

func NewFuelEntryRepository() fuel_entry.Repository {
	return &FuelEntryRepository{
		fieldMap: map[fuel_entry.Field]string{
			fuel_entry.FieldID:        "id",
			fuel_entry.FieldTenantID:  "tenant_id",
			fuel_entry.FieldVehicleID: "vehicle_id",
			fuel_entry.FieldDriverID:  "driver_id",
			fuel_entry.FieldDate:      "date",
			fuel_entry.FieldQuantity:  "quantity",
			fuel_entry.FieldCost:      "cost",
			fuel_entry.FieldOdometer:  "odometer",
			fuel_entry.FieldFuelType:  "fuel_type",
			fuel_entry.FieldLocation:  "location",
			fuel_entry.FieldCreatedAt: "created_at",
			fuel_entry.FieldUpdatedAt: "updated_at",
		},
	}
}

func (r *FuelEntryRepository) queryFuelEntries(
	ctx context.Context,
	query string,
	args ...interface{},
) ([]fuel_entry.FuelEntry, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]fuel_entry.FuelEntry, 0)
	for rows.Next() {
		var id, tenantID, vehicleID uuid.UUID
		var driverID sql.NullString
		var date, createdAt, updatedAt time.Time
		var quantity, cost float64
		var odometer, fuelType int
		var location sql.NullString

		if err := rows.Scan(
			&id,
			&tenantID,
			&vehicleID,
			&driverID,
			&date,
			&quantity,
			&cost,
			&odometer,
			&fuelType,
			&location,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		locationStr := ""
		if location.Valid {
			locationStr = location.String
		}

		opts := []fuel_entry.FuelEntryOption{
			fuel_entry.WithTimestamps(createdAt, updatedAt),
		}

		if driverID.Valid {
			did, _ := uuid.Parse(driverID.String)
			opts = append(opts, fuel_entry.WithDriverID(did))
		}

		f := fuel_entry.NewFuelEntry(
			id,
			tenantID,
			vehicleID,
			date,
			quantity,
			cost,
			odometer,
			enums.FuelType(fuelType),
			locationStr,
			opts...,
		)

		entries = append(entries, f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *FuelEntryRepository) GetByID(ctx context.Context, id uuid.UUID) (fuel_entry.FuelEntry, error) {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Join(selectFuelEntryQuery, "WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL")
	entries, err := r.queryFuelEntries(ctx, query, id, tenantID)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, ErrFuelEntryNotFound
	}
	return entries[0], nil
}

func (r *FuelEntryRepository) GetPaginated(ctx context.Context, params *fuel_entry.FindParams) ([]fuel_entry.FuelEntry, error) {
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

	if params.DriverID != nil {
		where = append(where, fmt.Sprintf("driver_id = $%d", len(args)+1))
		args = append(args, *params.DriverID)
	}

	if params.FuelType != nil {
		where = append(where, fmt.Sprintf("fuel_type = $%d", len(args)+1))
		args = append(args, int(*params.FuelType))
	}

	if params.DateFrom != nil {
		where = append(where, fmt.Sprintf("date >= $%d", len(args)+1))
		args = append(args, *params.DateFrom)
	}

	if params.DateTo != nil {
		where = append(where, fmt.Sprintf("date <= $%d", len(args)+1))
		args = append(args, *params.DateTo)
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
		selectFuelEntryQuery,
		repo.JoinWhere(where...),
		fmt.Sprintf("ORDER BY %s %s", sortColumn, sortDir),
		repo.FormatLimitOffset(params.Limit, params.Offset),
	)

	return r.queryFuelEntries(ctx, query, args...)
}

func (r *FuelEntryRepository) Count(ctx context.Context, params *fuel_entry.FindParams) (int64, error) {
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

	if params.DriverID != nil {
		where = append(where, fmt.Sprintf("driver_id = $%d", len(args)+1))
		args = append(args, *params.DriverID)
	}

	if params.FuelType != nil {
		where = append(where, fmt.Sprintf("fuel_type = $%d", len(args)+1))
		args = append(args, int(*params.FuelType))
	}

	if params.DateFrom != nil {
		where = append(where, fmt.Sprintf("date >= $%d", len(args)+1))
		args = append(args, *params.DateFrom)
	}

	if params.DateTo != nil {
		where = append(where, fmt.Sprintf("date <= $%d", len(args)+1))
		args = append(args, *params.DateTo)
	}

	query := repo.Join(countFuelEntryQuery, repo.JoinWhere(where...))

	var count int64
	if err := pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *FuelEntryRepository) GetByVehicle(ctx context.Context, vehicleID uuid.UUID) ([]fuel_entry.FuelEntry, error) {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Join(
		selectFuelEntryQuery,
		"WHERE vehicle_id = $1 AND tenant_id = $2 AND deleted_at IS NULL",
		"ORDER BY date DESC",
	)
	return r.queryFuelEntries(ctx, query, vehicleID, tenantID)
}

func (r *FuelEntryRepository) GetByDriver(ctx context.Context, driverID uuid.UUID) ([]fuel_entry.FuelEntry, error) {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Join(
		selectFuelEntryQuery,
		"WHERE driver_id = $1 AND tenant_id = $2 AND deleted_at IS NULL",
		"ORDER BY date DESC",
	)
	return r.queryFuelEntries(ctx, query, driverID, tenantID)
}

func (r *FuelEntryRepository) GetLastEntry(ctx context.Context, vehicleID uuid.UUID) (fuel_entry.FuelEntry, error) {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Join(
		selectFuelEntryQuery,
		"WHERE vehicle_id = $1 AND tenant_id = $2 AND deleted_at IS NULL",
		"ORDER BY date DESC, odometer DESC",
		"LIMIT 1",
	)
	entries, err := r.queryFuelEntries(ctx, query, vehicleID, tenantID)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, ErrFuelEntryNotFound
	}
	return entries[0], nil
}

func (r *FuelEntryRepository) Create(ctx context.Context, f fuel_entry.FuelEntry) (fuel_entry.FuelEntry, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	var driverID interface{}
	if f.DriverID() != nil {
		driverID = *f.DriverID()
	}

	var location interface{}
	if f.Location() != "" {
		location = f.Location()
	}

	query := repo.Insert(
		"fleet_fuel_entries",
		[]string{
			"id",
			"tenant_id",
			"vehicle_id",
			"driver_id",
			"date",
			"quantity",
			"cost",
			"odometer",
			"fuel_type",
			"location",
			"created_at",
			"updated_at",
		},
		"",
	)

	_, err = pool.Exec(
		ctx,
		query,
		f.ID(),
		f.TenantID(),
		f.VehicleID(),
		driverID,
		f.Date(),
		f.Quantity(),
		f.Cost(),
		f.Odometer(),
		int(f.FuelType()),
		location,
		f.CreatedAt(),
		f.UpdatedAt(),
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, f.ID())
}

func (r *FuelEntryRepository) Update(ctx context.Context, f fuel_entry.FuelEntry) (fuel_entry.FuelEntry, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	var driverID interface{}
	if f.DriverID() != nil {
		driverID = *f.DriverID()
	}

	var location interface{}
	if f.Location() != "" {
		location = f.Location()
	}

	query := repo.Update(
		"fleet_fuel_entries",
		[]string{
			"vehicle_id",
			"driver_id",
			"date",
			"quantity",
			"cost",
			"odometer",
			"fuel_type",
			"location",
			"updated_at",
		},
		fmt.Sprintf("id = $%d AND tenant_id = $%d", 10, 11),
	)

	_, err = pool.Exec(
		ctx,
		query,
		f.VehicleID(),
		driverID,
		f.Date(),
		f.Quantity(),
		f.Cost(),
		f.Odometer(),
		int(f.FuelType()),
		location,
		time.Now(),
		f.ID(),
		f.TenantID(),
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, f.ID())
}

func (r *FuelEntryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return err
	}

	pool, err := composables.UseTx(ctx)
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, deleteFuelEntryQuery, time.Now(), id, tenantID)
	return err
}
