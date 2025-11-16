# Fleet Management Module - Design Document

## Overview

The Fleet Management Module is a comprehensive solution for managing vehicle fleets within the IOTA SDK multi-tenant architecture. It follows Domain-Driven Design (DDD) principles with clear separation between domain logic, infrastructure, services, and presentation layers. The module integrates with the existing IOTA SDK event system, authentication, and RBAC framework.

## Architecture

### Layer Structure

```
modules/fleet/
├── domain/
│   ├── aggregates/
│   │   ├── vehicle/
│   │   ├── driver/
│   │   ├── trip/
│   │   ├── maintenance/
│   │   └── fuel_entry/
│   ├── value_objects/
│   │   ├── vehicle_status.go
│   │   ├── license_info.go
│   │   └── location.go
│   └── enums/
│       ├── vehicle_status.go
│       ├── fuel_type.go
│       └── service_type.go
├── infrastructure/
│   └── persistence/
│       ├── models/
│       ├── vehicle_repository.go
│       ├── driver_repository.go
│       ├── trip_repository.go
│       ├── maintenance_repository.go
│       ├── fuel_entry_repository.go
│       ├── fleet_mappers.go
│       └── schema/fleet-schema.sql
├── services/
│   ├── vehicle_service.go
│   ├── driver_service.go
│   ├── trip_service.go
│   ├── maintenance_service.go
│   ├── fuel_service.go
│   └── analytics_service.go
├── presentation/
│   ├── controllers/
│   │   ├── dtos/
│   │   ├── vehicle_controller.go
│   │   ├── driver_controller.go
│   │   ├── trip_controller.go
│   │   ├── maintenance_controller.go
│   │   ├── fuel_controller.go
│   │   └── dashboard_controller.go
│   ├── templates/
│   │   ├── pages/
│   │   │   ├── dashboard/
│   │   │   ├── vehicles/
│   │   │   ├── drivers/
│   │   │   ├── trips/
│   │   │   ├── maintenance/
│   │   │   └── fuel/
│   │   └── components/
│   ├── viewmodels/
│   ├── mappers/
│   └── locales/
├── module.go
├── links.go
└── permissions/
```

## Components and Interfaces

### Domain Aggregates

#### Vehicle Aggregate
```go
type Vehicle interface {
    ID() uuid.UUID
    TenantID() uuid.UUID
    Make() string
    Model() string
    Year() int
    VIN() string
    LicensePlate() string
    Status() VehicleStatus
    CurrentOdometer() int
    RegistrationExpiry() time.Time
    InsuranceExpiry() time.Time
    CreatedAt() time.Time
    UpdatedAt() time.Time
    
    UpdateStatus(status VehicleStatus) Vehicle
    UpdateOdometer(reading int) Vehicle
    UpdateDetails(make, model string, year int) Vehicle
}

type VehicleRepository interface {
    GetByID(ctx context.Context, id uuid.UUID) (Vehicle, error)
    GetPaginated(ctx context.Context, params *FindParams) ([]Vehicle, error)
    Count(ctx context.Context, params *FindParams) (int64, error)
    GetByStatus(ctx context.Context, status VehicleStatus) ([]Vehicle, error)
    GetExpiringRegistrations(ctx context.Context, days int) ([]Vehicle, error)
    Create(ctx context.Context, vehicle Vehicle) (Vehicle, error)
    Update(ctx context.Context, vehicle Vehicle) (Vehicle, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```

#### Driver Aggregate
```go
type Driver interface {
    ID() uuid.UUID
    TenantID() uuid.UUID
    UserID() uuid.UUID
    FirstName() string
    LastName() string
    LicenseNumber() string
    LicenseExpiry() time.Time
    Phone() string
    Email() string
    Status() DriverStatus
    CreatedAt() time.Time
    UpdatedAt() time.Time
    
    UpdateLicense(number string, expiry time.Time) Driver
    UpdateContact(phone, email string) Driver
    UpdateStatus(status DriverStatus) Driver
}

type DriverRepository interface {
    GetByID(ctx context.Context, id uuid.UUID) (Driver, error)
    GetByUserID(ctx context.Context, userID uuid.UUID) (Driver, error)
    GetPaginated(ctx context.Context, params *FindParams) ([]Driver, error)
    Count(ctx context.Context, params *FindParams) (int64, error)
    GetExpiringLicenses(ctx context.Context, days int) ([]Driver, error)
    GetAvailable(ctx context.Context, startTime, endTime time.Time) ([]Driver, error)
    Create(ctx context.Context, driver Driver) (Driver, error)
    Update(ctx context.Context, driver Driver) (Driver, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```

