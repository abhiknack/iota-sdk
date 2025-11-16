# Code Generation Guide

IOTA SDK includes a code generator to reduce boilerplate when creating new entities. This guide explains how to use it effectively.

## Quick Start

### Generate Complete CRUD

```bash
# Linux/Mac
./scripts/generate.sh crud -m fleet -e Vehicle -f "Make:string:required,Model:string:required,Year:int:min=1900"

# Windows
scripts\generate.bat crud -m fleet -e Vehicle -f "Make:string:required,Model:string:required,Year:int:min=1900"
```

This generates:
- Domain aggregate (interface, implementation, events)
- Repository (interface and implementation)
- Service layer
- Controller with CRUD endpoints
- DTOs (Create, Update, Filter)

### Generate Entity Only

```bash
./scripts/generate.sh entity -m crm -e Contact -f "Name:string:required,Email:string:email"
```

### Generate Migration

```bash
./scripts/generate.sh migration
```

## What Gets Generated

### Domain Layer

**Location**: `modules/{module}/domain/aggregates/{entity}/`

Files:
- `{entity}.go` - Entity interface and implementation with functional options
- `{entity}_repository.go` - Repository interface with CRUD methods
- `{entity}_events.go` - Domain events (Created, Updated, Deleted)

Example structure:
```go
type Vehicle interface {
    ID() uuid.UUID
    TenantID() uuid.UUID
    Make() string
    Model() string
    // ... other getters
}

type Repository interface {
    GetByID(ctx context.Context, id uuid.UUID) (Vehicle, error)
    GetPaginated(ctx context.Context, params *FindParams) ([]Vehicle, error)
    // ... other methods
}
```

### Infrastructure Layer

**Location**: `modules/{module}/infrastructure/persistence/`

Files:
- `{entity}_repository.go` - PostgreSQL repository implementation

Features:
- Multi-tenant isolation (automatic tenant_id filtering)
- Soft deletes (deleted_at column)
- Pagination support
- Search and filtering
- Query building with pkg/repo

### Service Layer

**Location**: `modules/{module}/services/`

Files:
- `{entity}_service.go` - Business logic orchestration

Features:
- RBAC permission checks
- Event publishing
- Transaction management
- Error handling

### Presentation Layer

**Location**: `modules/{module}/presentation/`

Files:
- `controllers/{entity}_controller.go` - HTTP handlers
- `controllers/dtos/{entity}_dto.go` - Data transfer objects

Features:
- HTMX support
- Form validation
- Pagination
- CRUD endpoints

## Field Definitions

### Syntax

```
FieldName:Type:Validation
```

### Supported Types

| Type | Go Type | Example |
|------|---------|---------|
| `string` | `string` | `Name:string:required` |
| `int` | `int` | `Age:int:min=0` |
| `int64` | `int64` | `Count:int64:min=0` |
| `float64` | `float64` | `Price:float64:min=0` |
| `bool` | `bool` | `Active:bool` |
| `time.Time` | `time.Time` | `BirthDate:time.Time:required` |
| `uuid.UUID` | `uuid.UUID` | `ParentID:uuid.UUID` |

### Validation Tags

Common validation tags for DTOs:

| Tag | Description | Example |
|-----|-------------|---------|
| `required` | Field is required | `Name:string:required` |
| `min=N` | Minimum value/length | `Age:int:min=0` |
| `max=N` | Maximum value/length | `Name:string:max=100` |
| `len=N` | Exact length | `VIN:string:len=17` |
| `email` | Email format | `Email:string:email` |
| `url` | URL format | `Website:string:url` |

Multiple validations:
```
Age:int:required,min=0,max=150
```

## Post-Generation Checklist

After generating code, complete these steps:

### 1. Create Database Migration

```bash
./scripts/generate.sh migration
```

Edit `migrations/changes-{timestamp}.sql`:

```sql
-- +migrate Up
CREATE TABLE fleet_vehicles (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    make VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_fleet_vehicles_tenant_id ON fleet_vehicles(tenant_id);
CREATE INDEX idx_fleet_vehicles_deleted_at ON fleet_vehicles(deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS fleet_vehicles;
```

Apply migration:
```bash
make db migrate up
```

### 2. Register in Module

Edit `modules/{module}/module.go`:

```go
func (m *Module) Register(app application.Application) error {
    // ... existing code ...

    // Register repository
    vehicleRepo := persistence.NewVehicleRepository()

    // Register service
    app.RegisterServices(
        services.NewVehicleService(vehicleRepo, app.EventPublisher()),
    )

    // Register controller
    app.RegisterControllers(
        controllers.NewVehicleController(app),
    )

    return nil
}
```

### 3. Add Permissions

Edit `modules/{module}/permissions/constants.go`:

```go
const (
    // Vehicle permissions
    VehicleRead   = "fleet.vehicle.read"
    VehicleCreate = "fleet.vehicle.create"
    VehicleUpdate = "fleet.vehicle.update"
    VehicleDelete = "fleet.vehicle.delete"
)
```

Register permissions in seed data (`modules/core/seed/seed_permissions.go`).

### 4. Create Templates

Create templ files in `modules/{module}/presentation/templates/pages/{entity}/`:

**list.templ**:
```templ
package vehicles

import (
    "github.com/iota-uz/iota-sdk/components/base/pagination"
)

type IndexPageProps struct {
    Vehicles        []VehicleViewModel
    PaginationState *pagination.State
}

templ IndexPage(props *IndexPageProps) {
    <div class="container">
        <h1>Vehicles</h1>
        @VehiclesTable(props)
    </div>
}
```

**new.templ**:
```templ
package vehicles

templ NewPage() {
    <div class="container">
        <h1>New Vehicle</h1>
        @VehicleForm(nil)
    </div>
}
```

**edit.templ**:
```templ
package vehicles

templ EditPage(vehicle *VehicleViewModel) {
    <div class="container">
        <h1>Edit Vehicle</h1>
        @VehicleForm(vehicle)
    </div>
}
```

### 5. Add Translations

Add to `modules/{module}/presentation/locales/en.json`:

```json
{
  "Fleet": {
    "NavigationLinks": {
      "Vehicles": "Vehicles"
    },
    "Vehicles": {
      "Meta": {
        "Title": "Vehicles"
      },
      "List": {
        "Title": "Vehicles",
        "New": "New Vehicle",
        "Make": "Make",
        "Model": "Model",
        "Year": "Year"
      },
      "Single": {
        "Edit": "Edit Vehicle",
        "Delete": "Delete Vehicle"
      }
    }
  }
}
```

Repeat for `ru.json` and `uz.json`.

### 6. Build

```bash
templ generate && make css
```

### 7. Test

```bash
go vet ./...
go test ./modules/{module}/...
```

## Real-World Examples

### Example 1: Simple Contact Entity

```bash
./scripts/generate.sh crud \
  -m crm \
  -e Contact \
  -f "FirstName:string:required,LastName:string:required,Email:string:email,Phone:string:max=20"
```

### Example 2: Product with Pricing

```bash
./scripts/generate.sh crud \
  -m warehouse \
  -e Product \
  -f "Name:string:required,SKU:string:required,Price:float64:min=0,Quantity:int:min=0,Description:string:max=500"
```

### Example 3: Employee with Dates

```bash
./scripts/generate.sh crud \
  -m hrm \
  -e Employee \
  -f "FirstName:string:required,LastName:string:required,HireDate:time.Time:required,Salary:float64:min=0,DepartmentID:uuid.UUID:required"
```

### Example 4: Document with Status

```bash
./scripts/generate.sh crud \
  -m documents \
  -e Invoice \
  -f "Number:string:required,Amount:float64:min=0,DueDate:time.Time:required,Status:int:required"
```

## Customization

The generated code is a starting point. Common customizations:

### Add Custom Domain Methods

```go
// In domain aggregate
func (v *vehicle) CalculateMileage() int {
    return v.currentOdometer - v.initialOdometer
}
```

### Add Custom Repository Queries

```go
// In repository
func (r *VehicleRepository) GetByVIN(ctx context.Context, vin string) (vehicle.Vehicle, error) {
    query := repo.Join(selectVehicleQuery, "WHERE vin = $1 AND deleted_at IS NULL")
    vehicles, err := r.queryVehicles(ctx, query, vin)
    // ...
}
```

### Add Business Logic to Service

```go
// In service
func (s *VehicleService) AssignDriver(ctx context.Context, vehicleID, driverID uuid.UUID) error {
    // Custom business logic
    vehicle, err := s.repo.GetByID(ctx, vehicleID)
    // ...
}
```

### Add Custom Controller Endpoints

```go
// In controller
func (c *VehicleController) AssignDriver(w http.ResponseWriter, r *http.Request) {
    // Custom endpoint logic
}
```

## Best Practices

1. **Use PascalCase** for entity and field names
2. **Add validation tags** appropriate for your use case
3. **Review generated code** before committing
4. **Customize as needed** - generated code is a starting point
5. **Follow DDD principles** - keep domain logic in domain layer
6. **Test thoroughly** - add unit and integration tests
7. **Document custom logic** - add comments for complex business rules

## Troubleshooting

### Import Errors

```bash
go mod tidy
```

### Template Compilation Errors

```bash
templ generate
```

### Database Errors

Check that:
- Migration was applied: `make db migrate status`
- Table name matches repository queries
- All columns exist in database

### Permission Errors

Ensure permissions are:
- Defined in `permissions/constants.go`
- Registered in seed data
- Assigned to appropriate roles

## Advanced Usage

### Generate Multiple Entities

```bash
# Generate multiple related entities
./scripts/generate.sh crud -m fleet -e Vehicle -f "Make:string:required,Model:string:required"
./scripts/generate.sh crud -m fleet -e Driver -f "FirstName:string:required,LastName:string:required"
./scripts/generate.sh crud -m fleet -e Trip -f "VehicleID:uuid.UUID:required,DriverID:uuid.UUID:required"
```

### Custom Generators

To add custom generators, create a new file in `cmd/codegen/generators/` and add the command to `cmd/codegen/main.go`.

## See Also

- [AGENTS.md](../AGENTS.md) - Module architecture guide
- [go-development.md](../.kiro/steering/go-development.md) - Go development standards
- [database-patterns.md](../.kiro/steering/database-patterns.md) - Database patterns
