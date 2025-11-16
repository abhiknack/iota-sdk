package controllers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/iota-uz/iota-sdk/components/base/pagination"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/vehicle"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/controllers/dtos"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/mappers"
	vehiclestpl "github.com/iota-uz/iota-sdk/modules/fleet/presentation/templates/pages/vehicles"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/viewmodels"
	"github.com/iota-uz/iota-sdk/modules/fleet/services"
	"github.com/iota-uz/iota-sdk/pkg/application"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/htmx"
	"github.com/iota-uz/iota-sdk/pkg/middleware"
	"github.com/iota-uz/iota-sdk/pkg/shared"
)

type VehicleController struct {
	app            application.Application
	vehicleService *services.VehicleService
	basePath       string
}

func NewVehicleController(app application.Application) application.Controller {
	return &VehicleController{
		app:            app,
		vehicleService: app.Service(services.VehicleService{}).(*services.VehicleService),
		basePath:       "/fleet/vehicles",
	}
}

func (c *VehicleController) Key() string {
	return c.basePath
}

func (c *VehicleController) Register(r *mux.Router) {
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
	getRouter.HandleFunc("/{id}/detail", c.Detail).Methods(http.MethodGet)

	setRouter := r.PathPrefix(c.basePath).Subrouter()
	setRouter.Use(commonMiddleware...)
	setRouter.Use(middleware.WithTransaction())
	setRouter.HandleFunc("", c.Create).Methods(http.MethodPost)
	setRouter.HandleFunc("/{id}", c.Update).Methods(http.MethodPost)
	setRouter.HandleFunc("/{id}/delete", c.Delete).Methods(http.MethodPost)
	setRouter.HandleFunc("/{id}/status", c.ChangeStatus).Methods(http.MethodPost)
}