#### Trip Aggregate
```go
type Trip interface {
    ID() uuid.UUID
    TenantID() uuid.UUID
    VehicleID() uuid.UUID
    DriverID() uuid.UUID
    Origin() string
    Destination() string
    Purpose() string
    StartTime() time.Time
    EndTime() *time.Time
    StartOdometer() int
    EndOdometer() *int
    Status() TripStatus
    CreatedAt() time.Time
    UpdatedAt() time.Time
    
    Complete(endTime time.Time, endOdometer int) Trip
    Cancel(reason string) Trip
    UpdateRoute(origin, destination string) Trip
}

type TripRepository interface {
    GetByID(ctx context.Context, id uuid.UUID) (Trip, error)
    GetPaginated(ctx context.Context, params *FindParams) ([]Trip, error)
    Count(ctx context.Context, params *FindParams) (int64, error)
    GetByVehicle(ctx context.Context, vehicleID uuid.UUID) ([]Trip, error)
    GetByDriver(ctx context.Context, driverID uuid.UUID) ([]Trip, error)
    GetActiveTrips(ctx context.Context) ([]Trip, error)
    CheckConflict(ctx context.Context, vehicleID uuid.UUID, startTime, endTime time.Time) (bool, error)
    Create(ctx context.Context, trip Trip) (Trip, error)
    Update(ctx context.Context, trip Trip) (Trip, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```

#### Maintenance Aggregate
```go
type Maintenance interface {
    ID() uuid.UUID
    TenantID() uuid.UUID
    VehicleID() uuid.UUID
    ServiceType() ServiceType
    ServiceDate() time.Time
    Odometer() int
    Cost() float64
    ServiceProvider() string
    Description() string
    NextServiceDue() *time.Time
    NextServiceOdometer() *int
    CreatedAt() time.Time
    UpdatedAt() time.Time
    
    UpdateCost(cost float64) Maintenance
    UpdateNextService(date *time.Time, odometer *int) Maintenance
}

type MaintenanceRepository interface {
    GetByID(ctx context.Context, id uuid.UUID) (Maintenance, error)
    GetPaginated(ctx context.Context, params *FindParams) ([]Maintenance, error)
    Count(ctx context.Context, params *FindParams) (int64, error)
    GetByVehicle(ctx context.Context, vehicleID uuid.UUID) ([]Maintenance, error)
    GetDueMaintenance(ctx context.Context) ([]Maintenance, error)
    Create(ctx context.Context, maintenance Maintenance) (Maintenance, error)
    Update(ctx context.Context, maintenance Maintenance) (Maintenance, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```

#### Fuel Entry Aggregate
```go
type FuelEntry interface {
    ID() uuid.UUID
    TenantID() uuid.UUID
    VehicleID() uuid.UUID
    DriverID() *uuid.UUID
    Date() time.Time
    Quantity() float64
    Cost() float64
    Odometer() int
    FuelType() FuelType
    Location() string
    CreatedAt() time.Time
    UpdatedAt() time.Time
    
    CalculateEfficiency(previousOdometer int) float64
    UpdateCost(cost float64) FuelEntry
}

type FuelEntryRepository interface {
    GetByID(ctx context.Context, id uuid.UUID) (FuelEntry, error)
    GetPaginated(ctx context.Context, params *FindParams) ([]FuelEntry, error)
    Count(ctx context.Context, params *FindParams) (int64, error)
    GetByVehicle(ctx context.Context, vehicleID uuid.UUID) ([]FuelEntry, error)
    GetByDriver(ctx context.Context, driverID uuid.UUID) ([]FuelEntry, error)
    GetLastEntry(ctx context.Context, vehicleID uuid.UUID) (FuelEntry, error)
    Create(ctx context.Context, entry FuelEntry) (FuelEntry, error)
    Update(ctx context.Context, entry FuelEntry) (FuelEntry, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### Value Objects

#### VehicleStatus
```go
type VehicleStatus int

const (
    VehicleStatusAvailable VehicleStatus = iota
    VehicleStatusInUse
    VehicleStatusMaintenance
    VehicleStatusOutOfService
    VehicleStatusRetired
)
```

#### FuelType
```go
type FuelType int

const (
    FuelTypeGasoline FuelType = iota
    FuelTypeDiesel
    FuelTypeElectric
    FuelTypeHybrid
    FuelTypeCNG
)
```

#### ServiceType
```go
type ServiceType int

