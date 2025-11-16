package mappers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/components/base/badge"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/driver"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/fuel_entry"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/maintenance"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/trip"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/vehicle"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/controllers/dtos"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/viewmodels"
)

func VehicleToListViewModel(v vehicle.Vehicle) viewmodels.VehicleListViewModel {
	return viewmodels.VehicleListViewModel{
		ID:                 v.ID().String(),
		Make:               v.Make(),
		Model:              v.Model(),
		Year:               v.Year(),
		LicensePlate:       v.LicensePlate(),
		Status:             v.Status().String(),
		StatusBadgeVariant: GetVehicleStatusBadgeVariant(v.Status()),
		CurrentOdometer:    v.CurrentOdometer(),
		RegistrationExpiry: v.RegistrationExpiry().Format("2006-01-02"),
		InsuranceExpiry:    v.InsuranceExpiry().Format("2006-01-02"),
	}
}

func VehicleToDetailViewModel(v vehicle.Vehicle) viewmodels.VehicleDetailViewModel {
	return viewmodels.VehicleDetailViewModel{
		ID:                 v.ID().String(),
		TenantID:           v.TenantID().String(),
		Make:               v.Make(),
		Model:              v.Model(),
		Year:               v.Year(),
		VIN:                v.VIN(),
		LicensePlate:       v.LicensePlate(),
		Status:             v.Status().String(),
		StatusBadgeVariant: GetVehicleStatusBadgeVariant(v.Status()),
		CurrentOdometer:    v.CurrentOdometer(),
		RegistrationExpiry: v.RegistrationExpiry().Format("2006-01-02"),
		InsuranceExpiry:    v.InsuranceExpiry().Format("2006-01-02"),
		CreatedAt:          v.CreatedAt().Format("2006-01-02 15:04:05"),
		UpdatedAt:          v.UpdatedAt().Format("2006-01-02 15:04:05"),
		FullName:           fmt.Sprintf("%s %s (%d)", v.Make(), v.Model(), v.Year()),
	}
}

func VehicleToFormViewModel(v vehicle.Vehicle, errors map[string]string) viewmodels.VehicleFormViewModel {
	return viewmodels.VehicleFormViewModel{
		ID:                 v.ID().String(),
		Make:               v.Make(),
		Model:              v.Model(),
		Year:               v.Year(),
		VIN:                v.VIN(),
		LicensePlate:       v.LicensePlate(),
		CurrentOdometer:    v.CurrentOdometer(),
		RegistrationExpiry: v.RegistrationExpiry().Format("2006-01-02"),
		InsuranceExpiry:    v.InsuranceExpiry().Format("2006-01-02"),
		IsEdit:             true,
		Errors:             errors,
	}
}

func GetVehicleStatusBadgeVariant(status enums.VehicleStatus) badge.Variant {
	switch status {
	case enums.VehicleStatusAvailable:
		return badge.VariantGreen
	case enums.VehicleStatusInUse:
		return badge.VariantBlue
	case enums.VehicleStatusMaintenance:
		return badge.VariantYellow
	case enums.VehicleStatusOutOfService:
		return badge.VariantPink
	case enums.VehicleStatusRetired:
		return badge.VariantGray
	default:
		return badge.VariantGray
	}
}

func DriverToListViewModel(d driver.Driver) viewmodels.DriverListViewModel {
	licenseExpiryWarning := time.Until(d.LicenseExpiry()) < 30*24*time.Hour
	return viewmodels.DriverListViewModel{
		ID:                   d.ID().String(),
		FirstName:            d.FirstName(),
		LastName:             d.LastName(),
		FullName:             fmt.Sprintf("%s %s", d.FirstName(), d.LastName()),
		LicenseNumber:        d.LicenseNumber(),
		LicenseExpiry:        d.LicenseExpiry().Format("2006-01-02"),
		LicenseExpiryWarning: licenseExpiryWarning,
		Phone:                d.Phone(),
		Email:                d.Email(),
		Status:               d.Status().String(),
		StatusBadgeVariant:   GetDriverStatusBadgeVariant(d.Status()),
	}
}

