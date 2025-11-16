package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/maintenance"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/vehicle"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
	"github.com/iota-uz/iota-sdk/modules/fleet/permissions"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/eventbus"
)

type MaintenanceService struct {
	repo        maintenance.Repository
	vehicleRepo vehicle.Repository
	publisher   eventbus.EventBus
}

func NewMaintenanceService(
	repo maintenance.Repository,
	vehicleRepo vehicle.Repository,
	publisher eventbus.EventBus,
) *MaintenanceService {
	return &MaintenanceService{
		repo:        repo,
		vehicleRepo: vehicleRepo,
		publisher:   publisher,
	}
}

func (s *MaintenanceService) GetByID(ctx context.Context, id uuid.UUID) (maintenance.Maintenance, error) {
	if err := composables.CanUser(ctx, permissions.MaintenanceRead); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *MaintenanceService) GetPaginated(ctx context.Context, params *maintenance.FindParams) ([]maintenance.Maintenance, error) {
	if err := composables.CanUser(ctx, permissions.MaintenanceRead); err != nil {
		return nil, err
	}
	return s.repo.GetPaginated(ctx, params)
}

func (s *MaintenanceService) Count(ctx context.Context, params *maintenance.FindParams) (int64, error) {
	if err := composables.CanUser(ctx, permissions.MaintenanceRead); err != nil {
		return 0, err
	}
	return s.repo.Count(ctx, params)
}

func (s *MaintenanceService) GetByVehicle(ctx context.Context, vehicleID uuid.UUID) ([]maintenance.Maintenance, error) {
	if err := composables.CanUser(ctx, permissions.MaintenanceRead); err != nil {
		return nil, err
	}
	return s.repo.GetByVehicle(ctx, vehicleID)
}

func (s *MaintenanceService) GetDueMaintenance(ctx context.Context, tenantID uuid.UUID) ([]maintenance.Maintenance, error) {
	if err := composables.CanUser(ctx, permissions.MaintenanceRead); err != nil {
		return nil, err
	}
	return s.repo.GetDueMaintenance(ctx, tenantID)
}

func (s *MaintenanceService) Create(ctx context.Context, m maintenance.Maintenance) (maintenance.Maintenance, error) {
	if err := composables.CanUser(ctx, permissions.MaintenanceCreate); err != nil {
		return nil, err
	}

	v, err := s.vehicleRepo.GetByID(ctx, m.VehicleID())
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}

	nextServiceDate, nextServiceOdometer := s.calculateNextService(m, v)
	if nextServiceDate != nil || nextServiceOdometer != nil {
		m = m.UpdateNextService(nextServiceDate, nextServiceOdometer)
	}

	created, err := s.repo.Create(ctx, m)
	if err != nil {
		return nil, fmt.Errorf("failed to create maintenance: %w", err)
	}

	event, err := maintenance.NewMaintenanceCreatedEvent(ctx, created)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	event.Result = created
	s.publisher.Publish(event)

	return created, nil
}

func (s *MaintenanceService) Update(ctx context.Context, m maintenance.Maintenance) (maintenance.Maintenance, error) {
	if err := composables.CanUser(ctx, permissions.MaintenanceUpdate); err != nil {
		return nil, err
	}

	updated, err := s.repo.Update(ctx, m)
	if err != nil {
		return nil, fmt.Errorf("failed to update maintenance: %w", err)
	}

	event, err := maintenance.NewMaintenanceUpdatedEvent(ctx, updated)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	event.Result = updated
	s.publisher.Publish(event)

	return updated, nil
}

func (s *MaintenanceService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := composables.CanUser(ctx, permissions.MaintenanceDelete); err != nil {
		return err
	}

	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get maintenance: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete maintenance: %w", err)
	}

	event, err := maintenance.NewMaintenanceDeletedEvent(ctx)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	event.Result = m
	s.publisher.Publish(event)

	return nil
}

func (s *MaintenanceService) calculateNextService(m maintenance.Maintenance, v vehicle.Vehicle) (*time.Time, *int) {
	var nextDate *time.Time
	var nextOdometer *int

	switch m.ServiceType() {
	case enums.ServiceTypeOilChange:
		date := m.ServiceDate().AddDate(0, 6, 0)
		nextDate = &date
		odo := m.Odometer() + 10000
		nextOdometer = &odo

	case enums.ServiceTypeTireRotation:
		date := m.ServiceDate().AddDate(0, 6, 0)
		nextDate = &date
		odo := m.Odometer() + 12000
		nextOdometer = &odo

	case enums.ServiceTypeBrakeService:
		date := m.ServiceDate().AddDate(1, 0, 0)
		nextDate = &date
		odo := m.Odometer() + 20000
		nextOdometer = &odo

	case enums.ServiceTypeInspection:
		date := m.ServiceDate().AddDate(1, 0, 0)
		nextDate = &date

	case enums.ServiceTypeRepair:
	case enums.ServiceTypeOther:
	}

	return nextDate, nextOdometer
}
