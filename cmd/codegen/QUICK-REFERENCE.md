# Code Generator Quick Reference Card

## Commands

### Generate Complete CRUD
```bash
./scripts/generate.sh crud -m MODULE -e ENTITY -f "FIELDS"
```

### Generate Entity Only
```bash
./scripts/generate.sh entity -m MODULE -e ENTITY -f "FIELDS"
```

### Generate Migration
```bash
./scripts/generate.sh migration
```

## Field Syntax

```
FieldName:Type:Validation
```

**Multiple fields:** Comma-separated
```
"Name:string:required,Age:int:min=0,Email:string:email"
```

## Types

| Type | Go Type |
|------|---------|
| `string` | `string` |
| `int` | `int` |
| `int64` | `int64` |
| `float64` | `float64` |
| `bool` | `bool` |
| `time.Time` | `time.Time` |
| `uuid.UUID` | `uuid.UUID` |

## Validations

| Tag | Description |
|-----|-------------|
| `required` | Field is required |
| `min=N` | Minimum value/length |
| `max=N` | Maximum value/length |
| `len=N` | Exact length |
| `email` | Email format |
| `url` | URL format |

**Multiple validations:**
```
"Age:int:required,min=0,max=150"
```

## Generated Structure

```
modules/{module}/
├── domain/aggregates/{entity}/
│   ├── {entity}.go
│   ├── {entity}_repository.go
│   └── {entity}_events.go
├── infrastructure/persistence/
│   └── {entity}_repository.go
├── services/
│   └── {entity}_service.go
└── presentation/controllers/
    ├── dtos/{entity}_dto.go
    └── {entity}_controller.go
```

## Post-Generation Checklist

- [ ] Create migration: `./scripts/generate.sh migration`
- [ ] Edit migration file with table structure
- [ ] Apply migration: `make db migrate up`
- [ ] Register service in `module.go`
- [ ] Register controller in `module.go`
- [ ] Add permissions to `permissions/constants.go`
- [ ] Create templates in `presentation/templates/pages/`
- [ ] Add translations to locale files
- [ ] Build: `templ generate && make css`
- [ ] Test: `go vet ./...`

## Quick Examples

### Simple Entity
```bash
./scripts/generate.sh crud -m crm -e Contact \
  -f "Name:string:required,Email:string:email"
```

### Entity with Validation
```bash
./scripts/generate.sh crud -m warehouse -e Product \
  -f "Name:string:required,SKU:string:required,Price:float64:min=0"
```

### Entity with Dates
```bash
./scripts/generate.sh crud -m hrm -e Employee \
  -f "FirstName:string:required,HireDate:time.Time:required"
```

### Entity with Relations
```bash
./scripts/generate.sh crud -m projects -e Task \
  -f "Title:string:required,ProjectID:uuid.UUID:required"
```

## Common Patterns

### Master-Detail
```bash
# Master
./scripts/generate.sh crud -m sales -e Order \
  -f "Number:string:required,Date:time.Time:required"

# Detail
./scripts/generate.sh crud -m sales -e OrderItem \
  -f "OrderID:uuid.UUID:required,Quantity:int:min=1"
```

### Hierarchical
```bash
./scripts/generate.sh crud -m org -e Department \
  -f "Name:string:required,ParentID:uuid.UUID"
```

## Naming Conventions

- **Module:** lowercase (e.g., `fleet`, `crm`, `warehouse`)
- **Entity:** PascalCase (e.g., `Vehicle`, `Contact`, `Product`)
- **Fields:** PascalCase (e.g., `FirstName`, `EmailAddress`)

## Tips

✅ **DO:**
- Use descriptive field names
- Include appropriate validations
- Review generated code
- Customize business logic
- Add tests

❌ **DON'T:**
- Use snake_case for fields
- Skip validation tags
- Commit without review
- Forget to register in module.go
- Skip documentation

## Help

```bash
./scripts/generate.sh help
```

## Full Documentation

- Quick Start: [CODEGEN.md](../../CODEGEN.md)
- Complete Guide: [docs/code-generation.md](../../docs/code-generation.md)
- Examples: [EXAMPLES.md](EXAMPLES.md)
- Generator README: [README.md](README.md)