func DriverToDetailViewModel(d driver.Driver, totalTrips, activeTrips int) viewmodels.DriverDetailViewModel {
	licenseExpiryWarning := time.Until(d.LicenseExpiry()) < 30*24*time.Hour
	userID := ""
	if d.UserID() != nil {
		userID = fmt.Sprintf("%d", *d.UserID())
	}
	return viewmodels.DriverDetailViewModel{
		ID:                   d.ID().String(),
		TenantID:             d.TenantID().String(),
		UserID:               userID,
		FirstName:            d.FirstName(),
		LastName:             d.LastName(),
		FullName:             fmt.Sprintf("%s %s", d.FirstName(), d.LastName()),
		LicenseNumber:        d.LicenseNumber(),
		LicenseExpiry:        d.LicenseExpiry().Format("2006-01-02"),
		LicenseExpiryWarning: licenseExpiryWarning,
		Phone:                d.Phone(),
		Email:                d.Email(),
		Status:               d.Status().String(),
		StatusBadgeVariant:   GetDriverStatusBadgeVariant(d.Status()),
		CreatedAt:            d.CreatedAt().Format("2006-01-02 15:04:05"),
		UpdatedAt:            d.UpdatedAt().Format("2006-01-02 15:04:05"),
		TotalTrips:           totalTrips,
		ActiveTrips:          activeTrips,
	}
}

func DriverToFormViewModel(d driver.Driver, errors map[string]string) viewmodels.DriverFormViewModel {
	userID := ""
	if d.UserID() != nil {
		userID = fmt.Sprintf("%d", *d.UserID())
	}
	return viewmodels.DriverFormViewModel{
		ID:            d.ID().String(),
		UserID:        userID,
		FirstName:     d.FirstName(),
		LastName:      d.LastName(),
		LicenseNumber: d.LicenseNumber(),
		LicenseExpiry: d.LicenseExpiry().Format("2006-01-02"),
		Phone:         d.Phone(),
		Email:         d.Email(),
		Status:        d.Status().String(),
		IsEdit:        true,
		Errors:        errors,
	}
}

func GetDriverStatusBadgeVariant(status enums.DriverStatus) badge.Variant {
	switch status {
	case enums.DriverStatusActive:
		return badge.VariantGreen
	case enums.DriverStatusInactive:
		return badge.VariantGray
	case enums.DriverStatusOnLeave:
		return badge.VariantYellow
	case enums.DriverStatusTerminated:
		return badge.VariantPink
	default:
		return badge.VariantGray
	}
}

func TripToListViewModel(t trip.Trip, vehicleName, driverName string) viewmodels.TripListViewModel {
	endTime := ""
	if t.EndTime() != nil {
		endTime = t.EndTime().Format("2006-01-02 15:04")
	}

	distance := 0
	if t.EndOdometer() != nil {
		distance = *t.EndOdometer() - t.StartOdometer()
	}

	duration := ""
	if t.EndTime() != nil {
		d := t.EndTime().Sub(t.StartTime())
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		duration = fmt.Sprintf("%dh %dm", hours, minutes)
	}

	return viewmodels.TripListViewModel{
		ID:                 t.ID().String(),
		VehicleID:          t.VehicleID().String(),
		VehicleName:        vehicleName,
		DriverID:           t.DriverID().String(),
		DriverName:         driverName,
		Origin:             t.Origin(),
		Destination:        t.Destination(),
		StartTime:          t.StartTime().Format("2006-01-02 15:04"),
		EndTime:            endTime,
		Status:             t.Status().String(),
		StatusBadgeVariant: GetTripStatusBadgeVariant(t.Status()),
		Distance:           distance,
		Duration:           duration,
	}
}

