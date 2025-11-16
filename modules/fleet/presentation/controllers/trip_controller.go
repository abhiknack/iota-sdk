package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/iota-uz/iota-sdk/components/base/pagination"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/driver"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/trip"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/aggregates/vehicle"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/controllers/dtos"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/mappers"
	tripstpl "github.com/iota-uz/iota-sdk/modules/fleet/presentation/templates/pages/trips"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/viewmodels"
	"github.com/iota-uz/iota-sdk/modules/fleet/services"
	"github.com/iota-uz/iota-sdk/pkg/application"
	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/htmx"
	"github.com/iota-uz/iota-sdk/pkg/middleware"
	"github.com/iota-uz/iota-sdk/pkg/shared"
)

type TripController struct {
	app         application.Application
	tripService *services.TripService
	basePath    string
}

func NewTripController(app application.Application) application.Controller {
	return &TripController{
		app:         app,
		tripService: app.Service(services.TripService{}).(*services.TripService),
		basePath:    "/fleet/trips",
	}
}

func (c *TripController) Key() string {
	return c.basePath
}

func (c *TripController) Register(r *mux.Router) {
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
	getRouter.HandleFunc("/{id}", c.Detail).Methods(http.MethodGet)

	setRouter := r.PathPrefix(c.basePath).Subrouter()
	setRouter.Use(commonMiddleware...)
	setRouter.Use(middleware.WithTransaction())
	setRouter.HandleFunc("", c.Create).Methods(http.MethodPost)
	setRouter.HandleFunc("/{id}/complete", c.Complete).Methods(http.MethodPost)
	setRouter.HandleFunc("/{id}/cancel", c.Cancel).Methods(http.MethodPost)
}

