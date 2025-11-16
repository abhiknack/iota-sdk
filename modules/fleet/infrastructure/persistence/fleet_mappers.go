package persistence

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/driver"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/fuel_entry"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/maintenance"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/trip"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/vehicle"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
	"github.com/iota-uz/iota-sdk/modules/fleet/infrastructure/persistence/models"
)

func ToDomainVehicle(dbRow *models.Vehicle) vehicle.Vehicle {
	return vehicle.NewVehicle(
		dbRow.ID,
		dbRow.TenantID,
		dbRow.Make,
		dbRow.Model,
		dbRow.Year,
		dbRow.VIN,
		dbRow.LicensePlate,
		vehicle.WithStatus(enums.VehicleStatus(dbRow.Status)),
		vehicle.WithOdometer(dbRow.CurrentOdometer),
		vehicle.WithRegistrationExpiry(dbRow.RegistrationExpiry),
		vehicle.WithInsuranceExpiry(dbRow.InsuranceExpiry),
		vehicle.WithTimestamps(dbRow.CreatedAt, dbRow.UpdatedAt),
	)
}

func ToDBVehicle(domainEntity vehicle.Vehicle) *models.Vehicle {
	dbRow := &models.Vehicle{
		ID:                 domainEntity.ID(),
		TenantID:           domainEntity.TenantID(),
		Make:               domainEntity.Make(),
		Model:              domainEntity.Model(),
		Year:               domainEntity.Year(),
		VIN:                domainEntity.VIN(),
		LicensePlate:       domainEntity.LicensePlate(),
		Status:             int(domainEntity.Status()),
		CurrentOdometer:    domainEntity.CurrentOdometer(),
		RegistrationExpiry: domainEntity.RegistrationExpiry(),
		InsuranceExpiry:    domainEntity.InsuranceExpiry(),
		CreatedAt:          domainEntity.CreatedAt(),
		UpdatedAt:          domainEntity.UpdatedAt(),
	}
	return dbRow
}

func ToDomainDriver(dbRow *models.Driver) driver.Driver {
	opts := []driver.DriverOption{
		driver.WithDriverStatus(enums.DriverStatus(dbRow.Status)),
		driver.WithDriverTimestamps(dbRow.CreatedAt, dbRow.UpdatedAt),
	}

	if dbRow.UserID.Valid {
		opts = append(opts, driver.WithUserID(dbRow.UserID.Int64))
	}

	if dbRow.Phone.Valid {
		opts = append(opts, driver.WithPhone(dbRow.Phone.String))
	}

	if dbRow.Email.Valid {
		opts = append(opts, driver.WithEmail(dbRow.Email.String))
	}

	return driver.NewDriver(
		dbRow.ID,
		dbRow.TenantID,
		dbRow.FirstName,
		dbRow.LastName,
		dbRow.LicenseNumber,
		dbRow.LicenseExpiry,
		opts...,
	)
}

func ToDBDriver(domainEntity driver.Driver) *models.Driver {
	dbRow := &models.Driver{
		ID:            domainEntity.ID(),
		TenantID:      domainEntity.TenantID(),
		FirstName:     domainEntity.FirstName(),
		LastName:      domainEntity.LastName(),
		LicenseNumber: domainEntity.LicenseNumber(),
		LicenseExpiry: domainEntity.LicenseExpiry(),
		Status:        int(domainEntity.Status()),
		CreatedAt:     domainEntity.CreatedAt(),
		UpdatedAt:     domainEntity.UpdatedAt(),
	}

	if domainEntity.UserID() != nil {
		dbRow.UserID = sql.NullInt64{
			Int64: *domainEntity.UserID(),
			Valid: true,
		}
	}

	if domainEntity.Phone() != "" {
		dbRow.Phone = sql.NullString{
			String: domainEntity.Phone(),
			Valid:  true,
		}
	}

	if domainEntity.Email() != "" {
		dbRow.Email = sql.NullString{
			String: domainEntity.Email(),
			Valid:  true,
		}
	}

	return dbRow
}

func ToDomainTrip(dbRow *models.Trip) trip.Trip {
	purposeStr := ""
	if dbRow.Purpose.Valid {
		purposeStr = dbRow.Purpose.String
	}

	opts := []trip.TripOption{
		trip.WithStatus(enums.TripStatus(dbRow.Status)),
		trip.WithTimestamps(dbRow.CreatedAt, dbRow.UpdatedAt),
	}

	if dbRow.EndTime.Valid {
		opts = append(opts, trip.WithEndTime(dbRow.EndTime.Time))
	}

	if dbRow.EndOdometer.Valid {
		opts = append(opts, trip.WithEndOdometer(int(dbRow.EndOdometer.Int32)))
	}

	return trip.NewTrip(
		dbRow.ID,
		dbRow.TenantID,
		dbRow.VehicleID,
		dbRow.DriverID,
		dbRow.Origin,
		dbRow.Destination,
		purposeStr,
		dbRow.StartTime,
		dbRow.StartOdometer,
		opts...,
	)
}