func TripToDetailViewModel(t trip.Trip, vehicleName, driverName string) viewmodels.TripDetailViewModel {
	endTime := ""
	if t.EndTime() != nil {
		endTime = t.EndTime().Format("2006-01-02 15:04")
	}

	endOdometer := 0
	if t.EndOdometer() != nil {
		endOdometer = *t.EndOdometer()
	}

	distance := 0
	if t.EndOdometer() != nil {
		distance = *t.EndOdometer() - t.StartOdometer()
	}

	duration := ""
	averageSpeed := 0.0
	if t.EndTime() != nil {
		d := t.EndTime().Sub(t.StartTime())
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		duration = fmt.Sprintf("%dh %dm", hours, minutes)

		if d.Hours() > 0 && distance > 0 {
			averageSpeed = float64(distance) / d.Hours()
		}
	}

	canComplete := t.Status() == enums.TripStatusInProgress
	canCancel := t.Status() == enums.TripStatusScheduled || t.Status() == enums.TripStatusInProgress

	return viewmodels.TripDetailViewModel{
		ID:                 t.ID().String(),
		TenantID:           t.TenantID().String(),
		VehicleID:          t.VehicleID().String(),
		VehicleName:        vehicleName,
		DriverID:           t.DriverID().String(),
		DriverName:         driverName,
		Origin:             t.Origin(),
		Destination:        t.Destination(),
		Purpose:            t.Purpose(),
		StartTime:          t.StartTime().Format("2006-01-02 15:04"),
		EndTime:            endTime,
		StartOdometer:      t.StartOdometer(),
		EndOdometer:        endOdometer,
		Status:             t.Status().String(),
		StatusBadgeVariant: GetTripStatusBadgeVariant(t.Status()),
		Distance:           distance,
		Duration:           duration,
		AverageSpeed:       averageSpeed,
		CreatedAt:          t.CreatedAt().Format("2006-01-02 15:04:05"),
		UpdatedAt:          t.UpdatedAt().Format("2006-01-02 15:04:05"),
		CanComplete:        canComplete,
		CanCancel:          canCancel,
	}
}

func TripToFormViewModel(t trip.Trip, errors map[string]string, vehicles []viewmodels.VehicleOption, drivers []viewmodels.DriverOption) viewmodels.TripFormViewModel {
	return viewmodels.TripFormViewModel{
		ID:            t.ID().String(),
		VehicleID:     t.VehicleID().String(),
		DriverID:      t.DriverID().String(),
		Origin:        t.Origin(),
		Destination:   t.Destination(),
		Purpose:       t.Purpose(),
		StartTime:     t.StartTime().Format("2006-01-02T15:04"),
		StartOdometer: t.StartOdometer(),
		IsEdit:        true,
		Errors:        errors,
		Vehicles:      vehicles,
		Drivers:       drivers,
		HasConflict:   false,
		ConflictMsg:   "",
	}
}

func GetTripStatusBadgeVariant(status enums.TripStatus) badge.Variant {
	switch status {
	case enums.TripStatusScheduled:
		return badge.VariantBlue
	case enums.TripStatusInProgress:
		return badge.VariantYellow
	case enums.TripStatusCompleted:
		return badge.VariantGreen
	case enums.TripStatusCancelled:
		return badge.VariantPink
	default:
		return badge.VariantGray
	}
}

func MaintenanceToListViewModel(m maintenance.Maintenance, vehicleName string) viewmodels.MaintenanceListViewModel {
	nextServiceDue := ""
	if m.NextServiceDue() != nil {
		nextServiceDue = m.NextServiceDue().Format("2006-01-02")
	}

	nextServiceOdometer := 0
	if m.NextServiceOdometer() != nil {
		nextServiceOdometer = *m.NextServiceOdometer()
	}

	isDue := false
	if m.NextServiceDue() != nil {
		isDue = time.Now().After(*m.NextServiceDue())
	}

	return viewmodels.MaintenanceListViewModel{
		ID:                  m.ID().String(),
		VehicleID:           m.VehicleID().String(),
		VehicleName:         vehicleName,
		ServiceType:         m.ServiceType().String(),
		ServiceDate:         m.ServiceDate().Format("2006-01-02"),
		Odometer:            m.Odometer(),
		Cost:                fmt.Sprintf("%.2f", m.Cost()),
		ServiceProvider:     m.ServiceProvider(),
		NextServiceDue:      nextServiceDue,
		NextServiceOdometer: nextServiceOdometer,
		IsDue:               isDue,
	}
}