func (c *TripController) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	paginationParams := composables.UsePaginated(r)
	filterDTO, err := composables.UseQuery(&dtos.TripFilterDTO{
		Page:     paginationParams.Page,
		PageSize: paginationParams.Limit,
	}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := &trip.FindParams{
		TenantID: tenantID,
		Limit:    filterDTO.PageSize,
		Offset:   (filterDTO.Page - 1) * filterDTO.PageSize,
		SortBy:   trip.FieldStartTime,
		SortDesc: true,
	}

	if filterDTO.VehicleID != uuid.Nil {
		params.VehicleID = &filterDTO.VehicleID
	}
	if filterDTO.DriverID != uuid.Nil {
		params.DriverID = &filterDTO.DriverID
	}
	if filterDTO.Status != "" {
		status, err := enums.ParseTripStatus(filterDTO.Status)
		if err == nil {
			params.Status = &status
		}
	}
	if !filterDTO.StartDate.IsZero() {
		params.StartTimeFrom = &filterDTO.StartDate
	}
	if !filterDTO.EndDate.IsZero() {
		params.StartTimeTo = &filterDTO.EndDate
	}

	trips, err := c.tripService.GetPaginated(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	total, err := c.tripService.Count(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	viewModels, err := c.buildTripListViewModels(ctx, trips)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	props := c.buildListPageProps(r.URL.Path, viewModels, total, filterDTO, vehicleOptions, driverOptions)

	if htmx.IsHxRequest(r) {
		w.WriteHeader(http.StatusOK)
		tripstpl.TripsTable(props).Render(ctx, w)
	} else {
		w.WriteHeader(http.StatusOK)
		tripstpl.Index(props).Render(ctx, w)
	}
}

func (c *TripController) New(w http.ResponseWriter, r *http.Request) {
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
	tripstpl.New(props).Render(ctx, w)
}

func (c *TripController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, err := composables.UseTenantID(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dto, err := composables.UseForm(&dtos.TripCreateDTO{}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t := trip.NewTrip(
		uuid.New(),
		tenantID,
		dto.VehicleID,
		dto.DriverID,
		dto.Origin,
		dto.Destination,
		dto.Purpose,
		dto.StartTime,
		dto.StartOdometer,
		trip.WithStatus(enums.TripStatusScheduled),
	)

	created, err := c.tripService.Create(ctx, t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if htmx.IsHxRequest(r) {
		htmx.Redirect(w, c.basePath)
	} else {
		shared.Redirect(w, r, c.basePath)
	}
	fmt.Fprintf(w, "Trip created: %s", created.ID())
}

func (c *TripController) Detail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid trip ID", http.StatusBadRequest)
		return
	}

	t, err := c.tripService.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	vehicleName, driverName, err := c.getVehicleAndDriverNames(ctx, t.VehicleID(), t.DriverID())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	props := c.buildDetailPageProps(t, vehicleName, driverName)
	w.WriteHeader(http.StatusOK)
	tripstpl.Detail(props).Render(ctx, w)
}

func (c *TripController) Complete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid trip ID", http.StatusBadRequest)
		return
	}

	dto, err := composables.UseForm(&dtos.TripCompleteDTO{}, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = c.tripService.Complete(ctx, id, dto.EndTime, dto.EndOdometer)
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

func (c *TripController) Cancel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid trip ID", http.StatusBadRequest)
		return
	}

	reason := r.FormValue("Reason")
	if reason == "" {
		reason = "Cancelled by user"
	}

	_, err = c.tripService.Cancel(ctx, id, reason)
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

func (c *TripController) buildTripListViewModels(ctx context.Context, trips []trip.Trip) ([]viewmodels.TripListViewModel, error) {
	vehicleService := c.app.Service(services.VehicleService{}).(*services.VehicleService)
	driverService := c.app.Service(services.DriverService{}).(*services.DriverService)

	viewModels := make([]viewmodels.TripListViewModel, 0, len(trips))
	for _, t := range trips {
		v, err := vehicleService.GetByID(ctx, t.VehicleID())
		if err != nil {
			return nil, err
		}

		d, err := driverService.GetByID(ctx, t.DriverID())
		if err != nil {
			return nil, err
		}

		vehicleName := fmt.Sprintf("%s %s (%d)", v.Make(), v.Model(), v.Year())
		driverName := fmt.Sprintf("%s %s", d.FirstName(), d.LastName())

		viewModels = append(viewModels, mappers.TripToListViewModel(t, vehicleName, driverName))
	}

	return viewModels, nil
}

func (c *TripController) getVehicleOptions(ctx context.Context, tenantID uuid.UUID) ([]viewmodels.VehicleOption, error) {
	vehicleService := c.app.Service(services.VehicleService{}).(*services.VehicleService)

	params := &vehicle.FindParams{
		TenantID: tenantID,
		Limit:    1000,
		Offset:   0,
	}

	vehicles, err := vehicleService.GetPaginated(ctx, params)
	if err != nil {
		return nil, err
	}

	options := make([]viewmodels.VehicleOption, 0, len(vehicles))
	for _, v := range vehicles {
		options = append(options, viewmodels.VehicleOption{
			ID:    v.ID().String(),
			Label: fmt.Sprintf("%s %s (%d) - %s", v.Make(), v.Model(), v.Year(), v.LicensePlate()),
		})
	}

	return options, nil
}

func (c *TripController) getDriverOptions(ctx context.Context, tenantID uuid.UUID) ([]viewmodels.DriverOption, error) {
	driverService := c.app.Service(services.DriverService{}).(*services.DriverService)

	params := &driver.FindParams{
		TenantID: tenantID,
		Limit:    1000,
		Offset:   0,
	}

	drivers, err := driverService.GetPaginated(ctx, params)
	if err != nil {
		return nil, err
	}

	options := make([]viewmodels.DriverOption, 0, len(drivers))
	for _, d := range drivers {
		options = append(options, viewmodels.DriverOption{
			ID:    d.ID().String(),
			Label: fmt.Sprintf("%s %s - %s", d.FirstName(), d.LastName(), d.LicenseNumber()),
		})
	}

	return options, nil
}

func (c *TripController) getVehicleAndDriverNames(ctx context.Context, vehicleID, driverID uuid.UUID) (string, string, error) {
	vehicleService := c.app.Service(services.VehicleService{}).(*services.VehicleService)
	driverService := c.app.Service(services.DriverService{}).(*services.DriverService)

	v, err := vehicleService.GetByID(ctx, vehicleID)
	if err != nil {
		return "", "", err
	}

	d, err := driverService.GetByID(ctx, driverID)
	if err != nil {
		return "", "", err
	}

	vehicleName := fmt.Sprintf("%s %s (%d)", v.Make(), v.Model(), v.Year())
	driverName := fmt.Sprintf("%s %s", d.FirstName(), d.LastName())

	return vehicleName, driverName, nil
}

func (c *TripController) buildListPageProps(
	urlPath string,
	tripsList []viewmodels.TripListViewModel,
	total int64,
	filterDTO *dtos.TripFilterDTO,
	vehicleOptions []viewmodels.VehicleOption,
	driverOptions []viewmodels.DriverOption,
) *tripstpl.IndexPageProps {
	return &tripstpl.IndexPageProps{
		Trips: tripsList,
		PaginationState: pagination.New(
			urlPath,
			filterDTO.Page,
			int(total),
			filterDTO.PageSize,
		),
		StatusOptions: []struct {
			Value     string
			MessageID string
		}{
			{Value: "Scheduled", MessageID: "Fleet.Enums.TripStatus.Scheduled"},
			{Value: "InProgress", MessageID: "Fleet.Enums.TripStatus.InProgress"},
			{Value: "Completed", MessageID: "Fleet.Enums.TripStatus.Completed"},
			{Value: "Cancelled", MessageID: "Fleet.Enums.TripStatus.Cancelled"},
		},
		VehicleOptions: convertVehicleOptions(vehicleOptions),
		DriverOptions:  convertDriverOptions(driverOptions),
	}
}

func (c *TripController) buildNewPageProps(
	vehicleOptions []viewmodels.VehicleOption,
	driverOptions []viewmodels.DriverOption,
) *tripstpl.CreatePageProps {
	return &tripstpl.CreatePageProps{
		Trip: &viewmodels.TripFormViewModel{
			Errors:      make(map[string]string),
			Vehicles:    vehicleOptions,
			Drivers:     driverOptions,
			HasConflict: false,
		},
		PostPath: c.basePath,
	}
}

func (c *TripController) buildDetailPageProps(
	t trip.Trip,
	vehicleName, driverName string,
) *tripstpl.DetailPageProps {
	vm := mappers.TripToDetailViewModel(t, vehicleName, driverName)
	return &tripstpl.DetailPageProps{
		Trip: &vm,
	}
}

func convertVehicleOptions(options []viewmodels.VehicleOption) []struct {
	Value string
	Label string
} {
	result := make([]struct {
		Value string
		Label string
	}, len(options))
	for i, opt := range options {
		result[i] = struct {
			Value string
			Label string
		}{
			Value: opt.ID,
			Label: opt.Label,
		}
	}
	return result
}

func convertDriverOptions(options []viewmodels.DriverOption) []struct {
	Value string
	Label string
} {
	result := make([]struct {
		Value string
		Label string
	}, len(options))
	for i, opt := range options {
		result[i] = struct {
			Value string
			Label string
		}{
			Value: opt.ID,
			Label: opt.Label,
		}
	}
	return result
}
