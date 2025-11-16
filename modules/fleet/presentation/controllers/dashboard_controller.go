package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/iota-uz/iota-sdk/components/charts"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/templates/pages/dashboard"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/viewmodels"
	"github.com/iota-uz/iota-sdk/modules/fleet/services"
	"github.com/iota-uz/iota-sdk/pkg/application"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/configuration"
	"github.com/iota-uz/iota-sdk/pkg/htmx"
	"github.com/iota-uz/iota-sdk/pkg/middleware"
)

type DashboardController struct {
	app              application.Application
	analyticsService *services.AnalyticsService
	basePath         string
}

func NewDashboardController(app application.Application) application.Controller {
	logger := configuration.Use().Logger()

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("Failed to create DashboardController: %v", r)
			panic(fmt.Sprintf("Failed to create DashboardController: %v", r))
		}
	}()

	analyticsService := app.Service(services.AnalyticsService{})
	if analyticsService == nil {
		logger.Error("AnalyticsService is nil - service not registered")
		panic("AnalyticsService is nil - ensure it is registered in module.go")
	}

	typedService, ok := analyticsService.(*services.AnalyticsService)
	if !ok {
		logger.Errorf("Failed to cast AnalyticsService to correct type, got %T", analyticsService)
		panic(fmt.Sprintf("Failed to cast AnalyticsService to correct type, got %T", analyticsService))
	}

	logger.Info("DashboardController created successfully")

	return &DashboardController{
		app:              app,
		analyticsService: typedService,
		basePath:         "/fleet/dashboard",
	}
}

func (c *DashboardController) Key() string {
	return c.basePath
}

func (c *DashboardController) Register(r *mux.Router) {
	commonMiddleware := []mux.MiddlewareFunc{
		middleware.Authorize(),
		middleware.RedirectNotAuthenticated(),
		middleware.ProvideUser(),
		middleware.ProvideDynamicLogo(c.app),
		middleware.ProvideLocalizer(c.app.Bundle()),
		middleware.NavItems(),
		middleware.WithPageContext(),
	}

	getRouter := r.PathPrefix(c.basePath).Subrouter()
	getRouter.Use(commonMiddleware...)
	getRouter.HandleFunc("", c.Index).Methods(http.MethodGet)
	getRouter.HandleFunc("/utilization", c.GetUtilizationData).Methods(http.MethodGet)
	getRouter.HandleFunc("/costs", c.GetCostTrends).Methods(http.MethodGet)
	getRouter.HandleFunc("/export", c.ExportReport).Methods(http.MethodGet)
}

func (c *DashboardController) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := configuration.Use().Logger()

	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		logger.Errorf("Failed to get tenant ID: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	logger.Infof("Dashboard Index called for tenant: %s", tenantID)

	stats, err := c.analyticsService.GetDashboardStats(ctx, tenantID)
	if err != nil {
		logger.Errorf("Failed to get dashboard stats for tenant %s: %v", tenantID, err)
		http.Error(w, fmt.Sprintf("Failed to get dashboard stats: %v", err), http.StatusInternalServerError)
		return
	}

	logger.Infof("Dashboard stats retrieved: %+v", stats)

	viewModel, err := c.buildDashboardViewModel(ctx, tenantID, stats)
	if err != nil {
		logger.Errorf("Failed to build dashboard view model for tenant %s: %v", tenantID, err)
		http.Error(w, fmt.Sprintf("Failed to build dashboard view model: %v", err), http.StatusInternalServerError)
		return
	}

	logger.Info("Dashboard view model built successfully")

	templ.Handler(dashboard.Index(viewModel)).ServeHTTP(w, r)
}

func (c *DashboardController) GetUtilizationData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time

	if startDateStr != "" {
		var parseErr error
		startDate, parseErr = time.Parse("2006-01-02", startDateStr)
		if parseErr != nil {
			http.Error(w, "Invalid start date", http.StatusBadRequest)
			return
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0)
	}

	if endDateStr != "" {
		var parseErr error
		endDate, parseErr = time.Parse("2006-01-02", endDateStr)
		if parseErr != nil {
			http.Error(w, "Invalid end date", http.StatusBadRequest)
			return
		}
	} else {
		endDate = time.Now()
	}

	report, err := c.analyticsService.GetUtilizationReport(ctx, tenantID, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
}

func (c *DashboardController) GetCostTrends(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time

	if startDateStr != "" {
		var parseErr error
		startDate, parseErr = time.Parse("2006-01-02", startDateStr)
		if parseErr != nil {
			http.Error(w, "Invalid start date", http.StatusBadRequest)
			return
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0)
	}

	if endDateStr != "" {
		var parseErr error
		endDate, parseErr = time.Parse("2006-01-02", endDateStr)
		if parseErr != nil {
			http.Error(w, "Invalid end date", http.StatusBadRequest)
			return
		}
	} else {
		endDate = time.Now()
	}

	trends, err := c.analyticsService.GetTrendData(ctx, tenantID, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trends)
}

