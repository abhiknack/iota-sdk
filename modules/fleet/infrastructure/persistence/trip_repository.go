package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/trip"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/repo"
)

var (
	ErrTripNotFound = errors.New("trip not found")
)

const (
	selectTripQuery = `
		SELECT
			id,
			tenant_id,
			vehicle_id,
			driver_id,
			origin,
			destination,
			purpose,
			start_time,
			end_time,
			start_odometer,
			end_odometer,
			status,
			created_at,
			updated_at
		FROM fleet_trips
	`
	countTripQuery  = `SELECT COUNT(*) FROM fleet_trips`
	deleteTripQuery = `UPDATE fleet_trips SET deleted_at = $1 WHERE id = $2 AND tenant_id = $3`
)

type TripRepository struct {
	fieldMap map[trip.Field]string
}

func NewTripRepository() trip.Repository {
	return &TripRepository{
		fieldMap: map[trip.Field]string{
			trip.FieldID:            "id",
			trip.FieldTenantID:      "tenant_id",
			trip.FieldVehicleID:     "vehicle_id",
			trip.FieldDriverID:      "driver_id",
			trip.FieldOrigin:        "origin",
			trip.FieldDestination:   "destination",
			trip.FieldPurpose:       "purpose",
			trip.FieldStartTime:     "start_time",
			trip.FieldEndTime:       "end_time",
			trip.FieldStartOdometer: "start_odometer",
			trip.FieldEndOdometer:   "end_odometer",
			trip.FieldStatus:        "status",
			trip.FieldCreatedAt:     "created_at",
			trip.FieldUpdatedAt:     "updated_at",
		},
	}
}

