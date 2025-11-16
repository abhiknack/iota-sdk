package fleet

import (
	"embed"

	icons "github.com/iota-uz/icons/phosphor"
	"github.com/iota-uz/iota-sdk/modules/fleet/infrastructure/persistence"
	"github.com/iota-uz/iota-sdk/modules/fleet/presentation/controllers"
	"github.com/iota-uz/iota-sdk/modules/fleet/services"
	"github.com/iota-uz/iota-sdk/pkg/application"
	"github.com/iota-uz/iota-sdk/pkg/configuration"
	"github.com/iota-uz/iota-sdk/pkg/spotlight"
)

//go:embed presentation/locales/*.json
var localeFiles embed.FS

//go:embed infrastructure/persistence/schema/*.sql
var migrationFiles embed.FS

func NewModule() application.Module {
	return &Module{}
}

type Module struct {
}

func (m *Module) Register(app application.Application) error {
	vehicleRepo := persistence.NewVehicleRepository()
	driverRepo := persistence.NewDriverRepository()
	tripRepo := persistence.NewTripRepository()
	maintenanceRepo := persistence.NewMaintenanceRepository()
	fuelEntryRepo := persistence.NewFuelEntryRepository()

	logger := configuration.Use().Logger()

	vehicleService := services.NewVehicleService(vehicleRepo, app.EventPublisher())
	driverService := services.NewDriverService(driverRepo, app.EventPublisher())
	tripService := services.NewTripService(tripRepo, vehicleRepo, driverRepo, app.EventPublisher())
	maintenanceService := services.NewMaintenanceService(maintenanceRepo, vehicleRepo, app.EventPublisher())
	fuelService := services.NewFuelService(fuelEntryRepo, app.EventPublisher())
	analyticsService := services.NewAnalyticsService(vehicleRepo, driverRepo, tripRepo, maintenanceRepo, fuelEntryRepo)
	notificationService := services.NewNotificationService(vehicleService, driverService, maintenanceService, fuelService, app.EventPublisher(), logger)
	schedulerService := services.NewSchedulerService(notificationService, logger)

	app.RegisterServices(
		vehicleService,
		driverService,
		tripService,
		maintenanceService,
		fuelService,
		analyticsService,
		notificationService,
		schedulerService,
	)

	app.RegisterControllers(
		controllers.NewVehicleController(app),
		controllers.NewDriverController(app),
		controllers.NewTripController(app),
		controllers.NewMaintenanceController(app),
		controllers.NewFuelController(app),
		controllers.NewDashboardController(app),
	)

	app.QuickLinks().Add(
		spotlight.NewQuickLink(
			icons.Gauge(icons.Props{Size: "24"}),
			DashboardItem.Name,
			DashboardItem.Href,
		),
		spotlight.NewQuickLink(
			icons.Car(icons.Props{Size: "24"}),
			VehiclesItem.Name,
			VehiclesItem.Href,
		),
		spotlight.NewQuickLink(
			icons.User(icons.Props{Size: "24"}),
			DriversItem.Name,
			DriversItem.Href,
		),
		spotlight.NewQuickLink(
			icons.MapPin(icons.Props{Size: "24"}),
			TripsItem.Name,
			TripsItem.Href,
		),
		spotlight.NewQuickLink(
			icons.Wrench(icons.Props{Size: "24"}),
			MaintenanceItem.Name,
			MaintenanceItem.Href,
		),
		spotlight.NewQuickLink(
			icons.GasPump(icons.Props{Size: "24"}),
			FuelItem.Name,
			FuelItem.Href,
		),
		spotlight.NewQuickLink(
			icons.PlusCircle(icons.Props{Size: "24"}),
			"Vehicles.List.New",
			"/fleet/vehicles/new",
		),
		spotlight.NewQuickLink(
			icons.PlusCircle(icons.Props{Size: "24"}),
			"Drivers.List.New",
			"/fleet/drivers/new",
		),
		spotlight.NewQuickLink(
			icons.PlusCircle(icons.Props{Size: "24"}),
			"Trips.List.New",
			"/fleet/trips/new",
		),
	)

	app.RegisterLocaleFiles(&localeFiles)
	app.Migrations().RegisterSchema(&migrationFiles)
	return nil
}

func (m *Module) Name() string {
	return "fleet"
}
