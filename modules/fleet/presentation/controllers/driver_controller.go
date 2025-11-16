package controllers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/iota-uz/iota-sdk/components/base/pagination"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/driver"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/controllers/dtos"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/mappers"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/templates/pages/drivers"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/viewmodels"
	"github.com/iota-uz/iota-sdk/modules/fleet/services"
	"github.com/iota-uz/iota-sdk/pkg/application"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/htmx"
	"github.com/iota-uz/iota-sdk/pkg/middleware"
	"github.com/iota-uz/iota-sdk/pkg/shared"
)

type DriverController struct {
	app           application.Application
	driverService *services.DriverService
	tripService   *services.TripService
	basePath      string
}

func NewDriverController(app application.Application) application.Controller {
	return &DriverController{
		app:           app,
		driverService: app.Service(services.DriverService{}).(*services.DriverService),
		tripService:   app.Service(services.TripService{}).(*services.TripService),
		basePath:      "/fleet/drivers",
	}
}

func (c *DriverController) Key() string {
	return c.basePath
}

func (c *DriverController) Register(r *mux.Router) {
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
	getRouter.HandleFunc("/{id}/edit", c.Edit).Methods(http.MethodGet)
	getRouter.HandleFunc("/{id}", c.Detail).Methods(http.MethodGet)

	setRouter := r.PathPrefix(c.basePath).Subrouter()
	setRouter.Use(commonMiddleware...)
	setRouter.Use(middleware.WithTransaction())
	setRouter.HandleFunc("", c.Create).Methods(http.MethodPost)
	setRouter.HandleFunc("/{id}", c.Update).Methods(http.MethodPost)
	setRouter.HandleFunc("/{id}/delete", c.Delete).Methods(http.MethodPost)
}