func (c *DashboardController) ExportReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	reportType := r.URL.Query().Get("type")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate time.Time

	if startDateStr != "" {
		var parseErr error
		startDate, parseErr = time.Parse("2006-01-02", startDateStr)
		if parseErr != nil {
			http.Error(w, "Invalid start date", http.StatusBadRequest)
			return
		}
	} else {
		startDate = time.Now().AddDate(0, -1, 0)
	}

	if endDateStr != "" {
		var parseErr error
		endDate, parseErr = time.Parse("2006-01-02", endDateStr)
		if parseErr != nil {
			http.Error(w, "Invalid end date", http.StatusBadRequest)
			return
		}
	} else {
		endDate = time.Now()
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=fleet_report_%s.csv", time.Now().Format("2006-01-02")))

	switch reportType {
	case "utilization":
		report, err := c.analyticsService.GetUtilizationReport(ctx, tenantID, startDate, endDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Vehicle ID,Vehicle Name,Total Trips,Total Hours,Total Distance,Utilization %")
		for _, r := range report {
			fmt.Fprintf(w, "%s,%s,%d,%.2f,%d,%.2f\n",
				r.VehicleID,
				r.VehicleName,
				r.TotalTrips,
				r.TotalHours,
				r.TotalDistance,
				r.UtilizationPct,
			)
		}

	case "costs":
		analysis, err := c.analyticsService.GetCostAnalysis(ctx, tenantID, startDate, endDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Vehicle ID,Vehicle Name,Fuel Cost,Maintenance Cost,Total Cost,Cost Per Km")
		for _, a := range analysis {
			fmt.Fprintf(w, "%s,%s,%.2f,%.2f,%.2f,%.2f\n",
				a.VehicleID,
				a.VehicleName,
				a.FuelCost,
				a.MaintenanceCost,
				a.TotalCost,
				a.CostPerKm,
			)
		}

	default:
		http.Error(w, "Invalid report type", http.StatusBadRequest)
		return
	}

	if htmx.IsHxRequest(r) {
		htmx.SetTrigger(w, "reportExported", "")
	}
}

func (c *DashboardController) buildDashboardViewModel(ctx context.Context, tenantID uuid.UUID, stats *services.DashboardStats) (*viewmodels.DashboardViewModel, error) {
	logger := configuration.Use().Logger()
	now := time.Now()
	startDate := now.AddDate(0, -1, 0)
	endDate := now

	logger.Infof("Building dashboard view model for tenant %s from %s to %s", tenantID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	utilizationReport, err := c.analyticsService.GetUtilizationReport(ctx, tenantID, startDate, endDate)
	if err != nil {
		logger.Errorf("Failed to get utilization report for tenant %s: %v", tenantID, err)
		return nil, fmt.Errorf("failed to get utilization report: %w", err)
	}

	logger.Infof("Utilization report retrieved: %d records", len(utilizationReport))

	trendData, err := c.analyticsService.GetTrendData(ctx, tenantID, startDate, endDate)
	if err != nil {
		logger.Errorf("Failed to get trend data for tenant %s: %v", tenantID, err)
		return nil, fmt.Errorf("failed to get trend data: %w", err)
	}

	logger.Infof("Trend data retrieved: %d records", len(trendData))

	utilizationLabels := make([]string, 0)
	utilizationData := make([]float64, 0)
	tripCountData := make([]float64, 0)

	if utilizationReport != nil && len(utilizationReport) > 0 {
		for _, report := range utilizationReport {
			utilizationLabels = append(utilizationLabels, report.VehicleName)
			utilizationData = append(utilizationData, report.UtilizationPct)
			tripCountData = append(tripCountData, float64(report.TotalTrips))
		}
	} else {
		utilizationLabels = []string{"No Data"}
		utilizationData = []float64{0}
		tripCountData = []float64{0}
	}

	utilizationChart := charts.NewBarChart().
		WithSeries("Utilization %", utilizationData).
		WithSeries("Trip Count", tripCountData).
		WithCategories(utilizationLabels).
		WithHeight("320px").
		WithColors("#3b82f6", "#10b981").
		Build()

	trendLabels := make([]string, 0)
	fuelCostData := make([]float64, 0)
	maintenanceCostData := make([]float64, 0)

	if trendData != nil && len(trendData) > 0 {
		for _, trend := range trendData {
			trendLabels = append(trendLabels, trend.Date.Format("Jan 02"))
			fuelCostData = append(fuelCostData, trend.FuelCost)
			maintenanceCostData = append(maintenanceCostData, trend.MaintenanceCost)
		}
	} else {
		trendLabels = []string{"No Data"}
		fuelCostData = []float64{0}
		maintenanceCostData = []float64{0}
	}

	costTrendChart := charts.NewLineChart().
		WithSeries("Fuel Cost", fuelCostData).
		WithSeries("Maintenance Cost", maintenanceCostData).
		WithCategories(trendLabels).
		WithHeight("320px").
		WithColors("#8b5cf6", "#f59e0b").
		Build()

	return &viewmodels.DashboardViewModel{
		TotalVehicles:        stats.TotalVehicles,
		ActiveVehicles:       stats.ActiveVehicles,
		VehiclesInUse:        stats.VehiclesInUse,
		VehiclesMaintenance:  stats.VehiclesMaintenance,
		ActiveDrivers:        stats.ActiveDrivers,
		UpcomingMaintenance:  stats.UpcomingMaintenance,
		FuelCostMonth:        fmt.Sprintf("$%.2f", stats.FuelCostMonth),
		MaintenanceCostMonth: fmt.Sprintf("$%.2f", stats.MaintenanceCostMonth),
		TotalCostMonth:       fmt.Sprintf("$%.2f", stats.FuelCostMonth+stats.MaintenanceCostMonth),
		UtilizationChart: viewmodels.UtilizationChartViewModel{
			Options: utilizationChart,
		},
		CostTrendChart: viewmodels.CostTrendViewModel{
			Options: costTrendChart,
		},
	}, nil
}
