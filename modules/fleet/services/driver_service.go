package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/driver"
	"github.com/iota-uz/iota-sdk/modules/fleet/permissions"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/eventbus"
)

type DriverService struct {
	repo      driver.Repository
	publisher eventbus.EventBus
}

func NewDriverService(
	repo driver.Repository,
	publisher eventbus.EventBus,
) *DriverService {
	return &DriverService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *DriverService) GetByID(ctx context.Context, id uuid.UUID) (driver.Driver, error) {
	if err := composables.CanUser(ctx, permissions.DriverRead); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *DriverService) GetByUserID(ctx context.Context, userID uuid.UUID) (driver.Driver, error) {
	if err := composables.CanUser(ctx, permissions.DriverRead); err != nil {
		return nil, err
	}
	return s.repo.GetByUserID(ctx, userID)
}

func (s *DriverService) GetPaginated(ctx context.Context, params *driver.FindParams) ([]driver.Driver, error) {
	if err := composables.CanUser(ctx, permissions.DriverRead); err != nil {
		return nil, err
	}
	return s.repo.GetPaginated(ctx, params)
}

func (s *DriverService) Count(ctx context.Context, params *driver.FindParams) (int64, error) {
	if err := composables.CanUser(ctx, permissions.DriverRead); err != nil {
		return 0, err
	}
	return s.repo.Count(ctx, params)
}

func (s *DriverService) GetExpiringLicenses(ctx context.Context, tenantID uuid.UUID, days int) ([]driver.Driver, error) {
	if err := composables.CanUser(ctx, permissions.DriverRead); err != nil {
		return nil, err
	}
	return s.repo.GetExpiringLicenses(ctx, tenantID, days)
}

func (s *DriverService) GetAvailable(ctx context.Context, tenantID uuid.UUID, startTime, endTime time.Time) ([]driver.Driver, error) {
	if err := composables.CanUser(ctx, permissions.DriverRead); err != nil {
		return nil, err
	}
	return s.repo.GetAvailable(ctx, tenantID, startTime, endTime)
}

func (s *DriverService) Create(ctx context.Context, d driver.Driver) (driver.Driver, error) {
	if err := composables.CanUser(ctx, permissions.DriverCreate); err != nil {
		return nil, err
	}

	if err := s.validateLicense(d); err != nil {
		return nil, err
	}

	created, err := s.repo.Create(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	event, err := driver.NewDriverCreatedEvent(ctx, created)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	event.Result = created
	s.publisher.Publish(event)

	return created, nil
}

func (s *DriverService) Update(ctx context.Context, d driver.Driver) (driver.Driver, error) {
	if err := composables.CanUser(ctx, permissions.DriverUpdate); err != nil {
		return nil, err
	}

	if err := s.validateLicense(d); err != nil {
		return nil, err
	}

	updated, err := s.repo.Update(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("failed to update driver: %w", err)
	}

	event, err := driver.NewDriverUpdatedEvent(ctx, updated)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	event.Result = updated
	s.publisher.Publish(event)

	return updated, nil
}

func (s *DriverService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := composables.CanUser(ctx, permissions.DriverDelete); err != nil {
		return err
	}

	d, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get driver: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete driver: %w", err)
	}

	event, err := driver.NewDriverDeletedEvent(ctx)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	event.Result = d
	s.publisher.Publish(event)

	return nil
}

func (s *DriverService) validateLicense(d driver.Driver) error {
	if d.LicenseNumber() == "" {
		return fmt.Errorf("license number is required")
	}

	if d.LicenseExpiry().Before(time.Now()) {
		return fmt.Errorf("driver license has expired")
	}

	return nil
}