func ToDBTrip(domainEntity trip.Trip) *models.Trip {
	dbRow := &models.Trip{
		ID:            domainEntity.ID(),
		TenantID:      domainEntity.TenantID(),
		VehicleID:     domainEntity.VehicleID(),
		DriverID:      domainEntity.DriverID(),
		Origin:        domainEntity.Origin(),
		Destination:   domainEntity.Destination(),
		StartTime:     domainEntity.StartTime(),
		StartOdometer: domainEntity.StartOdometer(),
		Status:        int(domainEntity.Status()),
		CreatedAt:     domainEntity.CreatedAt(),
		UpdatedAt:     domainEntity.UpdatedAt(),
	}

	if domainEntity.Purpose() != "" {
		dbRow.Purpose = sql.NullString{
			String: domainEntity.Purpose(),
			Valid:  true,
		}
	}

	if domainEntity.EndTime() != nil {
		dbRow.EndTime = sql.NullTime{
			Time:  *domainEntity.EndTime(),
			Valid: true,
		}
	}

	if domainEntity.EndOdometer() != nil {
		dbRow.EndOdometer = sql.NullInt32{
			Int32: int32(*domainEntity.EndOdometer()),
			Valid: true,
		}
	}

	return dbRow
}

func ToDomainMaintenance(dbRow *models.Maintenance) maintenance.Maintenance {
	serviceProviderStr := ""
	if dbRow.ServiceProvider.Valid {
		serviceProviderStr = dbRow.ServiceProvider.String
	}

	descriptionStr := ""
	if dbRow.Description.Valid {
		descriptionStr = dbRow.Description.String
	}

	opts := []maintenance.MaintenanceOption{
		maintenance.WithTimestamps(dbRow.CreatedAt, dbRow.UpdatedAt),
	}

	if dbRow.NextServiceDue.Valid {
		opts = append(opts, maintenance.WithNextServiceDue(dbRow.NextServiceDue.Time))
	}

	if dbRow.NextServiceOdometer.Valid {
		odo := int(dbRow.NextServiceOdometer.Int32)
		opts = append(opts, maintenance.WithNextServiceOdometer(odo))
	}

	return maintenance.NewMaintenance(
		dbRow.ID,
		dbRow.TenantID,
		dbRow.VehicleID,
		enums.ServiceType(dbRow.ServiceType),
		dbRow.ServiceDate,
		dbRow.Odometer,
		dbRow.Cost,
		serviceProviderStr,
		descriptionStr,
		opts...,
	)
}

func ToDBMaintenance(domainEntity maintenance.Maintenance) *models.Maintenance {
	dbRow := &models.Maintenance{
		ID:          domainEntity.ID(),
		TenantID:    domainEntity.TenantID(),
		VehicleID:   domainEntity.VehicleID(),
		ServiceType: int(domainEntity.ServiceType()),
		ServiceDate: domainEntity.ServiceDate(),
		Odometer:    domainEntity.Odometer(),
		Cost:        domainEntity.Cost(),
		CreatedAt:   domainEntity.CreatedAt(),
		UpdatedAt:   domainEntity.UpdatedAt(),
	}

	if domainEntity.ServiceProvider() != "" {
		dbRow.ServiceProvider = sql.NullString{
			String: domainEntity.ServiceProvider(),
			Valid:  true,
		}
	}

	if domainEntity.Description() != "" {
		dbRow.Description = sql.NullString{
			String: domainEntity.Description(),
			Valid:  true,
		}
	}

	if domainEntity.NextServiceDue() != nil {
		dbRow.NextServiceDue = sql.NullTime{
			Time:  *domainEntity.NextServiceDue(),
			Valid: true,
		}
	}

	if domainEntity.NextServiceOdometer() != nil {
		dbRow.NextServiceOdometer = sql.NullInt32{
			Int32: int32(*domainEntity.NextServiceOdometer()),
			Valid: true,
		}
	}

	return dbRow
}

func ToDomainFuelEntry(dbRow *models.FuelEntry) fuel_entry.FuelEntry {
	locationStr := ""
	if dbRow.Location.Valid {
		locationStr = dbRow.Location.String
	}

	opts := []fuel_entry.FuelEntryOption{
		fuel_entry.WithTimestamps(dbRow.CreatedAt, dbRow.UpdatedAt),
	}

	if dbRow.DriverID.Valid {
		opts = append(opts, fuel_entry.WithDriverID(dbRow.DriverID.UUID))
	}

	return fuel_entry.NewFuelEntry(
		dbRow.ID,
		dbRow.TenantID,
		dbRow.VehicleID,
		dbRow.Date,
		dbRow.Quantity,
		dbRow.Cost,
		dbRow.Odometer,
		enums.FuelType(dbRow.FuelType),
		locationStr,
		opts...,
	)
}

func ToDBFuelEntry(domainEntity fuel_entry.FuelEntry) *models.FuelEntry {
	dbRow := &models.FuelEntry{
		ID:        domainEntity.ID(),
		TenantID:  domainEntity.TenantID(),
		VehicleID: domainEntity.VehicleID(),
		Date:      domainEntity.Date(),
		Quantity:  domainEntity.Quantity(),
		Cost:      domainEntity.Cost(),
		Odometer:  domainEntity.Odometer(),
		FuelType:  int(domainEntity.FuelType()),
		CreatedAt: domainEntity.CreatedAt(),
		UpdatedAt: domainEntity.UpdatedAt(),
	}

	if domainEntity.DriverID() != nil {
		dbRow.DriverID = uuid.NullUUID{
			UUID:  *domainEntity.DriverID(),
			Valid: true,
		}
	}

	if domainEntity.Location() != "" {
		dbRow.Location = sql.NullString{
			String: domainEntity.Location(),
			Valid:  true,
		}
	}

	return dbRow
}