func (r *TripRepository) queryTrips(
	ctx context.Context,
	query string,
	args ...interface{},
) ([]trip.Trip, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trips := make([]trip.Trip, 0)
	for rows.Next() {
		var id, tenantID, vehicleID, driverID uuid.UUID
		var origin, destination string
		var purpose sql.NullString
		var startTime, createdAt, updatedAt time.Time
		var endTime sql.NullTime
		var startOdometer, status int
		var endOdometer sql.NullInt32

		if err := rows.Scan(
			&id,
			&tenantID,
			&vehicleID,
			&driverID,
			&origin,
			&destination,
			&purpose,
			&startTime,
			&endTime,
			&startOdometer,
			&endOdometer,
			&status,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		purposeStr := ""
		if purpose.Valid {
			purposeStr = purpose.String
		}

		opts := []trip.TripOption{
			trip.WithStatus(enums.TripStatus(status)),
			trip.WithTimestamps(createdAt, updatedAt),
		}

		if endTime.Valid {
			opts = append(opts, trip.WithEndTime(endTime.Time))
		}

		if endOdometer.Valid {
			opts = append(opts, trip.WithEndOdometer(int(endOdometer.Int32)))
		}

		t := trip.NewTrip(
			id,
			tenantID,
			vehicleID,
			driverID,
			origin,
			destination,
			purposeStr,
			startTime,
			startOdometer,
			opts...,
		)

		trips = append(trips, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trips, nil
}

func (r *TripRepository) GetByID(ctx context.Context, id uuid.UUID) (trip.Trip, error) {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Join(selectTripQuery, "WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL")
	trips, err := r.queryTrips(ctx, query, id, tenantID)
	if err != nil {
		return nil, err
	}
	if len(trips) == 0 {
		return nil, ErrTripNotFound
	}
	return trips[0], nil
}

func (r *TripRepository) GetPaginated(ctx context.Context, params *trip.FindParams) ([]trip.Trip, error) {
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

	if params.Status != nil {
		where = append(where, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, int(*params.Status))
	}

	if params.StartTimeFrom != nil {
		where = append(where, fmt.Sprintf("start_time >= $%d", len(args)+1))
		args = append(args, *params.StartTimeFrom)
	}

	if params.StartTimeTo != nil {
		where = append(where, fmt.Sprintf("start_time <= $%d", len(args)+1))
		args = append(args, *params.StartTimeTo)
	}

	if params.Search != nil && *params.Search != "" {
		searchPlaceholder := fmt.Sprintf("$%d", len(args)+1)
		where = append(where, fmt.Sprintf("(origin ILIKE %s OR destination ILIKE %s OR purpose ILIKE %s)", searchPlaceholder, searchPlaceholder, searchPlaceholder))
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
		selectTripQuery,
		repo.JoinWhere(where...),
		fmt.Sprintf("ORDER BY %s %s", sortColumn, sortDir),
		repo.FormatLimitOffset(params.Limit, params.Offset),
	)

	return r.queryTrips(ctx, query, args...)
}

func (r *TripRepository) Count(ctx context.Context, params *trip.FindParams) (int64, error) {
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

	if params.Status != nil {
		where = append(where, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, int(*params.Status))
	}

	if params.StartTimeFrom != nil {
		where = append(where, fmt.Sprintf("start_time >= $%d", len(args)+1))
		args = append(args, *params.StartTimeFrom)
	}

	if params.StartTimeTo != nil {
		where = append(where, fmt.Sprintf("start_time <= $%d", len(args)+1))
		args = append(args, *params.StartTimeTo)
	}

	if params.Search != nil && *params.Search != "" {
		searchPlaceholder := fmt.Sprintf("$%d", len(args)+1)
		where = append(where, fmt.Sprintf("(origin ILIKE %s OR destination ILIKE %s OR purpose ILIKE %s)", searchPlaceholder, searchPlaceholder, searchPlaceholder))
		args = append(args, "%"+*params.Search+"%")
	}

	query := repo.Join(countTripQuery, repo.JoinWhere(where...))

	var count int64
	if err := pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *TripRepository) GetByVehicle(ctx context.Context, vehicleID uuid.UUID) ([]trip.Trip, error) {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Join(
		selectTripQuery,
		"WHERE vehicle_id = $1 AND tenant_id = $2 AND deleted_at IS NULL",
		"ORDER BY start_time DESC",
	)
	return r.queryTrips(ctx, query, vehicleID, tenantID)
}

func (r *TripRepository) GetByDriver(ctx context.Context, driverID uuid.UUID) ([]trip.Trip, error) {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Join(
		selectTripQuery,
		"WHERE driver_id = $1 AND tenant_id = $2 AND deleted_at IS NULL",
		"ORDER BY start_time DESC",
	)
	return r.queryTrips(ctx, query, driverID, tenantID)
}

func (r *TripRepository) GetActiveTrips(ctx context.Context, tenantID uuid.UUID) ([]trip.Trip, error) {
	query := repo.Join(
		selectTripQuery,
		"WHERE tenant_id = $1 AND status IN ($2, $3) AND deleted_at IS NULL",
		"ORDER BY start_time ASC",
	)
	return r.queryTrips(ctx, query, tenantID, int(enums.TripStatusScheduled), int(enums.TripStatusInProgress))
}

func (r *TripRepository) CheckConflict(ctx context.Context, vehicleID uuid.UUID, startTime, endTime time.Time, excludeTripID *uuid.UUID) (bool, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return false, err
	}

	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return false, err
	}

	query := `
		SELECT EXISTS(
			SELECT 1 FROM fleet_trips
			WHERE vehicle_id = $1
			AND tenant_id = $2
			AND deleted_at IS NULL
			AND status IN ($3, $4)
			AND (
				(start_time <= $5 AND (end_time IS NULL OR end_time >= $5))
				OR (start_time <= $6 AND (end_time IS NULL OR end_time >= $6))
				OR (start_time >= $5 AND start_time <= $6)
			)
	`

	args := []interface{}{
		vehicleID,
		tenantID,
		int(enums.TripStatusScheduled),
		int(enums.TripStatusInProgress),
		startTime,
		endTime,
	}

	if excludeTripID != nil {
		query += " AND id != $7"
		args = append(args, *excludeTripID)
	}

	query += ")"

	var exists bool
	if err := pool.QueryRow(ctx, query, args...).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *TripRepository) Create(ctx context.Context, t trip.Trip) (trip.Trip, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	var purpose interface{}
	if t.Purpose() != "" {
		purpose = t.Purpose()
	}

	var endTime interface{}
	if t.EndTime() != nil {
		endTime = *t.EndTime()
	}

	var endOdometer interface{}
	if t.EndOdometer() != nil {
		endOdometer = *t.EndOdometer()
	}

	query := repo.Insert(
		"fleet_trips",
		[]string{
			"id",
			"tenant_id",
			"vehicle_id",
			"driver_id",
			"origin",
			"destination",
			"purpose",
			"start_time",
			"end_time",
			"start_odometer",
			"end_odometer",
			"status",
			"created_at",
			"updated_at",
		},
		"",
	)

	_, err = pool.Exec(
		ctx,
		query,
		t.ID(),
		t.TenantID(),
		t.VehicleID(),
		t.DriverID(),
		t.Origin(),
		t.Destination(),
		purpose,
		t.StartTime(),
		endTime,
		t.StartOdometer(),
		endOdometer,
		int(t.Status()),
		t.CreatedAt(),
		t.UpdatedAt(),
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, t.ID())
}

func (r *TripRepository) Update(ctx context.Context, t trip.Trip) (trip.Trip, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	var purpose interface{}
	if t.Purpose() != "" {
		purpose = t.Purpose()
	}

	var endTime interface{}
	if t.EndTime() != nil {
		endTime = *t.EndTime()
	}

	var endOdometer interface{}
	if t.EndOdometer() != nil {
		endOdometer = *t.EndOdometer()
	}

	query := repo.Update(
		"fleet_trips",
		[]string{
			"vehicle_id",
			"driver_id",
			"origin",
			"destination",
			"purpose",
			"start_time",
			"end_time",
			"start_odometer",
			"end_odometer",
			"status",
			"updated_at",
		},
		fmt.Sprintf("id = $%d AND tenant_id = $%d", 12, 13),
	)

	_, err = pool.Exec(
		ctx,
		query,
		t.VehicleID(),
		t.DriverID(),
		t.Origin(),
		t.Destination(),
		purpose,
		t.StartTime(),
		endTime,
		t.StartOdometer(),
		endOdometer,
		int(t.Status()),
		time.Now(),
		t.ID(),
		t.TenantID(),
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, t.ID())
}

func (r *TripRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return err
	}

	pool, err := composables.UseTx(ctx)
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, deleteTripQuery, time.Now(), id, tenantID)
	return err
}