func MaintenanceToDetailViewModel(m maintenance.Maintenance, vehicleName string) viewmodels.MaintenanceDetailViewModel {
	nextServiceDue := ""
	if m.NextServiceDue() != nil {
		nextServiceDue = m.NextServiceDue().Format("2006-01-02")
	}

	nextServiceOdometer := 0
	if m.NextServiceOdometer() != nil {
		nextServiceOdometer = *m.NextServiceOdometer()
	}

	return viewmodels.MaintenanceDetailViewModel{
		ID:                  m.ID().String(),
		TenantID:            m.TenantID().String(),
		VehicleID:           m.VehicleID().String(),
		VehicleName:         vehicleName,
		ServiceType:         m.ServiceType().String(),
		ServiceDate:         m.ServiceDate().Format("2006-01-02"),
		Odometer:            m.Odometer(),
		Cost:                fmt.Sprintf("%.2f", m.Cost()),
		ServiceProvider:     m.ServiceProvider(),
		Description:         m.Description(),
		NextServiceDue:      nextServiceDue,
		NextServiceOdometer: nextServiceOdometer,
		CreatedAt:           m.CreatedAt().Format("2006-01-02 15:04:05"),
		UpdatedAt:           m.UpdatedAt().Format("2006-01-02 15:04:05"),
	}
}

func MaintenanceToFormViewModel(m maintenance.Maintenance, errors map[string]string) viewmodels.MaintenanceFormViewModel {
	nextServiceDue := ""
	if m.NextServiceDue() != nil {
		nextServiceDue = m.NextServiceDue().Format("2006-01-02")
	}

	nextServiceOdometer := ""
	if m.NextServiceOdometer() != nil {
		nextServiceOdometer = fmt.Sprintf("%d", *m.NextServiceOdometer())
	}

	return viewmodels.MaintenanceFormViewModel{
		ID:                  m.ID().String(),
		VehicleID:           m.VehicleID().String(),
		ServiceType:         m.ServiceType().String(),
		ServiceDate:         m.ServiceDate().Format("2006-01-02"),
		Odometer:            m.Odometer(),
		Cost:                fmt.Sprintf("%.2f", m.Cost()),
		ServiceProvider:     m.ServiceProvider(),
		Description:         m.Description(),
		NextServiceDue:      nextServiceDue,
		NextServiceOdometer: nextServiceOdometer,
		IsEdit:              true,
		Errors:              errors,
	}
}

func FuelEntryToListViewModel(f fuel_entry.FuelEntry, vehicleName, driverName string, efficiency float64) viewmodels.FuelEntryListViewModel {
	return viewmodels.FuelEntryListViewModel{
		ID:          f.ID().String(),
		VehicleID:   f.VehicleID().String(),
		VehicleName: vehicleName,
		DriverName:  driverName,
		Date:        f.Date().Format("2006-01-02"),
		Quantity:    fmt.Sprintf("%.2f", f.Quantity()),
		Cost:        fmt.Sprintf("%.2f", f.Cost()),
		Odometer:    f.Odometer(),
		FuelType:    f.FuelType().String(),
		Location:    f.Location(),
		Efficiency:  fmt.Sprintf("%.2f", efficiency),
	}
}

func FuelEntryToDetailViewModel(f fuel_entry.FuelEntry, vehicleName, driverName string, efficiency float64) viewmodels.FuelEntryDetailViewModel {
	driverID := ""
	if f.DriverID() != nil {
		driverID = f.DriverID().String()
	}

	return viewmodels.FuelEntryDetailViewModel{
		ID:          f.ID().String(),
		TenantID:    f.TenantID().String(),
		VehicleID:   f.VehicleID().String(),
		VehicleName: vehicleName,
		DriverID:    driverID,
		DriverName:  driverName,
		Date:        f.Date().Format("2006-01-02"),
		Quantity:    fmt.Sprintf("%.2f", f.Quantity()),
		Cost:        fmt.Sprintf("%.2f", f.Cost()),
		Odometer:    f.Odometer(),
		FuelType:    f.FuelType().String(),
		Location:    f.Location(),
		Efficiency:  fmt.Sprintf("%.2f", efficiency),
		CreatedAt:   f.CreatedAt().Format("2006-01-02 15:04:05"),
		UpdatedAt:   f.UpdatedAt().Format("2006-01-02 15:04:05"),
	}
}

