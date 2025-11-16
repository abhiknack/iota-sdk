package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/driver"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/fuel_entry"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/maintenance"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/vehicle"
	"github.com/iota-uz/iota-sdk/pkg/eventbus"
	"github.com/sirupsen/logrus"
)

type NotificationType string

const (
	NotificationTypeLicenseExpiry      NotificationType = "license_expiry"
	NotificationTypeRegistrationExpiry NotificationType = "registration_expiry"
	NotificationTypeInsuranceExpiry    NotificationType = "insurance_expiry"
	NotificationTypeMaintenanceDue     NotificationType = "maintenance_due"
	NotificationTypeFuelAnomaly        NotificationType = "fuel_anomaly"
)

type Notification struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Type      NotificationType
	Title     string
	Message   string
	Data      map[string]interface{}
	CreatedAt time.Time
}

type NotificationService struct {
	vehicleService     *VehicleService
	driverService      *DriverService
	maintenanceService *MaintenanceService
	fuelService        *FuelService
	publisher          eventbus.EventBus
	logger             *logrus.Logger
}

func NewNotificationService(
	vehicleService *VehicleService,
	driverService *DriverService,
	maintenanceService *MaintenanceService,
	fuelService *FuelService,
	publisher eventbus.EventBus,
	logger *logrus.Logger,
) *NotificationService {
	return &NotificationService{
		vehicleService:     vehicleService,
		driverService:      driverService,
		maintenanceService: maintenanceService,
		fuelService:        fuelService,
		publisher:          publisher,
		logger:             logger,
	}
}

func (s *NotificationService) CheckExpiringLicenses(ctx context.Context, tenantID uuid.UUID, days int) error {
	drivers, err := s.driverService.GetExpiringLicenses(ctx, tenantID, days)
	if err != nil {
		return fmt.Errorf("failed to get expiring licenses: %w", err)
	}

	for _, d := range drivers {
		notification := s.createLicenseExpiryNotification(d, days)
		if err := s.sendNotification(ctx, notification); err != nil {
			s.logger.WithError(err).Error("failed to send license expiry notification")
		}
	}

	return nil
}

func (s *NotificationService) CheckExpiringRegistrations(ctx context.Context, tenantID uuid.UUID, days int) error {
	vehicles, err := s.vehicleService.GetExpiringRegistrations(ctx, tenantID, days)
	if err != nil {
		return fmt.Errorf("failed to get expiring registrations: %w", err)
	}

	for _, v := range vehicles {
		notification := s.createRegistrationExpiryNotification(v, days)
		if err := s.sendNotification(ctx, notification); err != nil {
			s.logger.WithError(err).Error("failed to send registration expiry notification")
		}
	}

	return nil
}

func (s *NotificationService) CheckExpiringInsurance(ctx context.Context, tenantID uuid.UUID, days int) error {
	vehicles, err := s.vehicleService.GetExpiringInsurance(ctx, tenantID, days)
	if err != nil {
		return fmt.Errorf("failed to get expiring insurance: %w", err)
	}

	for _, v := range vehicles {
		notification := s.createInsuranceExpiryNotification(v, days)
		if err := s.sendNotification(ctx, notification); err != nil {
			s.logger.WithError(err).Error("failed to send insurance expiry notification")
		}
	}

	return nil
}

func (s *NotificationService) CheckDueMaintenance(ctx context.Context, tenantID uuid.UUID) error {
	maintenanceRecords, err := s.maintenanceService.GetDueMaintenance(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get due maintenance: %w", err)
	}

	for _, m := range maintenanceRecords {
		notification := s.createMaintenanceDueNotification(m)
		if err := s.sendNotification(ctx, notification); err != nil {
			s.logger.WithError(err).Error("failed to send maintenance due notification")
		}
	}

	return nil
}

func (s *NotificationService) CheckFuelAnomaly(ctx context.Context, entry fuel_entry.FuelEntry) error {
	isAnomaly, err := s.fuelService.DetectAnomaly(ctx, entry)
	if err != nil {
		return fmt.Errorf("failed to detect fuel anomaly: %w", err)
	}

	if isAnomaly {
		notification := s.createFuelAnomalyNotification(entry)
		if err := s.sendNotification(ctx, notification); err != nil {
			s.logger.WithError(err).Error("failed to send fuel anomaly notification")
		}
	}

	return nil
}

func (s *NotificationService) createLicenseExpiryNotification(d driver.Driver, days int) *Notification {
	return &Notification{
		ID:       uuid.New(),
		TenantID: d.TenantID(),
		Type:     NotificationTypeLicenseExpiry,
		Title:    "Driver License Expiring Soon",
		Message:  fmt.Sprintf("Driver %s %s's license (number: %s) will expire in %d days on %s", d.FirstName(), d.LastName(), d.LicenseNumber(), days, d.LicenseExpiry().Format("2006-01-02")),
		Data: map[string]interface{}{
			"driver_id":         d.ID().String(),
			"driver_name":       fmt.Sprintf("%s %s", d.FirstName(), d.LastName()),
			"license_number":    d.LicenseNumber(),
			"expiry_date":       d.LicenseExpiry(),
			"days_until_expiry": days,
		},
		CreatedAt: time.Now(),
	}
}

