package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/vehicle"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/repo"
)

var (
	ErrVehicleNotFound = errors.New("vehicle not found")
)

const (
	selectVehicleQuery = `
		SELECT
			id,
			tenant_id,
			make,
			model,
			year,
			vin,
			license_plate,
			status,
			current_odometer,
			registration_expiry,
			insurance_expiry,
			created_at,
			updated_at
		FROM fleet_vehicles
	`
	countVehicleQuery  = `SELECT COUNT(*) FROM fleet_vehicles`
	deleteVehicleQuery = `UPDATE fleet_vehicles SET deleted_at = $1 WHERE id = $2 AND tenant_id = $3`
)

type VehicleRepository struct {
	fieldMap map[vehicle.Field]string
}

func NewVehicleRepository() vehicle.Repository {
	return &VehicleRepository{
		fieldMap: map[vehicle.Field]string{
			vehicle.FieldID:                 "id",
			vehicle.FieldTenantID:           "tenant_id",
			vehicle.FieldMake:               "make",
			vehicle.FieldModel:              "model",
			vehicle.FieldYear:               "year",
			vehicle.FieldVIN:                "vin",
			vehicle.FieldLicensePlate:       "license_plate",
			vehicle.FieldStatus:             "status",
			vehicle.FieldCurrentOdometer:    "current_odometer",
			vehicle.FieldRegistrationExpiry: "registration_expiry",
			vehicle.FieldInsuranceExpiry:    "insurance_expiry",
			vehicle.FieldCreatedAt:          "created_at",
			vehicle.FieldUpdatedAt:          "updated_at",
		},
	}
}

