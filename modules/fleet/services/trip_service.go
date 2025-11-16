package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/driver"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/trip"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/vehicle"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
	"github.com/iota-uz/iota-sdk/modules/fleet/permissions"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/eventbus"
)

type TripService struct {
	tripRepo    trip.Repository
	vehicleRepo vehicle.Repository
	driverRepo  driver.Repository
	publisher   eventbus.EventBus
}

func NewTripService(
	tripRepo trip.Repository,
	vehicleRepo vehicle.Repository,
	driverRepo driver.Repository,
	publisher eventbus.EventBus,
) *TripService {
	return &TripService{
		tripRepo:    tripRepo,
		vehicleRepo: vehicleRepo,
		driverRepo:  driverRepo,
		publisher:   publisher,
	}
}

func (s *TripService) GetByID(ctx context.Context, id uuid.UUID) (trip.Trip, error) {
	if err := composables.CanUser(ctx, permissions.TripRead); err != nil {
		return nil, err
	}
	return s.tripRepo.GetByID(ctx, id)
}

func (s *TripService) GetPaginated(ctx context.Context, params *trip.FindParams) ([]trip.Trip, error) {
	if err := composables.CanUser(ctx, permissions.TripRead); err != nil {
		return nil, err
	}
	return s.tripRepo.GetPaginated(ctx, params)
}

func (s *TripService) Count(ctx context.Context, params *trip.FindParams) (int64, error) {
	if err := composables.CanUser(ctx, permissions.TripRead); err != nil {
		return 0, err
	}
	return s.tripRepo.Count(ctx, params)
}

func (s *TripService) GetByVehicle(ctx context.Context, vehicleID uuid.UUID) ([]trip.Trip, error) {
	if err := composables.CanUser(ctx, permissions.TripRead); err != nil {
		return nil, err
	}
	return s.tripRepo.GetByVehicle(ctx, vehicleID)
}

func (s *TripService) GetByDriver(ctx context.Context, driverID uuid.UUID) ([]trip.Trip, error) {
	if err := composables.CanUser(ctx, permissions.TripRead); err != nil {
		return nil, err
	}
	return s.tripRepo.GetByDriver(ctx, driverID)
}

func (s *TripService) GetActiveTrips(ctx context.Context, tenantID uuid.UUID) ([]trip.Trip, error) {
	if err := composables.CanUser(ctx, permissions.TripRead); err != nil {
		return nil, err
	}
	return s.tripRepo.GetActiveTrips(ctx, tenantID)
}

func (s *TripService) Create(ctx context.Context, t trip.Trip) (trip.Trip, error) {
	if err := composables.CanUser(ctx, permissions.TripCreate); err != nil {
		return nil, err
	}

	v, err := s.vehicleRepo.GetByID(ctx, t.VehicleID())
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}

	if v.Status() == enums.VehicleStatusMaintenance || v.Status() == enums.VehicleStatusOutOfService {
		return nil, fmt.Errorf("vehicle is not available for trips")
	}

	d, err := s.driverRepo.GetByID(ctx, t.DriverID())
	if err != nil {
		return nil, fmt.Errorf("failed to get driver: %w", err)
	}

	if d.Status() != enums.DriverStatusActive {
		return nil, fmt.Errorf("driver is not active")
	}

	if d.LicenseExpiry().Before(time.Now()) {
		return nil, fmt.Errorf("driver license has expired")
	}

	endTime := t.StartTime().Add(24 * time.Hour)
	if t.EndTime() != nil {
		endTime = *t.EndTime()
	}

	hasConflict, err := s.tripRepo.CheckConflict(ctx, t.VehicleID(), t.StartTime(), endTime, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check for conflicts: %w", err)
	}
	if hasConflict {
		return nil, fmt.Errorf("vehicle is already assigned to another trip during this time")
	}

	created, err := s.tripRepo.Create(ctx, t)
	if err != nil {
		return nil, fmt.Errorf("failed to create trip: %w", err)
	}

	updatedVehicle := v.UpdateStatus(enums.VehicleStatusInUse)
	if _, err := s.vehicleRepo.Update(ctx, updatedVehicle); err != nil {
		return nil, fmt.Errorf("failed to update vehicle status: %w", err)
	}

	event, err := trip.NewTripCreatedEvent(ctx, created)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	event.Result = created
	s.publisher.Publish(event)

	return created, nil
}