func (s *NotificationService) createRegistrationExpiryNotification(v vehicle.Vehicle, days int) *Notification {
	return &Notification{
		ID:       uuid.New(),
		TenantID: v.TenantID(),
		Type:     NotificationTypeRegistrationExpiry,
		Title:    "Vehicle Registration Expiring Soon",
		Message:  fmt.Sprintf("Vehicle %s %s (plate: %s) registration will expire in %d days on %s", v.Make(), v.Model(), v.LicensePlate(), days, v.RegistrationExpiry().Format("2006-01-02")),
		Data: map[string]interface{}{
			"vehicle_id":        v.ID().String(),
			"vehicle_name":      fmt.Sprintf("%s %s", v.Make(), v.Model()),
			"license_plate":     v.LicensePlate(),
			"expiry_date":       v.RegistrationExpiry(),
			"days_until_expiry": days,
		},
		CreatedAt: time.Now(),
	}
}

func (s *NotificationService) createInsuranceExpiryNotification(v vehicle.Vehicle, days int) *Notification {
	return &Notification{
		ID:       uuid.New(),
		TenantID: v.TenantID(),
		Type:     NotificationTypeInsuranceExpiry,
		Title:    "Vehicle Insurance Expiring Soon",
		Message:  fmt.Sprintf("Vehicle %s %s (plate: %s) insurance will expire in %d days on %s", v.Make(), v.Model(), v.LicensePlate(), days, v.InsuranceExpiry().Format("2006-01-02")),
		Data: map[string]interface{}{
			"vehicle_id":        v.ID().String(),
			"vehicle_name":      fmt.Sprintf("%s %s", v.Make(), v.Model()),
			"license_plate":     v.LicensePlate(),
			"expiry_date":       v.InsuranceExpiry(),
			"days_until_expiry": days,
		},
		CreatedAt: time.Now(),
	}
}

func (s *NotificationService) createMaintenanceDueNotification(m maintenance.Maintenance) *Notification {
	var dueInfo string
	if m.NextServiceDue() != nil {
		dueInfo = fmt.Sprintf("due on %s", m.NextServiceDue().Format("2006-01-02"))
	} else if m.NextServiceOdometer() != nil {
		dueInfo = fmt.Sprintf("due at %d km", *m.NextServiceOdometer())
	} else {
		dueInfo = "due now"
	}

	return &Notification{
		ID:       uuid.New(),
		TenantID: m.TenantID(),
		Type:     NotificationTypeMaintenanceDue,
		Title:    "Vehicle Maintenance Due",
		Message:  fmt.Sprintf("Maintenance for vehicle (ID: %s) is %s. Service type: %s", m.VehicleID().String(), dueInfo, m.ServiceType().String()),
		Data: map[string]interface{}{
			"maintenance_id":        m.ID().String(),
			"vehicle_id":            m.VehicleID().String(),
			"service_type":          m.ServiceType().String(),
			"next_service_due":      m.NextServiceDue(),
			"next_service_odometer": m.NextServiceOdometer(),
		},
		CreatedAt: time.Now(),
	}
}

func (s *NotificationService) createFuelAnomalyNotification(f fuel_entry.FuelEntry) *Notification {
	return &Notification{
		ID:       uuid.New(),
		TenantID: f.TenantID(),
		Type:     NotificationTypeFuelAnomaly,
		Title:    "Fuel Efficiency Anomaly Detected",
		Message:  fmt.Sprintf("Unusual fuel efficiency detected for vehicle (ID: %s). Fuel entry on %s shows significant deviation from average.", f.VehicleID().String(), f.Date().Format("2006-01-02")),
		Data: map[string]interface{}{
			"fuel_entry_id": f.ID().String(),
			"vehicle_id":    f.VehicleID().String(),
			"date":          f.Date(),
			"quantity":      f.Quantity(),
			"cost":          f.Cost(),
			"odometer":      f.Odometer(),
		},
		CreatedAt: time.Now(),
	}
}

func (s *NotificationService) sendNotification(ctx context.Context, notification *Notification) error {
	s.logger.WithFields(logrus.Fields{
		"notification_id":   notification.ID,
		"tenant_id":         notification.TenantID,
		"notification_type": notification.Type,
		"title":             notification.Title,
		"message":           notification.Message,
		"data":              notification.Data,
	}).Info("Fleet notification")

	s.publisher.Publish(notification)

	return nil
}
