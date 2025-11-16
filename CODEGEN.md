# Code Generator Quick Start

Generate boilerplate code for new entities in IOTA SDK.

## Quick Commands

### Generate Complete CRUD

```bash
# Linux/Mac
./scripts/generate.sh crud -m MODULE -e ENTITY -f "Field1:type:validation,Field2:type:validation"

# Windows
scripts\generate.bat crud -m MODULE -e ENTITY -f "Field1:type:validation,Field2:type:validation"

# Direct Go command
go run cmd/codegen/main.go -type=crud -module=MODULE -entity=ENTITY -fields="Field1:type:validation"
```

### Generate Entity Only

```bash
./scripts/generate.sh entity -m MODULE -e ENTITY -f "Field1:type:validation"
```

### Generate Migration

```bash
./scripts/generate.sh migration
```

## Examples

### Example 1: Vehicle Entity

```bash
./scripts/generate.sh crud \
  -m fleet \
  -e Vehicle \
  -f "Make:string:required,Model:string:required,Year:int:min=1900,VIN:string:len=17"
```

Generates:
- Domain aggregate with Make, Model, Year, VIN fields
- Repository with CRUD operations
- Service with business logic
- Controller with HTTP endpoints
- DTOs with validation

### Example 2: Contact Entity

```bash
./scripts/generate.sh crud \
  -m crm \
  -e Contact \
  -f "FirstName:string:required,LastName:string:required,Email:string:email,Phone:string:max=20"
```

### Example 3: Product Entity

```bash
./scripts/generate.sh crud \
  -m warehouse \
  -e Product \
  -f "Name:string:required,SKU:string:required,Price:float64:min=0,Quantity:int:min=0"
```

## Field Types

| Type | Go Type | Example |
|------|---------|---------|
| `string` | `string` | `Name:string:required` |
| `int` | `int` | `Age:int:min=0` |
| `float64` | `float64` | `Price:float64:min=0` |
| `bool` | `bool` | `Active:bool` |
| `time.Time` | `time.Time` | `CreatedDate:time.Time:required` |
| `uuid.UUID` | `uuid.UUID` | `ParentID:uuid.UUID` |

## Validation Tags

- `required` - Field is required
- `min=N` - Minimum value/length
- `max=N` - Maximum value/length
- `len=N` - Exact length
- `email` - Email format
- `url` - URL format

Multiple validations: `Age:int:required,min=0,max=150`

## What Gets Generated

```
modules/{module}/
├── domain/aggregates/{entity}/
│   ├── {entity}.go              # Entity interface & implementation
│   ├── {entity}_repository.go   # Repository interface
│   └── {entity}_events.go       # Domain events
├── infrastructure/persistence/
│   └── {entity}_repository.go   # Repository implementation
├── services/
│   └── {entity}_service.go      # Service layer
└── presentation/controllers/
    ├── dtos/{entity}_dto.go     # DTOs
    └── {entity}_controller.go   # HTTP controller
```

## Post-Generation Steps

1. **Create migration**: `./scripts/generate.sh migration`
2. **Register in module.go**: Add service and controller registration
3. **Add permissions**: Define in `permissions/constants.go`
4. **Create templates**: Add `.templ` files for UI
5. **Add translations**: Update locale files (en.json, ru.json, uz.json)
6. **Build**: `templ generate && make css`
7. **Test**: `go vet ./...`

## Full Documentation

See [docs/code-generation.md](docs/code-generation.md) for complete documentation.

## Help

```bash
./scripts/generate.sh help
```

## Tips

- Use PascalCase for entity names: `Vehicle`, `FuelEntry`
- Use PascalCase for field names: `FirstName`, `LicenseNumber`
- Review and customize generated code before committing
- Add custom business logic to domain aggregates
- Add custom queries to repositories
- Add custom endpoints to controllers
