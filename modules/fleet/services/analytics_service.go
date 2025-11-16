package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/driver"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/fuel_entry"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/maintenance"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/trip"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/vehicle"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
	"github.com/iota-uz/iota-sdk/modules/fleet/permissions"
	"github.com/iota-uz/iota-sdk/pkg/composables"
)

type DashboardStats struct {
	TotalVehicles        int
	ActiveVehicles       int
	VehiclesInUse        int
	VehiclesMaintenance  int
	ActiveDrivers        int
	UpcomingMaintenance  int
	FuelCostMonth        float64
	MaintenanceCostMonth float64
}

type UtilizationReport struct {
	VehicleID      uuid.UUID
	VehicleName    string
	TotalTrips     int
	TotalHours     float64
	TotalDistance  int
	UtilizationPct float64
}

type CostAnalysis struct {
	VehicleID       uuid.UUID
	VehicleName     string
	FuelCost        float64
	MaintenanceCost float64
	TotalCost       float64
	CostPerKm       float64
}

type TrendData struct {
	Date            time.Time
	FuelCost        float64
	MaintenanceCost float64
	TripCount       int
	Distance        int
}

type AnalyticsService struct {
	vehicleRepo     vehicle.Repository
	driverRepo      driver.Repository
	tripRepo        trip.Repository
	maintenanceRepo maintenance.Repository
	fuelRepo        fuel_entry.Repository
}

func NewAnalyticsService(
	vehicleRepo vehicle.Repository,
	driverRepo driver.Repository,
	tripRepo trip.Repository,
	maintenanceRepo maintenance.Repository,
	fuelRepo fuel_entry.Repository,
) *AnalyticsService {
	return &AnalyticsService{
		vehicleRepo:     vehicleRepo,
		driverRepo:      driverRepo,
		tripRepo:        tripRepo,
		maintenanceRepo: maintenanceRepo,
		fuelRepo:        fuelRepo,
	}
}