func FuelEntryToFormViewModel(f fuel_entry.FuelEntry, errors map[string]string) viewmodels.FuelEntryFormViewModel {
	driverID := ""
	if f.DriverID() != nil {
		driverID = f.DriverID().String()
	}

	return viewmodels.FuelEntryFormViewModel{
		ID:        f.ID().String(),
		VehicleID: f.VehicleID().String(),
		DriverID:  driverID,
		Date:      f.Date().Format("2006-01-02"),
		Quantity:  fmt.Sprintf("%.2f", f.Quantity()),
		Cost:      fmt.Sprintf("%.2f", f.Cost()),
		Odometer:  f.Odometer(),
		FuelType:  f.FuelType().String(),
		Location:  f.Location(),
		IsEdit:    true,
		Errors:    errors,
	}
}

func VehicleCreateDTOToDomain(dto dtos.VehicleCreateDTO, tenantID uuid.UUID) vehicle.Vehicle {
	return vehicle.NewVehicle(
		uuid.New(),
		tenantID,
		dto.Make,
		dto.Model,
		dto.Year,
		dto.VIN,
		dto.LicensePlate,
		vehicle.WithOdometer(dto.CurrentOdometer),
		vehicle.WithRegistrationExpiry(dto.RegistrationExpiry),
		vehicle.WithInsuranceExpiry(dto.InsuranceExpiry),
	)
}

func VehicleUpdateDTOToDomain(dto dtos.VehicleUpdateDTO, existing vehicle.Vehicle) vehicle.Vehicle {
	v := existing.UpdateDetails(dto.Make, dto.Model, dto.Year)
	v = v.UpdateOdometer(dto.CurrentOdometer)
	return vehicle.NewVehicle(
		dto.ID,
		existing.TenantID(),
		dto.Make,
		dto.Model,
		dto.Year,
		dto.VIN,
		dto.LicensePlate,
		vehicle.WithStatus(existing.Status()),
		vehicle.WithOdometer(dto.CurrentOdometer),
		vehicle.WithRegistrationExpiry(dto.RegistrationExpiry),
		vehicle.WithInsuranceExpiry(dto.InsuranceExpiry),
		vehicle.WithTimestamps(existing.CreatedAt(), time.Now()),
	)
}

func DriverCreateDTOToDomain(dto dtos.DriverCreateDTO, tenantID uuid.UUID) driver.Driver {
	opts := []driver.DriverOption{
		driver.WithPhone(dto.Phone),
		driver.WithEmail(dto.Email),
	}
	if dto.UserID != nil {
		opts = append(opts, driver.WithUserID(*dto.UserID))
	}
	return driver.NewDriver(
		uuid.New(),
		tenantID,
		dto.FirstName,
		dto.LastName,
		dto.LicenseNumber,
		dto.LicenseExpiry,
		opts...,
	)
}

func DriverUpdateDTOToDomain(dto dtos.DriverUpdateDTO, existing driver.Driver) (driver.Driver, error) {
	status, err := enums.ParseDriverStatus(dto.Status)
	if err != nil {
		return nil, err
	}

	opts := []driver.DriverOption{
		driver.WithPhone(dto.Phone),
		driver.WithEmail(dto.Email),
		driver.WithDriverStatus(status),
		driver.WithDriverTimestamps(existing.CreatedAt(), time.Now()),
	}
	if dto.UserID != nil {
		opts = append(opts, driver.WithUserID(*dto.UserID))
	}

	return driver.NewDriver(
		dto.ID,
		existing.TenantID(),
		dto.FirstName,
		dto.LastName,
		dto.LicenseNumber,
		dto.LicenseExpiry,
		opts...,
	), nil
}