const (
    ServiceTypeOilChange ServiceType = iota
    ServiceTypeTireRotation
    ServiceTypeBrakeService
    ServiceTypeInspection
    ServiceTypeRepair
    ServiceTypeOther
)
```

### Service Layer

#### VehicleService
- Manages vehicle lifecycle (create, update, delete)
- Handles status transitions with validation
- Publishes domain events (VehicleCreated, VehicleUpdated, StatusChanged)
- Checks for expiring registrations and insurance

#### DriverService
- Manages driver registration and updates
- Validates license information
- Checks for expiring licenses
- Manages driver availability

#### TripService
- Creates and manages trips
- Validates vehicle and driver availability
- Checks for scheduling conflicts
- Calculates trip statistics (duration, distance)
- Updates vehicle status during trip lifecycle

#### MaintenanceService
- Records maintenance activities
- Calculates next service due dates
- Identifies overdue maintenance
- Tracks maintenance costs

#### FuelService
- Records fuel entries
- Calculates fuel efficiency
- Identifies efficiency anomalies
- Aggregates fuel costs

#### AnalyticsService
- Generates dashboard statistics
- Calculates vehicle utilization rates
- Produces cost analysis reports
- Generates trend data for charts

## Data Models

### Database Schema

```sql
-- Vehicles table
CREATE TABLE fleet_vehicles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    make VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INT NOT NULL,
    vin VARCHAR(17) UNIQUE NOT NULL,
    license_plate VARCHAR(20) NOT NULL,
    status INT NOT NULL DEFAULT 0,
    current_odometer INT NOT NULL DEFAULT 0,
    registration_expiry TIMESTAMPTZ NOT NULL,
    insurance_expiry TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_fleet_vehicles_tenant_id ON fleet_vehicles(tenant_id);
CREATE INDEX idx_fleet_vehicles_status ON fleet_vehicles(status);
CREATE INDEX idx_fleet_vehicles_deleted_at ON fleet_vehicles(deleted_at);

-- Drivers table
CREATE TABLE fleet_drivers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    license_number VARCHAR(50) NOT NULL,
    license_expiry TIMESTAMPTZ NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255),
    status INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_fleet_drivers_tenant_id ON fleet_drivers(tenant_id);
CREATE INDEX idx_fleet_drivers_user_id ON fleet_drivers(user_id);
CREATE INDEX idx_fleet_drivers_deleted_at ON fleet_drivers(deleted_at);

-- Trips table
CREATE TABLE fleet_trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    vehicle_id UUID NOT NULL REFERENCES fleet_vehicles(id) ON DELETE CASCADE,
    driver_id UUID NOT NULL REFERENCES fleet_drivers(id) ON DELETE CASCADE,
    origin VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    purpose TEXT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    start_odometer INT NOT NULL,
    end_odometer INT,
    status INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_fleet_trips_tenant_id ON fleet_trips(tenant_id);
CREATE INDEX idx_fleet_trips_vehicle_id ON fleet_trips(vehicle_id);
CREATE INDEX idx_fleet_trips_driver_id ON fleet_trips(driver_id);
CREATE INDEX idx_fleet_trips_status ON fleet_trips(status);
CREATE INDEX idx_fleet_trips_deleted_at ON fleet_trips(deleted_at);

-- Maintenance records table
CREATE TABLE fleet_maintenance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    vehicle_id UUID NOT NULL REFERENCES fleet_vehicles(id) ON DELETE CASCADE,
    service_type INT NOT NULL,
    service_date TIMESTAMPTZ NOT NULL,
    odometer INT NOT NULL,
    cost DECIMAL(10,2) NOT NULL DEFAULT 0,
    service_provider VARCHAR(255),
    description TEXT,
    next_service_due TIMESTAMPTZ,
    next_service_odometer INT,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_fleet_maintenance_tenant_id ON fleet_maintenance(tenant_id);
CREATE INDEX idx_fleet_maintenance_vehicle_id ON fleet_maintenance(vehicle_id);
CREATE INDEX idx_fleet_maintenance_deleted_at ON fleet_maintenance(deleted_at);

-- Fuel entries table
CREATE TABLE fleet_fuel_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    vehicle_id UUID NOT NULL REFERENCES fleet_vehicles(id) ON DELETE CASCADE,
    driver_id UUID REFERENCES fleet_drivers(id) ON DELETE SET NULL,
    date TIMESTAMPTZ NOT NULL,
    quantity DECIMAL(10,2) NOT NULL,
    cost DECIMAL(10,2) NOT NULL,
    odometer INT NOT NULL,
    fuel_type INT NOT NULL,
    location VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_fleet_fuel_entries_tenant_id ON fleet_fuel_entries(tenant_id);
