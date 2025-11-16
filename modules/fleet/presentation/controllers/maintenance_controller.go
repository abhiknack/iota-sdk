package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/iota-uz/iota-sdk/components/base/pagination"
	maintenancedom "github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/maintenance"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/vehicle"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/controllers/dtos"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/mappers"
	maintenancetpl "github.com/iota-uz/iota-sdk/modules/fleet/presentation/templates/pages/maintenance"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/viewmodels"
	"github.com/iota-uz/iota-sdk/modules/fleet/services"
	"github.com/iota-uz/iota-sdk/pkg/application"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/htmx"
	"github.com/iota-uz/iota-sdk/pkg/middleware"
	"github.com/iota-uz/iota-sdk/pkg/shared"
)

type MaintenanceController struct {
	app                application.Application
	maintenanceService *services.MaintenanceService
	basePath           string
}

func NewMaintenanceController(app application.Application) application.Controller {
	return &MaintenanceController{
		app:                app,
		maintenanceService: app.Service(services.MaintenanceService{}).(*services.MaintenanceService),
		basePath:           "/fleet/maintenance",
	}
}

func (c *MaintenanceController) Key() string {
	return c.basePath
}

func (c *MaintenanceController) Register(r *mux.Router) {
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
	getRouter.HandleFunc("", c.List).Methods(http.MethodGet)
	getRouter.HandleFunc("/new", c.New).Methods(http.MethodGet)
	getRouter.HandleFunc("/{id}", c.Edit).Methods(http.MethodGet)

	setRouter := r.PathPrefix(c.basePath).Subrouter()
	setRouter.Use(commonMiddleware...)
	setRouter.Use(middleware.WithTransaction())
	setRouter.HandleFunc("", c.Create).Methods(http.MethodPost)
	setRouter.HandleFunc("/{id}", c.Update).Methods(http.MethodPost)
	setRouter.HandleFunc("/{id}/delete", c.Delete).Methods(http.MethodPost)
}