func (r *VehicleRepository) queryVehicles(
	ctx context.Context,
	query string,
	args ...interface{},
) ([]vehicle.Vehicle, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vehicles := make([]vehicle.Vehicle, 0)
	for rows.Next() {
		var id, tenantID uuid.UUID
		var make, model, vin, licensePlate string
		var year, status, currentOdometer int
		var registrationExpiry, insuranceExpiry, createdAt, updatedAt time.Time

		if err := rows.Scan(
			&id,
			&tenantID,
			&make,
			&model,
			&year,
			&vin,
			&licensePlate,
			&status,
			&currentOdometer,
			&registrationExpiry,
			&insuranceExpiry,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		v := vehicle.NewVehicle(
			id,
			tenantID,
			make,
			model,
			year,
			vin,
			licensePlate,
			vehicle.WithStatus(enums.VehicleStatus(status)),
			vehicle.WithOdometer(currentOdometer),
			vehicle.WithRegistrationExpiry(registrationExpiry),
			vehicle.WithInsuranceExpiry(insuranceExpiry),
			vehicle.WithTimestamps(createdAt, updatedAt),
		)

		vehicles = append(vehicles, v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return vehicles, nil
}

func (r *VehicleRepository) GetByID(ctx context.Context, id uuid.UUID) (vehicle.Vehicle, error) {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Join(selectVehicleQuery, "WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL")
	vehicles, err := r.queryVehicles(ctx, query, id, tenantID)
	if err != nil {
		return nil, err
	}
	if len(vehicles) == 0 {
		return nil, ErrVehicleNotFound
	}
	return vehicles[0], nil
}

func (r *VehicleRepository) GetPaginated(ctx context.Context, params *vehicle.FindParams) ([]vehicle.Vehicle, error) {
	where := make([]string, 0)
	args := make([]interface{}, 0)

	where = append(where, fmt.Sprintf("tenant_id = $%d", len(args)+1))
	args = append(args, params.TenantID)

	if !params.IncludeDeleted {
		where = append(where, "deleted_at IS NULL")
	}

	if params.Status != nil {
		where = append(where, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, int(*params.Status))
	}

	if params.Make != nil {
		where = append(where, fmt.Sprintf("make ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Make+"%")
	}

	if params.Model != nil {
		where = append(where, fmt.Sprintf("model ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Model+"%")
	}

	if params.Year != nil {
		where = append(where, fmt.Sprintf("year = $%d", len(args)+1))
		args = append(args, *params.Year)
	}

	if params.Search != nil && *params.Search != "" {
		searchPlaceholder := fmt.Sprintf("$%d", len(args)+1)
		where = append(where, fmt.Sprintf("(make ILIKE %s OR model ILIKE %s OR vin ILIKE %s OR license_plate ILIKE %s)", searchPlaceholder, searchPlaceholder, searchPlaceholder, searchPlaceholder))
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
		selectVehicleQuery,
		repo.JoinWhere(where...),
		fmt.Sprintf("ORDER BY %s %s", sortColumn, sortDir),
		repo.FormatLimitOffset(params.Limit, params.Offset),
	)

	return r.queryVehicles(ctx, query, args...)
}

func (r *VehicleRepository) Count(ctx context.Context, params *vehicle.FindParams) (int64, error) {
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

	if params.Status != nil {
		where = append(where, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, int(*params.Status))
	}

	if params.Make != nil {
		where = append(where, fmt.Sprintf("make ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Make+"%")
	}

	if params.Model != nil {
		where = append(where, fmt.Sprintf("model ILIKE $%d", len(args)+1))
		args = append(args, "%"+*params.Model+"%")
	}

	if params.Year != nil {
		where = append(where, fmt.Sprintf("year = $%d", len(args)+1))
		args = append(args, *params.Year)
	}

	if params.Search != nil && *params.Search != "" {
		searchPlaceholder := fmt.Sprintf("$%d", len(args)+1)
		where = append(where, fmt.Sprintf("(make ILIKE %s OR model ILIKE %s OR vin ILIKE %s OR license_plate ILIKE %s)", searchPlaceholder, searchPlaceholder, searchPlaceholder, searchPlaceholder))
		args = append(args, "%"+*params.Search+"%")
	}

	query := repo.Join(countVehicleQuery, repo.JoinWhere(where...))

	var count int64
	if err := pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *VehicleRepository) GetByStatus(ctx context.Context, tenantID uuid.UUID, status enums.VehicleStatus) ([]vehicle.Vehicle, error) {
	query := repo.Join(selectVehicleQuery, "WHERE tenant_id = $1 AND status = $2 AND deleted_at IS NULL")
	return r.queryVehicles(ctx, query, tenantID, int(status))
}

func (r *VehicleRepository) GetExpiringRegistrations(ctx context.Context, tenantID uuid.UUID, days int) ([]vehicle.Vehicle, error) {
	query := repo.Join(
		selectVehicleQuery,
		"WHERE tenant_id = $1 AND registration_expiry <= $2 AND registration_expiry >= $3 AND deleted_at IS NULL",
		"ORDER BY registration_expiry ASC",
	)
	expiryDate := time.Now().AddDate(0, 0, days)
	today := time.Now()
	return r.queryVehicles(ctx, query, tenantID, expiryDate, today)
}

func (r *VehicleRepository) Create(ctx context.Context, v vehicle.Vehicle) (vehicle.Vehicle, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Insert(
		"fleet_vehicles",
		[]string{
			"id",
			"tenant_id",
			"make",
			"model",
			"year",
			"vin",
			"license_plate",
			"status",
			"current_odometer",
			"registration_expiry",
			"insurance_expiry",
			"created_at",
			"updated_at",
		},
		"",
	)

	_, err = pool.Exec(
		ctx,
		query,
		v.ID(),
		v.TenantID(),
		v.Make(),
		v.Model(),
		v.Year(),
		v.VIN(),
		v.LicensePlate(),
		int(v.Status()),
		v.CurrentOdometer(),
		v.RegistrationExpiry(),
		v.InsuranceExpiry(),
		v.CreatedAt(),
		v.UpdatedAt(),
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, v.ID())
}

func (r *VehicleRepository) Update(ctx context.Context, v vehicle.Vehicle) (vehicle.Vehicle, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Update(
		"fleet_vehicles",
		[]string{
			"make",
			"model",
			"year",
			"vin",
			"license_plate",
			"status",
			"current_odometer",
			"registration_expiry",
			"insurance_expiry",
			"updated_at",
		},
		fmt.Sprintf("id = $%d AND tenant_id = $%d", 11, 12),
	)

	_, err = pool.Exec(
		ctx,
		query,
		v.Make(),
		v.Model(),
		v.Year(),
		v.VIN(),
		v.LicensePlate(),
		int(v.Status()),
		v.CurrentOdometer(),
		v.RegistrationExpiry(),
		v.InsuranceExpiry(),
		time.Now(),
		v.ID(),
		v.TenantID(),
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, v.ID())
}

func (r *VehicleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return err
	}

	pool, err := composables.UseTx(ctx)
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, deleteVehicleQuery, time.Now(), id, tenantID)
	return err
}