CREATE INDEX idx_fleet_fuel_entries_vehicle_id ON fleet_fuel_entries(vehicle_id);
CREATE INDEX idx_fleet_fuel_entries_driver_id ON fleet_fuel_entries(driver_id);
CREATE INDEX idx_fleet_fuel_entries_deleted_at ON fleet_fuel_entries(deleted_at);
```

## Error Handling

### Domain Errors
- `ErrVehicleNotFound`: Vehicle with specified ID not found
- `ErrDriverNotFound`: Driver with specified ID not found
- `ErrTripNotFound`: Trip with specified ID not found
- `ErrInvalidVehicleStatus`: Invalid status transition
- `ErrVehicleNotAvailable`: Vehicle not available for assignment
- `ErrDriverNotAvailable`: Driver not available for assignment
- `ErrTripConflict`: Scheduling conflict detected
- `ErrExpiredLicense`: Driver license has expired
- `ErrInvalidOdometer`: Odometer reading is invalid

### Error Handling Pattern
```go
if err != nil {
    return nil, fmt.Errorf("failed to create vehicle: %w", err)
}
```

All errors are wrapped using `fmt.Errorf` with `%w` verb for error chain preservation.

## Testing Strategy

### Unit Tests
- Domain aggregate behavior tests
- Value object validation tests
- Service layer business logic tests
- Repository query building tests

### Integration Tests (ITF Framework)
- Controller endpoint tests with HTMX assertions
- Service integration with repository tests
- Transaction rollback tests
- Multi-tenant isolation tests

### Test Structure
```go
func TestVehicleController(t *testing.T) {
    suite := itf.NewSuiteBuilder(t).
        WithModules(modules.BuiltInModules...).
        AsAdmin().
        Build()
    
    controller := controllers.NewVehicleController(suite.App())
    suite.Register(controller)
    
    t.Run("CreateVehicle", func(t *testing.T) {
        suite.POST("/fleet/vehicles").
            FormString("Make", "Toyota").
            FormString("Model", "Camry").
            FormInt("Year", 2023).
            HTMX().
            Assert(t).
            ExpectOK().
            ExpectBodyContains("Vehicle created")
    })
}
```

## UI/UX Design

### Dashboard
- Total vehicles count with status breakdown
- Active trips count
- Upcoming maintenance alerts
- Fuel cost summary for current month
- Utilization chart
- Cost trend chart

### Vehicle Management
- List view with filters (status, make, model)
- Detail view with tabs: Info, Trips, Maintenance, Fuel
- Create/Edit forms with validation
- Status change actions

### Driver Management
- List view with license expiry warnings
- Detail view with assigned vehicles and trip history
- Create/Edit forms
- License document upload

### Trip Management
- Calendar view for trip scheduling
- List view with filters (status, vehicle, driver)
- Create form with conflict detection
- Trip completion form

### Maintenance Tracking
- List view with due date warnings
- Create form with service type selection
- History view per vehicle
- Cost analysis charts

### Fuel Tracking
- Entry form with efficiency calculation
- List view with filters
- Efficiency trend charts
- Cost analysis by vehicle

## RBAC Permissions

```go
const (
    PermissionFleetView           = "fleet.view"
    PermissionFleetManage         = "fleet.manage"
    PermissionVehicleCreate       = "fleet.vehicle.create"
    PermissionVehicleUpdate       = "fleet.vehicle.update"
    PermissionVehicleDelete       = "fleet.vehicle.delete"
    PermissionDriverCreate        = "fleet.driver.create"
    PermissionDriverUpdate        = "fleet.driver.update"
    PermissionDriverDelete        = "fleet.driver.delete"
    PermissionTripCreate          = "fleet.trip.create"
    PermissionTripUpdate          = "fleet.trip.update"
    PermissionTripDelete          = "fleet.trip.delete"
    PermissionMaintenanceCreate   = "fleet.maintenance.create"
    PermissionMaintenanceUpdate   = "fleet.maintenance.update"
    PermissionMaintenanceDelete   = "fleet.maintenance.delete"
    PermissionFuelCreate          = "fleet.fuel.create"
    PermissionFuelUpdate          = "fleet.fuel.update"
    PermissionFuelDelete          = "fleet.fuel.delete"
    PermissionAnalyticsView       = "fleet.analytics.view"
)
```

## Integration Points

### Event System
- Publishes domain events for vehicle, driver, trip, maintenance, and fuel operations
- Other modules can subscribe to fleet events (e.g., finance module for cost tracking)

### User Management
- Drivers linked to user accounts via `user_id`
- Authentication and authorization via existing RBAC

### Notification System
- License expiry notifications
- Registration/insurance expiry notifications
- Maintenance due notifications
- Fuel efficiency anomaly alerts

### File Upload
- Driver license documents
- Vehicle registration documents
- Maintenance receipts
- Integration with existing upload repository

## Localization

Translation keys structure:
```
Fleet.NavigationLinks.*
Fleet.Dashboard.*
Fleet.Vehicles.*
Fleet.Drivers.*
Fleet.Trips.*
Fleet.Maintenance.*
Fleet.Fuel.*
Fleet.Enums.VehicleStatus.*
Fleet.Enums.FuelType.*
Fleet.Enums.ServiceType.*
```

All three locale files (en.json, ru.json, uz.json) must be maintained.
