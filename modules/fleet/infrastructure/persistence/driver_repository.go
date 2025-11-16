package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/driver"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/repo"
)

var (
	ErrDriverNotFound = errors.New("driver not found")
)

const (
	selectDriverQuery = `
		SELECT
			id,
			tenant_id,
			user_id,
			first_name,
			last_name,
			license_number,
			license_expiry,
			phone,
			email,
			status,
			created_at,
			updated_at
		FROM fleet_drivers
	`
	countDriverQuery  = `SELECT COUNT(*) FROM fleet_drivers`
	deleteDriverQuery = `UPDATE fleet_drivers SET deleted_at = $1 WHERE id = $2 AND tenant_id = $3`
)

type DriverRepository struct {
	fieldMap map[driver.Field]string
}

func NewDriverRepository() driver.Repository {
	return &DriverRepository{
		fieldMap: map[driver.Field]string{
			driver.FieldID:            "id",
			driver.FieldTenantID:      "tenant_id",
			driver.FieldUserID:        "user_id",
			driver.FieldFirstName:     "first_name",
			driver.FieldLastName:      "last_name",
			driver.FieldLicenseNumber: "license_number",
			driver.FieldLicenseExpiry: "license_expiry",
			driver.FieldPhone:         "phone",
			driver.FieldEmail:         "email",
			driver.FieldStatus:        "status",
			driver.FieldCreatedAt:     "created_at",
			driver.FieldUpdatedAt:     "updated_at",
		},
	}
}