func (s *TripService) Complete(ctx context.Context, id uuid.UUID, endTime time.Time, endOdometer int) (trip.Trip, error) {
	if err := composables.CanUser(ctx, permissions.TripUpdate); err != nil {
		return nil, err
	}

	t, err := s.tripRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}

	if t.Status() != enums.TripStatusInProgress && t.Status() != enums.TripStatusScheduled {
		return nil, fmt.Errorf("trip cannot be completed in current status")
	}

	if endOdometer < t.StartOdometer() {
		return nil, fmt.Errorf("end odometer must be greater than start odometer")
	}

	completed := t.Complete(endTime, endOdometer)
	updated, err := s.tripRepo.Update(ctx, completed)
	if err != nil {
		return nil, fmt.Errorf("failed to update trip: %w", err)
	}

	v, err := s.vehicleRepo.GetByID(ctx, t.VehicleID())
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}

	updatedVehicle := v.UpdateStatus(enums.VehicleStatusAvailable).UpdateOdometer(endOdometer)
	if _, err := s.vehicleRepo.Update(ctx, updatedVehicle); err != nil {
		return nil, fmt.Errorf("failed to update vehicle: %w", err)
	}

	event, err := trip.NewTripCompletedEvent(ctx, updated)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	s.publisher.Publish(event)

	return updated, nil
}

func (s *TripService) Cancel(ctx context.Context, id uuid.UUID, reason string) (trip.Trip, error) {
	if err := composables.CanUser(ctx, permissions.TripUpdate); err != nil {
		return nil, err
	}

	t, err := s.tripRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip: %w", err)
	}

	if t.Status() == enums.TripStatusCompleted || t.Status() == enums.TripStatusCancelled {
		return nil, fmt.Errorf("trip cannot be cancelled in current status")
	}

	cancelled := t.Cancel(reason)
	updated, err := s.tripRepo.Update(ctx, cancelled)
	if err != nil {
		return nil, fmt.Errorf("failed to update trip: %w", err)
	}

	if t.Status() == enums.TripStatusInProgress {
		v, err := s.vehicleRepo.GetByID(ctx, t.VehicleID())
		if err != nil {
			return nil, fmt.Errorf("failed to get vehicle: %w", err)
		}

		updatedVehicle := v.UpdateStatus(enums.VehicleStatusAvailable)
		if _, err := s.vehicleRepo.Update(ctx, updatedVehicle); err != nil {
			return nil, fmt.Errorf("failed to update vehicle: %w", err)
		}
	}

	event, err := trip.NewTripCancelledEvent(ctx, updated, reason)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	s.publisher.Publish(event)

	return updated, nil
}

func (s *TripService) Update(ctx context.Context, t trip.Trip) (trip.Trip, error) {
	if err := composables.CanUser(ctx, permissions.TripUpdate); err != nil {
		return nil, err
	}

	updated, err := s.tripRepo.Update(ctx, t)
	if err != nil {
		return nil, fmt.Errorf("failed to update trip: %w", err)
	}

	return updated, nil
}

func (s *TripService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := composables.CanUser(ctx, permissions.TripDelete); err != nil {
		return err
	}

	if err := s.tripRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete trip: %w", err)
	}

	return nil
}

func (s *TripService) CalculateStatistics(t trip.Trip) map[string]interface{} {
	stats := make(map[string]interface{})

	if t.EndTime() != nil && t.EndOdometer() != nil {
		duration := t.EndTime().Sub(t.StartTime())
		distance := *t.EndOdometer() - t.StartOdometer()

		stats["duration_hours"] = duration.Hours()
		stats["distance_km"] = distance

		if duration.Hours() > 0 {
			stats["average_speed_kmh"] = float64(distance) / duration.Hours()
		}
	}

	return stats
}