func (s *AnalyticsService) GetDashboardStats(ctx context.Context, tenantID uuid.UUID) (*DashboardStats, error) {
	if err := composables.CanUser(ctx, permissions.VehicleRead); err != nil {
		return nil, fmt.Errorf("permission check failed: %w", err)
	}

	vehicles, err := s.vehicleRepo.GetPaginated(ctx, &vehicle.FindParams{
		TenantID: tenantID,
		Limit:    10000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicles: %w", err)
	}

	if vehicles == nil {
		vehicles = []vehicle.Vehicle{}
	}

	stats := &DashboardStats{
		TotalVehicles: len(vehicles),
	}

	for _, v := range vehicles {
		switch v.Status() {
		case enums.VehicleStatusAvailable:
			stats.ActiveVehicles++
		case enums.VehicleStatusInUse:
			stats.VehiclesInUse++
		case enums.VehicleStatusMaintenance:
			stats.VehiclesMaintenance++
		}
	}

	drivers, err := s.driverRepo.GetPaginated(ctx, &driver.FindParams{
		TenantID: tenantID,
		Limit:    10000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get drivers: %w", err)
	}

	if drivers == nil {
		drivers = []driver.Driver{}
	}

	for _, d := range drivers {
		if d.Status() == enums.DriverStatusActive {
			stats.ActiveDrivers++
		}
	}

	dueMaintenance, err := s.maintenanceRepo.GetDueMaintenance(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get due maintenance: %w", err)
	}

	if dueMaintenance == nil {
		dueMaintenance = []maintenance.Maintenance{}
	}
	stats.UpcomingMaintenance = len(dueMaintenance)

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	startOfMonthStr := startOfMonth.Format("2006-01-02")
	endOfMonthStr := endOfMonth.Format("2006-01-02")

	fuelEntries, err := s.fuelRepo.GetPaginated(ctx, &fuel_entry.FindParams{
		TenantID: tenantID,
		DateFrom: &startOfMonthStr,
		DateTo:   &endOfMonthStr,
		Limit:    10000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get fuel entries: %w", err)
	}

	if fuelEntries == nil {
		fuelEntries = []fuel_entry.FuelEntry{}
	}

	for _, f := range fuelEntries {
		stats.FuelCostMonth += f.Cost()
	}

	maintenanceRecords, err := s.maintenanceRepo.GetPaginated(ctx, &maintenance.FindParams{
		TenantID:        tenantID,
		ServiceDateFrom: &startOfMonthStr,
		ServiceDateTo:   &endOfMonthStr,
		Limit:           10000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get maintenance records: %w", err)
	}

	if maintenanceRecords == nil {
		maintenanceRecords = []maintenance.Maintenance{}
	}

	for _, m := range maintenanceRecords {
		stats.MaintenanceCostMonth += m.Cost()
	}

	return stats, nil
}

func (s *AnalyticsService) GetUtilizationReport(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) ([]UtilizationReport, error) {
	if err := composables.CanUser(ctx, permissions.VehicleRead); err != nil {
		return nil, err
	}

	vehicles, err := s.vehicleRepo.GetPaginated(ctx, &vehicle.FindParams{
		TenantID: tenantID,
		Limit:    10000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicles: %w", err)
	}

	if vehicles == nil {
		vehicles = []vehicle.Vehicle{}
	}

	reports := make([]UtilizationReport, 0)
	totalPeriodHours := endDate.Sub(startDate).Hours()

	for _, v := range vehicles {
		trips, err := s.tripRepo.GetByVehicle(ctx, v.ID())
		if err != nil {
			continue
		}

		if trips == nil {
			trips = []trip.Trip{}
		}

		report := UtilizationReport{
			VehicleID:   v.ID(),
			VehicleName: fmt.Sprintf("%s %s", v.Make(), v.Model()),
		}

		var totalHours float64
		for _, t := range trips {
			if t.StartTime().Before(startDate) || t.StartTime().After(endDate) {
				continue
			}

			report.TotalTrips++

			if t.EndTime() != nil {
				duration := t.EndTime().Sub(t.StartTime()).Hours()
				totalHours += duration

				if t.EndOdometer() != nil {
					report.TotalDistance += *t.EndOdometer() - t.StartOdometer()
				}
			}
		}

		report.TotalHours = totalHours
		if totalPeriodHours > 0 {
			report.UtilizationPct = (totalHours / totalPeriodHours) * 100
		}

		reports = append(reports, report)
	}

	return reports, nil
}

func (s *AnalyticsService) GetCostAnalysis(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) ([]CostAnalysis, error) {
	if err := composables.CanUser(ctx, permissions.VehicleRead); err != nil {
		return nil, err
	}

	vehicles, err := s.vehicleRepo.GetPaginated(ctx, &vehicle.FindParams{
		TenantID: tenantID,
		Limit:    10000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicles: %w", err)
	}

	if vehicles == nil {
		vehicles = []vehicle.Vehicle{}
	}

	analyses := make([]CostAnalysis, 0)

	for _, v := range vehicles {
		analysis := CostAnalysis{
			VehicleID:   v.ID(),
			VehicleName: fmt.Sprintf("%s %s", v.Make(), v.Model()),
		}

		fuelEntries, err := s.fuelRepo.GetByVehicle(ctx, v.ID())
		if err != nil {
			continue
		}

		if fuelEntries == nil {
			fuelEntries = []fuel_entry.FuelEntry{}
		}

		for _, f := range fuelEntries {
			if f.Date().After(startDate) && f.Date().Before(endDate) {
				analysis.FuelCost += f.Cost()
			}
		}

		maintenanceRecords, err := s.maintenanceRepo.GetByVehicle(ctx, v.ID())
		if err != nil {
			continue
		}

		if maintenanceRecords == nil {
			maintenanceRecords = []maintenance.Maintenance{}
		}

		for _, m := range maintenanceRecords {
			if m.ServiceDate().After(startDate) && m.ServiceDate().Before(endDate) {
				analysis.MaintenanceCost += m.Cost()
			}
		}

		analysis.TotalCost = analysis.FuelCost + analysis.MaintenanceCost

		trips, err := s.tripRepo.GetByVehicle(ctx, v.ID())
		if err != nil {
			continue
		}

		if trips == nil {
			trips = []trip.Trip{}
		}

		totalDistance := 0
		for _, t := range trips {
			if t.StartTime().After(startDate) && t.StartTime().Before(endDate) {
				if t.EndOdometer() != nil {
					totalDistance += *t.EndOdometer() - t.StartOdometer()
				}
			}
		}

		if totalDistance > 0 {
			analysis.CostPerKm = analysis.TotalCost / float64(totalDistance)
		}

		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

func (s *AnalyticsService) GetTrendData(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) ([]TrendData, error) {
	if err := composables.CanUser(ctx, permissions.VehicleRead); err != nil {
		return nil, err
	}

	trends := make([]TrendData, 0)
	current := startDate

	for current.Before(endDate) {
		nextDay := current.AddDate(0, 0, 1)

		trend := TrendData{
			Date: current,
		}

		currentStr := current.Format("2006-01-02")
		nextDayStr := nextDay.Format("2006-01-02")

		fuelEntries, err := s.fuelRepo.GetPaginated(ctx, &fuel_entry.FindParams{
			TenantID: tenantID,
			DateFrom: &currentStr,
			DateTo:   &nextDayStr,
			Limit:    10000,
		})
		if err == nil {
			if fuelEntries == nil {
				fuelEntries = []fuel_entry.FuelEntry{}
			}
			for _, f := range fuelEntries {
				trend.FuelCost += f.Cost()
			}
		}

		maintenanceRecords, err := s.maintenanceRepo.GetPaginated(ctx, &maintenance.FindParams{
			TenantID:        tenantID,
			ServiceDateFrom: &currentStr,
			ServiceDateTo:   &nextDayStr,
			Limit:           10000,
		})
		if err == nil {
			if maintenanceRecords == nil {
				maintenanceRecords = []maintenance.Maintenance{}
			}
			for _, m := range maintenanceRecords {
				trend.MaintenanceCost += m.Cost()
			}
		}

		trips, err := s.tripRepo.GetPaginated(ctx, &trip.FindParams{
			TenantID:      tenantID,
			StartTimeFrom: &current,
			StartTimeTo:   &nextDay,
			Limit:         10000,
		})
		if err == nil {
			if trips == nil {
				trips = []trip.Trip{}
			}
			trend.TripCount = len(trips)
			for _, t := range trips {
				if t.EndOdometer() != nil {
					trend.Distance += *t.EndOdometer() - t.StartOdometer()
				}
			}
		}

		trends = append(trends, trend)
		current = nextDay
	}

	return trends, nil
}
