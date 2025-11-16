package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/iota-uz/iota-sdk/components/base/pagination"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/driver"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/fuel_entry"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/vehicle"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/controllers/dtos"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/mappers"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/templates/pages/fuel"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/viewmodels"
	"github.com/iota-uz/iota-sdk/modules/fleet/services"
	"github.com/iota-uz/iota-sdk/pkg/application"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/htmx"
	"github.com/iota-uz/iota-sdk/pkg/middleware"
	"github.com/iota-uz/iota-sdk/pkg/shared"
)

type FuelController struct {
	app            application.Application
	fuelService    *services.FuelService
	vehicleService *services.VehicleService
	driverService  *services.DriverService
	basePath       string
}

func NewFuelController(app application.Application) application.Controller {
	return &FuelController{
		app:            app,
		fuelService:    app.Service(services.FuelService{}).(*services.FuelService),
		vehicleService: app.Service(services.VehicleService{}).(*services.VehicleService),
		driverService:  app.Service(services.DriverService{}).(*services.DriverService),
		basePath:       "/fleet/fuel",
	}
}

func (c *FuelController) Key() string {
	return c.basePath
}

func (c *FuelController) Register(r *mux.Router) {
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

func (c *FuelController) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	paginationParams := composables.UsePaginated(r)
	filterDTO, err := composables.UseQuery(&dtos.FuelEntryFilterDTO{
		Page:     paginationParams.Page,
		PageSize: paginationParams.Limit,
	}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := &fuel_entry.FindParams{
		TenantID: tenantID,
		Limit:    filterDTO.PageSize,
		Offset:   (filterDTO.Page - 1) * filterDTO.PageSize,
		SortBy:   fuel_entry.FieldDate,
		SortDesc: true,
	}

	if filterDTO.VehicleID != uuid.Nil {
		params.VehicleID = &filterDTO.VehicleID
	}
	if filterDTO.DriverID != uuid.Nil {
		params.DriverID = &filterDTO.DriverID
	}
	if filterDTO.FuelType != "" {
		fuelType, err := enums.ParseFuelType(filterDTO.FuelType)
		if err == nil {
			params.FuelType = &fuelType
		}
	}
	if !filterDTO.StartDate.IsZero() {
		startDateStr := filterDTO.StartDate.Format("2006-01-02")
		params.DateFrom = &startDateStr
	}
	if !filterDTO.EndDate.IsZero() {
		endDateStr := filterDTO.EndDate.Format("2006-01-02")
		params.DateTo = &endDateStr
	}

	fuelEntries, err := c.fuelService.GetPaginated(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	total, err := c.fuelService.Count(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	viewModels, err := c.buildFuelEntryListViewModels(ctx, fuelEntries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vehicleOptions, err := c.getVehicleOptions(ctx, tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	props := c.buildListPageProps(r.URL.Path, viewModels, total, filterDTO, vehicleOptions)

	if htmx.IsHxRequest(r) {
		w.WriteHeader(http.StatusOK)
		fuel.FuelTable(props).Render(ctx, w)
	} else {
		w.WriteHeader(http.StatusOK)
		fuel.Index(props).Render(ctx, w)
	}
}

func (c *FuelController) New(w http.ResponseWriter, r *http.Request) {
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

	driverOptions, err := c.getDriverOptions(ctx, tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	props := c.buildNewPageProps(vehicleOptions, driverOptions)
	w.WriteHeader(http.StatusOK)
	fuel.New(props).Render(ctx, w)
}

func (c *FuelController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dto, err := composables.UseForm(&dtos.FuelEntryCreateDTO{}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	f, err := mappers.FuelEntryCreateDTOToDomain(*dto, tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = c.fuelService.Create(ctx, f)
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

func (c *FuelController) Edit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid fuel entry ID", http.StatusBadRequest)
		return
	}

	f, err := c.fuelService.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	vehicleOptions, err := c.getVehicleOptions(ctx, tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	driverOptions, err := c.getDriverOptions(ctx, tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	props := c.buildEditPageProps(f, vehicleOptions, driverOptions)
	w.WriteHeader(http.StatusOK)
	fuel.Edit(props).Render(ctx, w)
}

func (c *FuelController) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid fuel entry ID", http.StatusBadRequest)
		return
	}

	dto, err := composables.UseForm(&dtos.FuelEntryUpdateDTO{}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto.ID = id

	existing, err := c.fuelService.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	updated, err := mappers.FuelEntryUpdateDTOToDomain(*dto, existing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = c.fuelService.Update(ctx, updated)
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

func (c *FuelController) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid fuel entry ID", http.StatusBadRequest)
		return
	}

	err = c.fuelService.Delete(ctx, id)
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

func (c *FuelController) buildFuelEntryListViewModels(ctx context.Context, entries []fuel_entry.FuelEntry) ([]viewmodels.FuelEntryListViewModel, error) {
	viewModels := make([]viewmodels.FuelEntryListViewModel, 0, len(entries))

	for _, entry := range entries {
		vehicleName := ""
		v, err := c.vehicleService.GetByID(ctx, entry.VehicleID())
		if err == nil {
			vehicleName = fmt.Sprintf("%s %s (%d)", v.Make(), v.Model(), v.Year())
		}

		driverName := ""
		if entry.DriverID() != nil {
			d, err := c.driverService.GetByID(ctx, *entry.DriverID())
			if err == nil {
				driverName = fmt.Sprintf("%s %s", d.FirstName(), d.LastName())
			}
		}

		efficiency, _ := c.fuelService.CalculateEfficiency(ctx, entry)
		viewModels = append(viewModels, mappers.FuelEntryToListViewModel(entry, vehicleName, driverName, efficiency))
	}

	return viewModels, nil
}

func (c *FuelController) getVehicleOptions(ctx context.Context, tenantID uuid.UUID) ([]struct{ Value, Label string }, error) {
	params := &vehicle.FindParams{
		TenantID: tenantID,
		Limit:    1000,
		Offset:   0,
		SortBy:   vehicle.FieldMake,
		SortDesc: false,
	}

	vehicles, err := c.vehicleService.GetPaginated(ctx, params)
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

func (c *FuelController) getDriverOptions(ctx context.Context, tenantID uuid.UUID) ([]struct{ Value, Label string }, error) {
	params := &driver.FindParams{
		TenantID: tenantID,
		Limit:    1000,
		Offset:   0,
		SortBy:   driver.FieldFirstName,
		SortDesc: false,
	}

	drivers, err := c.driverService.GetPaginated(ctx, params)
	if err != nil {
		return nil, err
	}

	options := make([]struct{ Value, Label string }, 0, len(drivers))
	for _, d := range drivers {
		options = append(options, struct{ Value, Label string }{
			Value: d.ID().String(),
			Label: fmt.Sprintf("%s %s - %s", d.FirstName(), d.LastName(), d.LicenseNumber()),
		})
	}

	return options, nil
}

func (c *FuelController) buildListPageProps(
	path string,
	entries []viewmodels.FuelEntryListViewModel,
	total int64,
	filter *dtos.FuelEntryFilterDTO,
	vehicleOptions []struct{ Value, Label string },
) *fuel.IndexPageProps {
	paginationState := pagination.New(path, filter.Page, filter.PageSize, int(total))

	fuelTypeOptions := []struct {
		Value     string
		MessageID string
	}{
		{Value: enums.FuelTypeGasoline.String(), MessageID: "Fleet.Enums.FuelType.Gasoline"},
		{Value: enums.FuelTypeDiesel.String(), MessageID: "Fleet.Enums.FuelType.Diesel"},
		{Value: enums.FuelTypeElectric.String(), MessageID: "Fleet.Enums.FuelType.Electric"},
		{Value: enums.FuelTypeHybrid.String(), MessageID: "Fleet.Enums.FuelType.Hybrid"},
		{Value: enums.FuelTypeCNG.String(), MessageID: "Fleet.Enums.FuelType.CNG"},
	}

	return &fuel.IndexPageProps{
		FuelEntries:     entries,
		PaginationState: paginationState,
		FuelTypeOptions: fuelTypeOptions,
		VehicleOptions:  vehicleOptions,
	}
}

func (c *FuelController) buildNewPageProps(
	vehicleOptions []struct{ Value, Label string },
	driverOptions []struct{ Value, Label string },
) *fuel.CreatePageProps {
	fuelTypeOptions := []struct {
		Value     string
		MessageID string
	}{
		{Value: enums.FuelTypeGasoline.String(), MessageID: "Fleet.Enums.FuelType.Gasoline"},
		{Value: enums.FuelTypeDiesel.String(), MessageID: "Fleet.Enums.FuelType.Diesel"},
		{Value: enums.FuelTypeElectric.String(), MessageID: "Fleet.Enums.FuelType.Electric"},
		{Value: enums.FuelTypeHybrid.String(), MessageID: "Fleet.Enums.FuelType.Hybrid"},
		{Value: enums.FuelTypeCNG.String(), MessageID: "Fleet.Enums.FuelType.CNG"},
	}

	return &fuel.CreatePageProps{
		FuelEntry: &viewmodels.FuelEntryFormViewModel{
			Errors: make(map[string]string),
		},
		PostPath:        "/fleet/fuel",
		VehicleOptions:  vehicleOptions,
		DriverOptions:   driverOptions,
		FuelTypeOptions: fuelTypeOptions,
	}
}

func (c *FuelController) buildEditPageProps(
	f fuel_entry.FuelEntry,
	vehicleOptions []struct{ Value, Label string },
	driverOptions []struct{ Value, Label string },
) *fuel.EditPageProps {
	fuelTypeOptions := []struct {
		Value     string
		MessageID string
	}{
		{Value: enums.FuelTypeGasoline.String(), MessageID: "Fleet.Enums.FuelType.Gasoline"},
		{Value: enums.FuelTypeDiesel.String(), MessageID: "Fleet.Enums.FuelType.Diesel"},
		{Value: enums.FuelTypeElectric.String(), MessageID: "Fleet.Enums.FuelType.Electric"},
		{Value: enums.FuelTypeHybrid.String(), MessageID: "Fleet.Enums.FuelType.Hybrid"},
		{Value: enums.FuelTypeCNG.String(), MessageID: "Fleet.Enums.FuelType.CNG"},
	}

	formVM := mappers.FuelEntryToFormViewModel(f, make(map[string]string))
	return &fuel.EditPageProps{
		FuelEntry:       &formVM,
		PostPath:        fmt.Sprintf("/fleet/fuel/%s", f.ID().String()),
		VehicleOptions:  vehicleOptions,
		DriverOptions:   driverOptions,
		FuelTypeOptions: fuelTypeOptions,
	}
}