func (c *DriverController) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	paginationParams := composables.UsePaginated(r)
	filterDTO, err := composables.UseQuery(&dtos.DriverFilterDTO{
		Page:     paginationParams.Page,
		PageSize: paginationParams.Limit,
	}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := &driver.FindParams{
		TenantID: tenantID,
		Limit:    filterDTO.PageSize,
		Offset:   (filterDTO.Page - 1) * filterDTO.PageSize,
		SortBy:   driver.FieldCreatedAt,
		SortDesc: true,
	}

	if filterDTO.Status != "" {
		status, err := enums.ParseDriverStatus(filterDTO.Status)
		if err == nil {
			params.Status = &status
		}
	}

	if filterDTO.FirstName != "" || filterDTO.LastName != "" || filterDTO.LicenseNumber != "" {
		searchTerm := filterDTO.FirstName
		if searchTerm == "" {
			searchTerm = filterDTO.LastName
		}
		if searchTerm == "" {
			searchTerm = filterDTO.LicenseNumber
		}
		params.Search = &searchTerm
	}

	driverList, err := c.driverService.GetPaginated(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	total, err := c.driverService.Count(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	driverViewModels := make([]viewmodels.DriverListViewModel, len(driverList))
	for i, d := range driverList {
		driverViewModels[i] = mappers.DriverToListViewModel(d)
	}

	statusOptions := []struct {
		Value     string
		MessageID string
	}{
		{Value: "Active", MessageID: "Fleet.Enums.DriverStatus.Active"},
		{Value: "Inactive", MessageID: "Fleet.Enums.DriverStatus.Inactive"},
		{Value: "OnLeave", MessageID: "Fleet.Enums.DriverStatus.OnLeave"},
		{Value: "Terminated", MessageID: "Fleet.Enums.DriverStatus.Terminated"},
	}

	paginationState := pagination.New(r.URL.Path, filterDTO.Page, int(total), filterDTO.PageSize)

	props := &drivers.IndexPageProps{
		Drivers:         driverViewModels,
		PaginationState: paginationState,
		StatusOptions:   statusOptions,
	}

	if htmx.IsHxRequest(r) {
		w.WriteHeader(http.StatusOK)
		drivers.DriversTable(props).Render(ctx, w)
	} else {
		w.WriteHeader(http.StatusOK)
		drivers.Index(props).Render(ctx, w)
	}
}

func (c *DriverController) New(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	props := &drivers.CreatePageProps{
		Driver: &viewmodels.DriverFormViewModel{
			Errors: make(map[string]string),
		},
		PostPath: c.basePath,
	}

	w.WriteHeader(http.StatusOK)
	drivers.New(props).Render(ctx, w)
}

func (c *DriverController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dto, err := composables.UseForm(&dtos.DriverCreateDTO{}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	opts := []driver.DriverOption{
		driver.WithPhone(dto.Phone),
		driver.WithEmail(dto.Email),
	}
	if dto.UserID != nil {
		opts = append(opts, driver.WithUserID(*dto.UserID))
	}

	d := driver.NewDriver(
		uuid.New(),
		tenantID,
		dto.FirstName,
		dto.LastName,
		dto.LicenseNumber,
		dto.LicenseExpiry,
		opts...,
	)

	_, err = c.driverService.Create(ctx, d)
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

func (c *DriverController) Detail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid driver ID", http.StatusBadRequest)
		return
	}

	d, err := c.driverService.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	trips, err := c.tripService.GetByDriver(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	activeTrips := 0
	for _, t := range trips {
		if t.Status() == enums.TripStatusInProgress || t.Status() == enums.TripStatusScheduled {
			activeTrips++
		}
	}

	driverName := fmt.Sprintf("%s %s", d.FirstName(), d.LastName())
	driverVM := mappers.DriverToDetailViewModel(d, len(trips), activeTrips)
	tripVMs := make([]viewmodels.TripListViewModel, len(trips))
	for i, t := range trips {
		tripVMs[i] = mappers.TripToListViewModel(t, "", driverName)
	}

	props := &drivers.DetailPageProps{
		Driver: &driverVM,
		Trips:  tripVMs,
	}

	w.WriteHeader(http.StatusOK)
	drivers.Detail(props).Render(ctx, w)
}

func (c *DriverController) Edit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid driver ID", http.StatusBadRequest)
		return
	}

	d, err := c.driverService.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	statusOptions := []struct {
		Value     string
		MessageID string
	}{
		{Value: "Active", MessageID: "Fleet.Enums.DriverStatus.Active"},
		{Value: "Inactive", MessageID: "Fleet.Enums.DriverStatus.Inactive"},
		{Value: "OnLeave", MessageID: "Fleet.Enums.DriverStatus.OnLeave"},
		{Value: "Terminated", MessageID: "Fleet.Enums.DriverStatus.Terminated"},
	}

	driverVM := mappers.DriverToFormViewModel(d, make(map[string]string))
	props := &drivers.EditPageProps{
		Driver:        &driverVM,
		StatusOptions: statusOptions,
	}

	w.WriteHeader(http.StatusOK)
	drivers.Edit(props).Render(ctx, w)
}

func (c *DriverController) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid driver ID", http.StatusBadRequest)
		return
	}

	dto, err := composables.UseForm(&dtos.DriverUpdateDTO{}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existing, err := c.driverService.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	status, err := enums.ParseDriverStatus(dto.Status)
	if err != nil {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	opts := []driver.DriverOption{
		driver.WithPhone(dto.Phone),
		driver.WithEmail(dto.Email),
		driver.WithDriverStatus(status),
		driver.WithDriverTimestamps(existing.CreatedAt(), existing.UpdatedAt()),
	}
	if dto.UserID != nil {
		opts = append(opts, driver.WithUserID(*dto.UserID))
	}

	updated := driver.NewDriver(
		existing.ID(),
		existing.TenantID(),
		dto.FirstName,
		dto.LastName,
		dto.LicenseNumber,
		dto.LicenseExpiry,
		opts...,
	)

	_, err = c.driverService.Update(ctx, updated)
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

func (c *DriverController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid driver ID", http.StatusBadRequest)
		return
	}

	err = c.driverService.Delete(ctx, id)
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