func (c *MaintenanceController) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	paginationParams := composables.UsePaginated(r)
	filterDTO, err := composables.UseQuery(&dtos.MaintenanceFilterDTO{
		Page:     paginationParams.Page,
		PageSize: paginationParams.Limit,
	}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := &maintenancedom.FindParams{
		TenantID: tenantID,
		Limit:    filterDTO.PageSize,
		Offset:   (filterDTO.Page - 1) * filterDTO.PageSize,
		SortBy:   maintenancedom.FieldServiceDate,
		SortDesc: true,
	}

	if filterDTO.VehicleID != uuid.Nil {
		params.VehicleID = &filterDTO.VehicleID
	}
	if filterDTO.ServiceType != "" {
		serviceType, err := enums.ParseServiceType(filterDTO.ServiceType)
		if err == nil {
			params.ServiceType = &serviceType
		}
	}
	if !filterDTO.StartDate.IsZero() {
		startDateStr := filterDTO.StartDate.Format("2006-01-02")
		params.ServiceDateFrom = &startDateStr
	}
	if !filterDTO.EndDate.IsZero() {
		endDateStr := filterDTO.EndDate.Format("2006-01-02")
		params.ServiceDateTo = &endDateStr
	}

	maintenanceRecords, err := c.maintenanceService.GetPaginated(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	total, err := c.maintenanceService.Count(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	viewModels, err := c.buildMaintenanceListViewModels(ctx, maintenanceRecords)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	props := c.buildMaintenanceListProps(viewModels, total, filterDTO.Page, filterDTO.PageSize, r.URL.Path)

	if htmx.IsHxRequest(r) {
		w.WriteHeader(http.StatusOK)
		maintenancetpl.MaintenanceTable(props).Render(ctx, w)
	} else {
		w.WriteHeader(http.StatusOK)
		maintenancetpl.Index(props).Render(ctx, w)
	}
}

func (c *MaintenanceController) New(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vehicleOptions, err := c.getVehicleOptions(ctx, tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	props := c.buildMaintenanceCreateProps(vehicleOptions)
	w.WriteHeader(http.StatusOK)
	maintenancetpl.New(props).Render(ctx, w)
}

func (c *MaintenanceController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dto, err := composables.UseForm(&dtos.MaintenanceCreateDTO{}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	serviceType, err := enums.ParseServiceType(dto.ServiceType)
	if err != nil {
		http.Error(w, "Invalid service type", http.StatusBadRequest)
		return
	}

	opts := []maintenancedom.MaintenanceOption{}
	if dto.NextServiceDue != nil {
		opts = append(opts, maintenancedom.WithNextServiceDue(*dto.NextServiceDue))
	}
	if dto.NextServiceOdometer != nil {
		opts = append(opts, maintenancedom.WithNextServiceOdometer(*dto.NextServiceOdometer))
	}

	m := maintenancedom.NewMaintenance(
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
	)

	created, err := c.maintenanceService.Create(ctx, m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if htmx.IsHxRequest(r) {
		htmx.Redirect(w, c.basePath)
	} else {
		shared.Redirect(w, r, c.basePath)
	}
	fmt.Fprintf(w, "Maintenance created: %s", created.ID())
}

func (c *MaintenanceController) Edit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid maintenance ID", http.StatusBadRequest)
		return
	}

	m, err := c.maintenanceService.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	vehicleOptions, err := c.getVehicleOptions(ctx, tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	props := c.buildMaintenanceEditProps(m, vehicleOptions)
	w.WriteHeader(http.StatusOK)
	maintenancetpl.Edit(props).Render(ctx, w)
}

func (c *MaintenanceController) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid maintenance ID", http.StatusBadRequest)
		return
	}

	dto, err := composables.UseForm(&dtos.MaintenanceUpdateDTO{}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existing, err := c.maintenanceService.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	serviceType, err := enums.ParseServiceType(dto.ServiceType)
	if err != nil {
		http.Error(w, "Invalid service type", http.StatusBadRequest)
		return
	}

	opts := []maintenancedom.MaintenanceOption{
		maintenancedom.WithTimestamps(existing.CreatedAt(), existing.UpdatedAt()),
	}
	if dto.NextServiceDue != nil {
		opts = append(opts, maintenancedom.WithNextServiceDue(*dto.NextServiceDue))
	}
	if dto.NextServiceOdometer != nil {
		opts = append(opts, maintenancedom.WithNextServiceOdometer(*dto.NextServiceOdometer))
	}

	updated := maintenancedom.NewMaintenance(
		existing.ID(),
		existing.TenantID(),
		dto.VehicleID,
		serviceType,
		dto.ServiceDate,
		dto.Odometer,
		dto.Cost,
		dto.ServiceProvider,
		dto.Description,
		opts...,
	)

	_, err = c.maintenanceService.Update(ctx, updated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if htmx.IsHxRequest(r) {
		htmx.Redirect(w, c.basePath)
	} else {
		shared.Redirect(w, r, c.basePath)
	}
}

func (c *MaintenanceController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid maintenance ID", http.StatusBadRequest)
		return
	}

	err = c.maintenanceService.Delete(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if htmx.IsHxRequest(r) {
		htmx.Redirect(w, c.basePath)
	} else {
		shared.Redirect(w, r, c.basePath)
	}
}

func (c *MaintenanceController) buildMaintenanceListViewModels(ctx context.Context, maintenanceRecords []maintenancedom.Maintenance) ([]viewmodels.MaintenanceListViewModel, error) {
	vehicleService := c.app.Service(services.VehicleService{}).(*services.VehicleService)
	viewModels := make([]viewmodels.MaintenanceListViewModel, 0, len(maintenanceRecords))

	for _, m := range maintenanceRecords {
		v, err := vehicleService.GetByID(ctx, m.VehicleID())
		if err != nil {
			return nil, fmt.Errorf("failed to get vehicle: %w", err)
		}

		vehicleName := fmt.Sprintf("%s %s (%d)", v.Make(), v.Model(), v.Year())
		vm := mappers.MaintenanceToListViewModel(m, vehicleName)
		viewModels = append(viewModels, vm)
	}

	return viewModels, nil
}

func (c *MaintenanceController) buildMaintenanceListProps(viewModels []viewmodels.MaintenanceListViewModel, total int64, page, pageSize int, urlPath string) *maintenancetpl.IndexPageProps {
	serviceTypeOptions := []struct {
		Value     string
		MessageID string
	}{
		{Value: enums.ServiceTypeOilChange.String(), MessageID: "Fleet.Enums.ServiceType.OilChange"},
		{Value: enums.ServiceTypeTireRotation.String(), MessageID: "Fleet.Enums.ServiceType.TireRotation"},
		{Value: enums.ServiceTypeBrakeService.String(), MessageID: "Fleet.Enums.ServiceType.BrakeService"},
		{Value: enums.ServiceTypeInspection.String(), MessageID: "Fleet.Enums.ServiceType.Inspection"},
		{Value: enums.ServiceTypeRepair.String(), MessageID: "Fleet.Enums.ServiceType.Repair"},
		{Value: enums.ServiceTypeOther.String(), MessageID: "Fleet.Enums.ServiceType.Other"},
	}

	paginationState := pagination.New(urlPath, page, int(total), pageSize)

	return &maintenancetpl.IndexPageProps{
		MaintenanceRecords: viewModels,
		PaginationState:    paginationState,
		ServiceTypeOptions: serviceTypeOptions,
	}
}

func (c *MaintenanceController) buildMaintenanceCreateProps(vehicleOptions []struct{ Value, Label string }) *maintenancetpl.CreatePageProps {
	serviceTypeOptions := []struct {
		Value     string
		MessageID string
	}{
		{Value: enums.ServiceTypeOilChange.String(), MessageID: "Fleet.Enums.ServiceType.OilChange"},
		{Value: enums.ServiceTypeTireRotation.String(), MessageID: "Fleet.Enums.ServiceType.TireRotation"},
		{Value: enums.ServiceTypeBrakeService.String(), MessageID: "Fleet.Enums.ServiceType.BrakeService"},
		{Value: enums.ServiceTypeInspection.String(), MessageID: "Fleet.Enums.ServiceType.Inspection"},
		{Value: enums.ServiceTypeRepair.String(), MessageID: "Fleet.Enums.ServiceType.Repair"},
		{Value: enums.ServiceTypeOther.String(), MessageID: "Fleet.Enums.ServiceType.Other"},
	}

	return &maintenancetpl.CreatePageProps{
		Maintenance: &viewmodels.MaintenanceFormViewModel{
			Errors: make(map[string]string),
		},
		PostPath:           c.basePath,
		VehicleOptions:     vehicleOptions,
		ServiceTypeOptions: serviceTypeOptions,
	}
}

func (c *MaintenanceController) buildMaintenanceEditProps(m maintenancedom.Maintenance, vehicleOptions []struct{ Value, Label string }) *maintenancetpl.EditPageProps {
	serviceTypeOptions := []struct {
		Value     string
		MessageID string
	}{
		{Value: enums.ServiceTypeOilChange.String(), MessageID: "Fleet.Enums.ServiceType.OilChange"},
		{Value: enums.ServiceTypeTireRotation.String(), MessageID: "Fleet.Enums.ServiceType.TireRotation"},
		{Value: enums.ServiceTypeBrakeService.String(), MessageID: "Fleet.Enums.ServiceType.BrakeService"},
		{Value: enums.ServiceTypeInspection.String(), MessageID: "Fleet.Enums.ServiceType.Inspection"},
		{Value: enums.ServiceTypeRepair.String(), MessageID: "Fleet.Enums.ServiceType.Repair"},
		{Value: enums.ServiceTypeOther.String(), MessageID: "Fleet.Enums.ServiceType.Other"},
	}

	vm := mappers.MaintenanceToFormViewModel(m, make(map[string]string))
	return &maintenancetpl.EditPageProps{
		Maintenance:        &vm,
		VehicleOptions:     vehicleOptions,
		ServiceTypeOptions: serviceTypeOptions,
	}
}

func (c *MaintenanceController) getVehicleOptions(ctx context.Context, tenantID uuid.UUID) ([]struct{ Value, Label string }, error) {
	vehicleService := c.app.Service(services.VehicleService{}).(*services.VehicleService)
	vehicles, err := vehicleService.GetPaginated(ctx, &vehicle.FindParams{
		TenantID: tenantID,
		Limit:    1000,
		Offset:   0,
	})
	if err != nil {
		return nil, err
	}

	options := make([]struct{ Value, Label string }, 0, len(vehicles))
	for _, v := range vehicles {
		options = append(options, struct{ Value, Label string }{
			Value: v.ID().String(),
			Label: fmt.Sprintf("%s %s (%d) - %s", v.Make(), v.Model(), v.Year(), v.LicensePlate()),
		})
	}

	return options, nil
}
