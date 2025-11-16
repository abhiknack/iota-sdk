# IOTA SDK Code Generator

A code generation tool to reduce boilerplate when creating new entities in IOTA SDK modules.

## Features

- **Entity Generation**: Generate domain aggregates with interfaces and implementations
- **CRUD Scaffolding**: Generate complete CRUD stack (domain, repository, service, controller, DTOs)
- **Migration Templates**: Generate timestamped migration files
- **DDD Compliance**: Follows Domain-Driven Design patterns used in IOTA SDK

## Installation

The code generator is part of the IOTA SDK project. No separate installation needed.

## Usage

### Generate Complete CRUD

Generate a complete CRUD implementation for a new entity:

```bash
go run cmd/codegen/main.go -type=crud \
  -module=fleet \
  -entity=Vehicle \
  -fields="Make:string:required,Model:string:required,Year:int:min=1900,VIN:string:len=17"
```

This generates:
- Domain aggregate (`modules/{module}/domain/aggregates/{entity}/`)
  - Entity interface and implementation
  - Repository interface
  - Domain events
- Infrastructure repository (`modules/{module}/infrastructure/persistence/`)
- Service layer (`modules/{module}/services/`)
- Controller (`modules/{module}/presentation/controllers/`)
- DTOs (`modules/{module}/presentation/controllers/dtos/`)

### Generate Entity Only

Generate just the domain aggregate:

```bash
go run cmd/codegen/main.go -type=entity \
  -module=fleet \
  -entity=Driver \
  -fields="FirstName:string:required,LastName:string:required,LicenseNumber:string:required"
```

### Generate Migration

Generate a timestamped migration file template:

```bash
go run cmd/codegen/main.go -type=migration
```

This creates `migrations/changes-{timestamp}.sql` with a template structure.

## Field Syntax

Fields are specified as comma-separated values with the format:

```
FieldName:Type:Validation
```

### Examples

```bash
# String field with required validation
"Name:string:required"

# Integer with min/max validation
"Age:int:min=0,max=150"

# String with length validation
"VIN:string:len=17"

# Multiple fields
"Make:string:required,Model:string:required,Year:int:min=1900"
```

### Supported Types

- `string`
- `int`
- `int64`
- `float64`
- `bool`
- `time.Time`
- `uuid.UUID`

### Validation Tags

Common validation tags (used in DTOs):
- `required` - Field is required
- `min=N` - Minimum value/length
- `max=N` - Maximum value/length
- `len=N` - Exact length
- `email` - Email format
- `url` - URL format

## Post-Generation Steps

After generating code, you need to:

### 1. Create Database Migration

```bash
go run cmd/codegen/main.go -type=migration
```

Edit the generated migration file to match your entity structure.

### 2. Register in Module

Edit `modules/{module}/module.go`:

```go
// Register service
app.RegisterServices(
    services.NewVehicleService(vehicleRepo, eventPublisher),
)

// Register controller
app.RegisterControllers(
    controllers.NewVehicleController(app),
)
```

### 3. Add Permissions

Edit `modules/{module}/permissions/constants.go`:

```go
const (
    VehicleRead   = "fleet.vehicle.read"
    VehicleCreate = "fleet.vehicle.create"
    VehicleUpdate = "fleet.vehicle.update"
    VehicleDelete = "fleet.vehicle.delete"
)
```

### 4. Create Templates

Create templ files in `modules/{module}/presentation/templates/pages/{entity}/`:
- `list.templ` - List view
- `new.templ` - Create form
- `edit.templ` - Edit form

### 5. Add Translations

Add translations to all locale files in `modules/{module}/presentation/locales/`:
- `en.json`
- `ru.json`
- `uz.json`

### 6. Build

```bash
templ generate && make css
```

## Examples

### Example 1: Simple Entity

```bash
go run cmd/codegen/main.go -type=crud \
  -module=crm \
  -entity=Contact \
  -fields="Name:string:required,Email:string:email,Phone:string:max=20"
```

### Example 2: Complex Entity

```bash
go run cmd/codegen/main.go -type=crud \
  -module=warehouse \
  -entity=Product \
  -fields="Name:string:required,SKU:string:required,Price:float64:min=0,Quantity:int:min=0,Description:string:max=500"
```

### Example 3: Entity with Time Fields

```bash
go run cmd/codegen/main.go -type=crud \
  -module=hrm \
  -entity=Employee \
  -fields="FirstName:string:required,LastName:string:required,HireDate:time.Time:required,Salary:float64:min=0"
```

## Generated Code Structure

```
modules/{module}/
├── domain/
│   └── aggregates/{entity}/
│       ├── {entity}.go              # Entity interface & implementation
│       ├── {entity}_repository.go   # Repository interface
│       └── {entity}_events.go       # Domain events
├── infrastructure/
│   └── persistence/
│       └── {entity}_repository.go   # Repository implementation
├── services/
│   └── {entity}_service.go          # Service layer
└── presentation/
    └── controllers/
        ├── dtos/
        │   └── {entity}_dto.go      # DTOs
        └── {entity}_controller.go   # HTTP controller
```

## Customization

The generated code follows IOTA SDK patterns but may need customization:

1. **Add custom methods** to domain aggregates
2. **Add custom queries** to repositories
3. **Add business logic** to services
4. **Customize validation** in DTOs
5. **Add custom endpoints** to controllers

## Tips

- Use PascalCase for entity names (e.g., `Vehicle`, `FuelEntry`)
- Use PascalCase for field names (e.g., `FirstName`, `LicenseNumber`)
- Keep field names descriptive and consistent
- Add validation tags appropriate for your use case
- Review and customize generated code before committing

## Troubleshooting

### Import Errors

After generation, run:
```bash
go mod tidy
```

### Template Errors

If templates don't compile:
```bash
templ generate
```

### Compilation Errors

Check that:
- Module name matches existing module
- Field types are valid Go types
- All imports are correct

## Contributing

To add new generators:

1. Create a new file in `cmd/codegen/generators/`
2. Implement the generator function
3. Add the generator type to `main.go`
4. Update this README

## License

Part of IOTA SDK - see main project LICENSE
