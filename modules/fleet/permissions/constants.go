package permissions

import (
	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/core/domain/entities/permission"
)

const (
	ResourceVehicle     permission.Resource = "vehicle"
	ResourceDriver      permission.Resource = "driver"
	ResourceTrip        permission.Resource = "trip"
	ResourceMaintenance permission.Resource = "maintenance"
	ResourceFuelEntry   permission.Resource = "fuel_entry"
)

var (
	VehicleCreate = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		Name:     "Vehicle.Create",
		Resource: ResourceVehicle,
		Action:   permission.ActionCreate,
		Modifier: permission.ModifierAll,
	}
	VehicleRead = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
		Name:     "Vehicle.Read",
		Resource: ResourceVehicle,
		Action:   permission.ActionRead,
		Modifier: permission.ModifierAll,
	}
	VehicleUpdate = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
		Name:     "Vehicle.Update",
		Resource: ResourceVehicle,
		Action:   permission.ActionUpdate,
		Modifier: permission.ModifierAll,
	}
	VehicleDelete = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"),
		Name:     "Vehicle.Delete",
		Resource: ResourceVehicle,
		Action:   permission.ActionDelete,
		Modifier: permission.ModifierAll,
	}
	DriverCreate = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"),
		Name:     "Driver.Create",
		Resource: ResourceDriver,
		Action:   permission.ActionCreate,
		Modifier: permission.ModifierAll,
	}
	DriverRead = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"),
		Name:     "Driver.Read",
		Resource: ResourceDriver,
		Action:   permission.ActionRead,
		Modifier: permission.ModifierAll,
	}
	DriverUpdate = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440007"),
		Name:     "Driver.Update",
		Resource: ResourceDriver,
		Action:   permission.ActionUpdate,
		Modifier: permission.ModifierAll,
	}
	DriverDelete = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440008"),
		Name:     "Driver.Delete",
		Resource: ResourceDriver,
		Action:   permission.ActionDelete,
		Modifier: permission.ModifierAll,
	}
	TripCreate = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440009"),
		Name:     "Trip.Create",
		Resource: ResourceTrip,
		Action:   permission.ActionCreate,
		Modifier: permission.ModifierAll,
	}
	TripRead = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440010"),
		Name:     "Trip.Read",
		Resource: ResourceTrip,
		Action:   permission.ActionRead,
		Modifier: permission.ModifierAll,
	}
	TripUpdate = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440011"),
		Name:     "Trip.Update",
		Resource: ResourceTrip,
		Action:   permission.ActionUpdate,
		Modifier: permission.ModifierAll,
	}
	TripDelete = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440012"),
		Name:     "Trip.Delete",
		Resource: ResourceTrip,
		Action:   permission.ActionDelete,
		Modifier: permission.ModifierAll,
	}
	MaintenanceCreate = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440013"),
		Name:     "Maintenance.Create",
		Resource: ResourceMaintenance,
		Action:   permission.ActionCreate,
		Modifier: permission.ModifierAll,
	}
	MaintenanceRead = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440014"),
		Name:     "Maintenance.Read",
		Resource: ResourceMaintenance,
		Action:   permission.ActionRead,
		Modifier: permission.ModifierAll,
	}
	MaintenanceUpdate = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440015"),
		Name:     "Maintenance.Update",
		Resource: ResourceMaintenance,
		Action:   permission.ActionUpdate,
		Modifier: permission.ModifierAll,
	}
	MaintenanceDelete = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440016"),
		Name:     "Maintenance.Delete",
		Resource: ResourceMaintenance,
		Action:   permission.ActionDelete,
		Modifier: permission.ModifierAll,
	}
	FuelEntryCreate = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440017"),
		Name:     "FuelEntry.Create",
		Resource: ResourceFuelEntry,
		Action:   permission.ActionCreate,
		Modifier: permission.ModifierAll,
	}
	FuelEntryRead = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440018"),
		Name:     "FuelEntry.Read",
		Resource: ResourceFuelEntry,
		Action:   permission.ActionRead,
		Modifier: permission.ModifierAll,
	}
	FuelEntryUpdate = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440019"),
		Name:     "FuelEntry.Update",
		Resource: ResourceFuelEntry,
		Action:   permission.ActionUpdate,
		Modifier: permission.ModifierAll,
	}
	FuelEntryDelete = &permission.Permission{
		ID:       uuid.MustParse("550e8400-e29b-41d4-a716-446655440020"),
		Name:     "FuelEntry.Delete",
		Resource: ResourceFuelEntry,
		Action:   permission.ActionDelete,
		Modifier: permission.ModifierAll,
	}
)

var Permissions = []*permission.Permission{
	VehicleCreate,
	VehicleRead,
	VehicleUpdate,
	VehicleDelete,
	DriverCreate,
	DriverRead,
	DriverUpdate,
	DriverDelete,
	TripCreate,
	TripRead,
	TripUpdate,
	TripDelete,
	MaintenanceCreate,
	MaintenanceRead,
	MaintenanceUpdate,
	MaintenanceDelete,
	FuelEntryCreate,
	FuelEntryRead,
	FuelEntryUpdate,
	FuelEntryDelete,
}
