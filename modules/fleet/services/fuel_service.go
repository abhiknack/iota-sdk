package services

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/fuel_entry"
	"github.com/iota-uz/iota-sdk/modules/fleet/permissions"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/eventbus"
)

type FuelService struct {
	repo      fuel_entry.Repository
	publisher eventbus.EventBus
}

func NewFuelService(
	repo fuel_entry.Repository,
	publisher eventbus.EventBus,
) *FuelService {
	return &FuelService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *FuelService) GetByID(ctx context.Context, id uuid.UUID) (fuel_entry.FuelEntry, error) {
	if err := composables.CanUser(ctx, permissions.FuelEntryRead); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *FuelService) GetPaginated(ctx context.Context, params *fuel_entry.FindParams) ([]fuel_entry.FuelEntry, error) {
	if err := composables.CanUser(ctx, permissions.FuelEntryRead); err != nil {
		return nil, err
	}
	return s.repo.GetPaginated(ctx, params)
}

func (s *FuelService) Count(ctx context.Context, params *fuel_entry.FindParams) (int64, error) {
	if err := composables.CanUser(ctx, permissions.FuelEntryRead); err != nil {
		return 0, err
	}
	return s.repo.Count(ctx, params)
}

func (s *FuelService) GetByVehicle(ctx context.Context, vehicleID uuid.UUID) ([]fuel_entry.FuelEntry, error) {
	if err := composables.CanUser(ctx, permissions.FuelEntryRead); err != nil {
		return nil, err
	}
	return s.repo.GetByVehicle(ctx, vehicleID)
}

func (s *FuelService) GetByDriver(ctx context.Context, driverID uuid.UUID) ([]fuel_entry.FuelEntry, error) {
	if err := composables.CanUser(ctx, permissions.FuelEntryRead); err != nil {
		return nil, err
	}
	return s.repo.GetByDriver(ctx, driverID)
}

func (s *FuelService) Create(ctx context.Context, f fuel_entry.FuelEntry) (fuel_entry.FuelEntry, error) {
	if err := composables.CanUser(ctx, permissions.FuelEntryCreate); err != nil {
		return nil, err
	}

	created, err := s.repo.Create(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("failed to create fuel entry: %w", err)
	}

	event, err := fuel_entry.NewFuelEntryCreatedEvent(ctx, created)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	event.Result = created
	s.publisher.Publish(event)

	return created, nil
}

func (s *FuelService) Update(ctx context.Context, f fuel_entry.FuelEntry) (fuel_entry.FuelEntry, error) {
	if err := composables.CanUser(ctx, permissions.FuelEntryUpdate); err != nil {
		return nil, err
	}

	updated, err := s.repo.Update(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("failed to update fuel entry: %w", err)
	}

	event, err := fuel_entry.NewFuelEntryUpdatedEvent(ctx, updated)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}
	event.Result = updated
	s.publisher.Publish(event)

	return updated, nil
}

func (s *FuelService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := composables.CanUser(ctx, permissions.FuelEntryDelete); err != nil {
		return err
	}

	f, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get fuel entry: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete fuel entry: %w", err)
	}

	event, err := fuel_entry.NewFuelEntryDeletedEvent(ctx)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	event.Result = f
	s.publisher.Publish(event)

	return nil
}

func (s *FuelService) CalculateEfficiency(ctx context.Context, f fuel_entry.FuelEntry) (float64, error) {
	lastEntry, err := s.repo.GetLastEntry(ctx, f.VehicleID())
	if err != nil {
		return 0, nil
	}

	return f.CalculateEfficiency(lastEntry.Odometer()), nil
}

func (s *FuelService) DetectAnomaly(ctx context.Context, f fuel_entry.FuelEntry) (bool, error) {
	entries, err := s.repo.GetByVehicle(ctx, f.VehicleID())
	if err != nil {
		return false, fmt.Errorf("failed to get fuel entries: %w", err)
	}

	if len(entries) < 3 {
		return false, nil
	}

	var efficiencies []float64
	for i := 1; i < len(entries); i++ {
		eff := entries[i].CalculateEfficiency(entries[i-1].Odometer())
		if eff > 0 {
			efficiencies = append(efficiencies, eff)
		}
	}

	if len(efficiencies) == 0 {
		return false, nil
	}

	avgEfficiency := average(efficiencies)
	currentEfficiency, err := s.CalculateEfficiency(ctx, f)
	if err != nil {
		return false, err
	}

	if currentEfficiency == 0 {
		return false, nil
	}

	deviation := math.Abs(currentEfficiency-avgEfficiency) / avgEfficiency
	return deviation > 0.20, nil
}

func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}
