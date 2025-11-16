package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/vehicle"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
	"github.com/iota-uz/iota-sdk/modules/fleet/permissions"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/eventbus"
)

type VehicleService struct {
	repo      vehicle.Repository
	publisher eventbus.EventBus
}

func NewVehicleService(
	repo vehicle.Repository,
	publisher eventbus.EventBus,
) *VehicleService {
	return &VehicleService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *VehicleService) GetByID(ctx context.Context, id uuid.UUID) (vehicle.Vehicle, error) {
	if err := composables.CanUser(ctx, permissions.VehicleRead); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *VehicleService) GetPaginated(ctx context.Context, params *vehicle.FindParams) ([]vehicle.Vehicle, error) {
	if err := composables.CanUser(ctx, permissions.VehicleRead); err != nil {
		return nil, err
	}
	return s.repo.GetPaginated(ctx, params)
}

func (s *VehicleService) Count(ctx context.Context, params *vehicle.FindParams) (int64, error) {
	if err := composables.CanUser(ctx, permissions.VehicleRead); err != nil {
		return 0, err
	}
	return s.repo.Count(ctx, params)
}

func (s *VehicleService) GetByStatus(ctx context.Context, tenantID uuid.UUID, status enums.VehicleStatus) ([]vehicle.Vehicle, error) {
	if err := composables.CanUser(ctx, permissions.VehicleRead); err != nil {
		return nil, err
	}
	return s.repo.GetByStatus(ctx, tenantID, status)
}

func (s *VehicleService) GetExpiringRegistrations(ctx context.Context, tenantID uuid.UUID, days int) ([]vehicle.Vehicle, error) {
	if err := composables.CanUser(ctx, permissions.VehicleRead); err != nil {
		return nil, err
	}
	return s.repo.GetExpiringRegistrations(ctx, tenantID, days)
}

func (s *VehicleService) Create(ctx context.Context, v vehicle.Vehicle) (vehicle.Vehicle, error) {
	if err := composables.CanUser(ctx, permissions.VehicleCreate); err != nil {
		return nil, err
	}

	created, err := s.repo.Create(ctx, v)
	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle: %w", err)
	}

	event, err := vehicle.NewVehicleCreatedEvent(ctx, created)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	event.Result = created
	s.publisher.Publish(event)

	return created, nil
}

func (s *VehicleService) Update(ctx context.Context, v vehicle.Vehicle) (vehicle.Vehicle, error) {
	if err := composables.CanUser(ctx, permissions.VehicleUpdate); err != nil {
		return nil, err
	}

	updated, err := s.repo.Update(ctx, v)
	if err != nil {
		return nil, fmt.Errorf("failed to update vehicle: %w", err)
	}

	event, err := vehicle.NewVehicleUpdatedEvent(ctx, updated)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	event.Result = updated
	s.publisher.Publish(event)

	return updated, nil
}

func (s *VehicleService) UpdateStatus(ctx context.Context, id uuid.UUID, newStatus enums.VehicleStatus) (vehicle.Vehicle, error) {
	if err := composables.CanUser(ctx, permissions.VehicleUpdate); err != nil {
		return nil, err
	}

	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}

	oldStatus := v.Status()
	if err := s.validateStatusTransition(oldStatus, newStatus); err != nil {
		return nil, err
	}

	updated := v.UpdateStatus(newStatus)
	updated, err = s.repo.Update(ctx, updated)
	if err != nil {
		return nil, fmt.Errorf("failed to update vehicle status: %w", err)
	}

	event, err := vehicle.NewVehicleStatusChangedEvent(ctx, updated, oldStatus, newStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	s.publisher.Publish(event)

	return updated, nil
}

func (s *VehicleService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := composables.CanUser(ctx, permissions.VehicleDelete); err != nil {
		return err
	}

	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get vehicle: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}

	event, err := vehicle.NewVehicleDeletedEvent(ctx)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	event.Result = v
	s.publisher.Publish(event)

	return nil
}

func (s *VehicleService) validateStatusTransition(oldStatus, newStatus enums.VehicleStatus) error {
	if oldStatus == newStatus {
		return nil
	}

	switch oldStatus {
	case enums.VehicleStatusRetired:
		return fmt.Errorf("cannot change status of retired vehicle")
	case enums.VehicleStatusOutOfService:
		if newStatus == enums.VehicleStatusInUse {
			return fmt.Errorf("vehicle must be available before being used")
		}
	case enums.VehicleStatusMaintenance:
		if newStatus == enums.VehicleStatusInUse {
			return fmt.Errorf("vehicle must be available before being used")
		}
	}

	return nil
}

func (s *VehicleService) GetExpiringInsurance(ctx context.Context, tenantID uuid.UUID, days int) ([]vehicle.Vehicle, error) {
	if err := composables.CanUser(ctx, permissions.VehicleRead); err != nil {
		return nil, err
	}

	vehicles, err := s.repo.GetPaginated(ctx, &vehicle.FindParams{
		TenantID: tenantID,
		Limit:    1000,
	})
	if err != nil {
		return nil, err
	}

	threshold := time.Now().AddDate(0, 0, days)
	var expiring []vehicle.Vehicle
	for _, v := range vehicles {
		if v.InsuranceExpiry().Before(threshold) {
			expiring = append(expiring, v)
		}
	}

	return expiring, nil
}