func (r *DriverRepository) queryDrivers(
	ctx context.Context,
	query string,
	args ...interface{},
) ([]driver.Driver, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	drivers := make([]driver.Driver, 0)
	for rows.Next() {
		var id, tenantID uuid.UUID
		var userID sql.NullInt64
		var firstName, lastName, licenseNumber string
		var phone, email sql.NullString
		var licenseExpiry, createdAt, updatedAt time.Time
		var status int

		if err := rows.Scan(
			&id,
			&tenantID,
			&userID,
			&firstName,
			&lastName,
			&licenseNumber,
			&licenseExpiry,
			&phone,
			&email,
			&status,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		opts := []driver.DriverOption{
			driver.WithDriverStatus(enums.DriverStatus(status)),
			driver.WithDriverTimestamps(createdAt, updatedAt),
		}

		if userID.Valid {
			opts = append(opts, driver.WithUserID(userID.Int64))
		}

		if phone.Valid {
			opts = append(opts, driver.WithPhone(phone.String))
		}

		if email.Valid {
			opts = append(opts, driver.WithEmail(email.String))
		}

		d := driver.NewDriver(
			id,
			tenantID,
			firstName,
			lastName,
			licenseNumber,
			licenseExpiry,
			opts...,
		)

		drivers = append(drivers, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return drivers, nil
}

func (r *DriverRepository) GetByID(ctx context.Context, id uuid.UUID) (driver.Driver, error) {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Join(selectDriverQuery, "WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL")
	drivers, err := r.queryDrivers(ctx, query, id, tenantID)
	if err != nil {
		return nil, err
	}
	if len(drivers) == 0 {
		return nil, ErrDriverNotFound
	}
	return drivers[0], nil
}

func (r *DriverRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (driver.Driver, error) {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return nil, err
	}

	query := repo.Join(selectDriverQuery, "WHERE user_id = $1 AND tenant_id = $2 AND deleted_at IS NULL")
	drivers, err := r.queryDrivers(ctx, query, userID, tenantID)
	if err != nil {
		return nil, err
	}
	if len(drivers) == 0 {
		return nil, ErrDriverNotFound
	}
	return drivers[0], nil
}

func (r *DriverRepository) GetPaginated(ctx context.Context, params *driver.FindParams) ([]driver.Driver, error) {
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

	if params.Search != nil && *params.Search != "" {
		searchPlaceholder := fmt.Sprintf("$%d", len(args)+1)
		where = append(where, fmt.Sprintf("(first_name ILIKE %s OR last_name ILIKE %s OR license_number ILIKE %s OR phone ILIKE %s OR email ILIKE %s)", searchPlaceholder, searchPlaceholder, searchPlaceholder, searchPlaceholder, searchPlaceholder))
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
		selectDriverQuery,
		repo.JoinWhere(where...),
		fmt.Sprintf("ORDER BY %s %s", sortColumn, sortDir),
		repo.FormatLimitOffset(params.Limit, params.Offset),
	)

	return r.queryDrivers(ctx, query, args...)
}

func (r *DriverRepository) Count(ctx context.Context, params *driver.FindParams) (int64, error) {
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

	if params.Search != nil && *params.Search != "" {
		searchPlaceholder := fmt.Sprintf("$%d", len(args)+1)
		where = append(where, fmt.Sprintf("(first_name ILIKE %s OR last_name ILIKE %s OR license_number ILIKE %s OR phone ILIKE %s OR email ILIKE %s)", searchPlaceholder, searchPlaceholder, searchPlaceholder, searchPlaceholder, searchPlaceholder))
		args = append(args, "%"+*params.Search+"%")
	}

	query := repo.Join(countDriverQuery, repo.JoinWhere(where...))

	var count int64
	if err := pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *DriverRepository) GetExpiringLicenses(ctx context.Context, tenantID uuid.UUID, days int) ([]driver.Driver, error) {
	query := repo.Join(
		selectDriverQuery,
		"WHERE tenant_id = $1 AND license_expiry <= $2 AND license_expiry >= $3 AND deleted_at IS NULL",
		"ORDER BY license_expiry ASC",
	)
	expiryDate := time.Now().AddDate(0, 0, days)
	today := time.Now()
	return r.queryDrivers(ctx, query, tenantID, expiryDate, today)
}

func (r *DriverRepository) GetAvailable(ctx context.Context, tenantID uuid.UUID, startTime, endTime time.Time) ([]driver.Driver, error) {
	query := `
		SELECT DISTINCT
			d.id,
			d.tenant_id,
			d.user_id,
			d.first_name,
			d.last_name,
			d.license_number,
			d.license_expiry,
			d.phone,
			d.email,
			d.status,
			d.created_at,
			d.updated_at
		FROM fleet_drivers d
		WHERE d.tenant_id = $1
		AND d.status = $2
		AND d.deleted_at IS NULL
		AND NOT EXISTS (
			SELECT 1 FROM fleet_trips t
			WHERE t.driver_id = d.id
			AND t.deleted_at IS NULL
			AND t.status IN (0, 1)
			AND (
				(t.start_time <= $3 AND (t.end_time IS NULL OR t.end_time >= $3))
				OR (t.start_time <= $4 AND (t.end_time IS NULL OR t.end_time >= $4))
				OR (t.start_time >= $3 AND t.start_time <= $4)
			)
		)
		ORDER BY d.first_name, d.last_name
	`
	return r.queryDrivers(ctx, query, tenantID, int(enums.DriverStatusActive), startTime, endTime)
}

func (r *DriverRepository) Create(ctx context.Context, d driver.Driver) (driver.Driver, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	var userID interface{}
	if d.UserID() != nil {
		userID = *d.UserID()
	}

	var phone, email interface{}
	if d.Phone() != "" {
		phone = d.Phone()
	}
	if d.Email() != "" {
		email = d.Email()
	}

	query := repo.Insert(
		"fleet_drivers",
		[]string{
			"id",
			"tenant_id",
			"user_id",
			"first_name",
			"last_name",
			"license_number",
			"license_expiry",
			"phone",
			"email",
			"status",
			"created_at",
			"updated_at",
		},
		"",
	)

	_, err = pool.Exec(
		ctx,
		query,
		d.ID(),
		d.TenantID(),
		userID,
		d.FirstName(),
		d.LastName(),
		d.LicenseNumber(),
		d.LicenseExpiry(),
		phone,
		email,
		int(d.Status()),
		d.CreatedAt(),
		d.UpdatedAt(),
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, d.ID())
}

func (r *DriverRepository) Update(ctx context.Context, d driver.Driver) (driver.Driver, error) {
	pool, err := composables.UseTx(ctx)
	if err != nil {
		return nil, err
	}

	var userID interface{}
	if d.UserID() != nil {
		userID = *d.UserID()
	}

	var phone, email interface{}
	if d.Phone() != "" {
		phone = d.Phone()
	}
	if d.Email() != "" {
		email = d.Email()
	}

	query := repo.Update(
		"fleet_drivers",
		[]string{
			"user_id",
			"first_name",
			"last_name",
			"license_number",
			"license_expiry",
			"phone",
			"email",
			"status",
			"updated_at",
		},
		fmt.Sprintf("id = $%d AND tenant_id = $%d", 10, 11),
	)

	_, err = pool.Exec(
		ctx,
		query,
		userID,
		d.FirstName(),
		d.LastName(),
		d.LicenseNumber(),
		d.LicenseExpiry(),
		phone,
		email,
		int(d.Status()),
		time.Now(),
		d.ID(),
		d.TenantID(),
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, d.ID())
}

func (r *DriverRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		return err
	}

	pool, err := composables.UseTx(ctx)
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, deleteDriverQuery, time.Now(), id, tenantID)
	return err
}