func TripCreateDTOToDomain(dto dtos.TripCreateDTO, tenantID uuid.UUID) trip.Trip {
	return trip.NewTrip(
		uuid.New(),
		tenantID,
		dto.VehicleID,
		dto.DriverID,
		dto.Origin,
		dto.Destination,
		dto.Purpose,
		dto.StartTime,
		dto.StartOdometer,
	)
}

func MaintenanceCreateDTOToDomain(dto dtos.MaintenanceCreateDTO, tenantID uuid.UUID) (maintenance.Maintenance, error) {
	serviceType, err := enums.ParseServiceType(dto.ServiceType)
	if err != nil {
		return nil, err
	}

	opts := []maintenance.MaintenanceOption{}
	if dto.NextServiceDue != nil {
		opts = append(opts, maintenance.WithNextServiceDue(*dto.NextServiceDue))
	}
	if dto.NextServiceOdometer != nil {
		opts = append(opts, maintenance.WithNextServiceOdometer(*dto.NextServiceOdometer))
	}

	return maintenance.NewMaintenance(
		uuid.New(),
		tenantID,
		dto.VehicleID,
		serviceType,
		dto.ServiceDate,
		dto.Odometer,
		dto.Cost,
		dto.ServiceProvider,
		dto.Description,
		opts...,
	), nil
}

func MaintenanceUpdateDTOToDomain(dto dtos.MaintenanceUpdateDTO, existing maintenance.Maintenance) (maintenance.Maintenance, error) {
	serviceType, err := enums.ParseServiceType(dto.ServiceType)
	if err != nil {
		return nil, err
	}

	opts := []maintenance.MaintenanceOption{
		maintenance.WithTimestamps(existing.CreatedAt(), time.Now()),
	}
	if dto.NextServiceDue != nil {
		opts = append(opts, maintenance.WithNextServiceDue(*dto.NextServiceDue))
	}
	if dto.NextServiceOdometer != nil {
		opts = append(opts, maintenance.WithNextServiceOdometer(*dto.NextServiceOdometer))
	}

	return maintenance.NewMaintenance(
		dto.ID,
		existing.TenantID(),
		dto.VehicleID,
		serviceType,
		dto.ServiceDate,
		dto.Odometer,
		dto.Cost,
		dto.ServiceProvider,
		dto.Description,
		opts...,
	), nil
}

func FuelEntryCreateDTOToDomain(dto dtos.FuelEntryCreateDTO, tenantID uuid.UUID) (fuel_entry.FuelEntry, error) {
	fuelType, err := enums.ParseFuelType(dto.FuelType)
	if err != nil {
		return nil, err
	}

	opts := []fuel_entry.FuelEntryOption{}
	if dto.DriverID != nil {
		opts = append(opts, fuel_entry.WithDriverID(*dto.DriverID))
	}

	return fuel_entry.NewFuelEntry(
		uuid.New(),
		tenantID,
		dto.VehicleID,
		dto.Date,
		dto.Quantity,
		dto.Cost,
		dto.Odometer,
		fuelType,
		dto.Location,
		opts...,
	), nil
}

func FuelEntryUpdateDTOToDomain(dto dtos.FuelEntryUpdateDTO, existing fuel_entry.FuelEntry) (fuel_entry.FuelEntry, error) {
	fuelType, err := enums.ParseFuelType(dto.FuelType)
	if err != nil {
		return nil, err
	}

	opts := []fuel_entry.FuelEntryOption{
		fuel_entry.WithTimestamps(existing.CreatedAt(), time.Now()),
	}
	if dto.DriverID != nil {
		opts = append(opts, fuel_entry.WithDriverID(*dto.DriverID))
	}

	return fuel_entry.NewFuelEntry(
		dto.ID,
		existing.TenantID(),
		dto.VehicleID,
		dto.Date,
		dto.Quantity,
		dto.Cost,
		dto.Odometer,
		fuelType,
		dto.Location,
		opts...,
	), nil
}