func (c *VehicleController) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	paginationParams := composables.UsePaginated(r)
	filterDTO, err := composables.UseQuery(&dtos.VehicleFilterDTO{
		Page:     paginationParams.Page,
		PageSize: paginationParams.Limit,
	}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := &vehicle.FindParams{
		TenantID: tenantID,
		Limit:    filterDTO.PageSize,
		Offset:   (filterDTO.Page - 1) * filterDTO.PageSize,
		SortBy:   vehicle.FieldCreatedAt,
		SortDesc: true,
	}

	if filterDTO.Status != "" {
		status, err := enums.ParseVehicleStatus(filterDTO.Status)
		if err == nil {
			params.Status = &status
		}
	}
	if filterDTO.Make != "" {
		params.Make = &filterDTO.Make
	}
	if filterDTO.Model != "" {
		params.Model = &filterDTO.Model
	}
	if filterDTO.Year > 0 {
		params.Year = &filterDTO.Year
	}

	vehiclesList, err := c.vehicleService.GetPaginated(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	total, err := c.vehicleService.Count(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vehicleVMs := make([]viewmodels.VehicleListViewModel, len(vehiclesList))
	for i, v := range vehiclesList {
		vehicleVMs[i] = mappers.VehicleToListViewModel(v)
	}

	paginationState := pagination.New(r.URL.Path, filterDTO.Page, int(total), filterDTO.PageSize)

	statusOptions := []struct {
		Value     string
		MessageID string
	}{
		{Value: "Available", MessageID: "Fleet.Enums.VehicleStatus.Available"},
		{Value: "InUse", MessageID: "Fleet.Enums.VehicleStatus.InUse"},
		{Value: "Maintenance", MessageID: "Fleet.Enums.VehicleStatus.Maintenance"},
		{Value: "OutOfService", MessageID: "Fleet.Enums.VehicleStatus.OutOfService"},
		{Value: "Retired", MessageID: "Fleet.Enums.VehicleStatus.Retired"},
	}

	props := &vehiclestpl.IndexPageProps{
		Vehicles:        vehicleVMs,
		PaginationState: paginationState,
		StatusOptions:   statusOptions,
	}

	if htmx.IsHxRequest(r) {
		w.WriteHeader(http.StatusOK)
		vehiclestpl.VehiclesTable(props).Render(ctx, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	vehiclestpl.Index(props).Render(ctx, w)
}

func (c *VehicleController) New(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	props := &vehiclestpl.CreatePageProps{
		Vehicle: &viewmodels.VehicleFormViewModel{
			Errors: make(map[string]string),
		},
		PostPath: c.basePath,
	}

	w.WriteHeader(http.StatusOK)
	vehiclestpl.New(props).Render(ctx, w)
}

func (c *VehicleController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dto, err := composables.UseForm(&dtos.VehicleCreateDTO{}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	v := vehicle.NewVehicle(
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

	created, err := c.vehicleService.Create(ctx, v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if htmx.IsHxRequest(r) {
		htmx.Redirect(w, c.basePath)
	} else {
		shared.Redirect(w, r, c.basePath)
	}
	fmt.Fprintf(w, "Vehicle created: %s", created.ID())
}

func (c *VehicleController) Edit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	v, err := c.vehicleService.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	vehicleVM := mappers.VehicleToFormViewModel(v, make(map[string]string))

	props := &vehiclestpl.EditPageProps{
		Vehicle: &vehicleVM,
	}

	w.WriteHeader(http.StatusOK)
	vehiclestpl.Edit(props).Render(ctx, w)
}

func (c *VehicleController) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	dto, err := composables.UseForm(&dtos.VehicleUpdateDTO{}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existing, err := c.vehicleService.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	updated := existing.UpdateDetails(dto.Make, dto.Model, dto.Year)
	updated = updated.UpdateOdometer(dto.CurrentOdometer)

	updatedVehicle := vehicle.NewVehicle(
		existing.ID(),
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
		vehicle.WithTimestamps(existing.CreatedAt(), existing.UpdatedAt()),
	)

	_, err = c.vehicleService.Update(ctx, updatedVehicle)
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

func (c *VehicleController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	err = c.vehicleService.Delete(ctx, id)
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

func (c *VehicleController) Detail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	v, err := c.vehicleService.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	vehicleVM := mappers.VehicleToDetailViewModel(v)

	currentStatus := v.Status()
	statusOptions := []struct {
		Value     string
		MessageID string
		Disabled  bool
	}{
		{
			Value:     "Available",
			MessageID: "Fleet.Enums.VehicleStatus.Available",
			Disabled:  !currentStatus.CanTransitionTo(enums.VehicleStatusAvailable),
		},
		{
			Value:     "InUse",
			MessageID: "Fleet.Enums.VehicleStatus.InUse",
			Disabled:  !currentStatus.CanTransitionTo(enums.VehicleStatusInUse),
		},
		{
			Value:     "Maintenance",
			MessageID: "Fleet.Enums.VehicleStatus.Maintenance",
			Disabled:  !currentStatus.CanTransitionTo(enums.VehicleStatusMaintenance),
		},
		{
			Value:     "OutOfService",
			MessageID: "Fleet.Enums.VehicleStatus.OutOfService",
			Disabled:  !currentStatus.CanTransitionTo(enums.VehicleStatusOutOfService),
		},
		{
			Value:     "Retired",
			MessageID: "Fleet.Enums.VehicleStatus.Retired",
			Disabled:  !currentStatus.CanTransitionTo(enums.VehicleStatusRetired),
		},
	}

	props := &vehiclestpl.DetailPageProps{
		Vehicle:       &vehicleVM,
		StatusOptions: statusOptions,
	}

	w.WriteHeader(http.StatusOK)
	vehiclestpl.Detail(props).Render(ctx, w)
}

func (c *VehicleController) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	statusStr := r.FormValue("Status")
	status, err := enums.ParseVehicleStatus(statusStr)
	if err != nil {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	_, err = c.vehicleService.UpdateStatus(ctx, id, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if htmx.IsHxRequest(r) {
		htmx.Redirect(w, fmt.Sprintf("%s/%s/detail", c.basePath, id))
	} else {
		shared.Redirect(w, r, fmt.Sprintf("%s/%s/detail", c.basePath, id))
	}
}
