---
inclusion: always
---

# Go Development Standards for IOTA SDK

## Core Principles
- Write minimal, maintainable code with iteratively-grown tests
- Follow Domain-Driven Design (DDD) with clear layer separation
- Use dependency injection via `di.H` pattern
- Implement multi-tenant architecture with organization_id/tenant_id isolation

## Layer Architecture

### Domain Layer (Pure Business Logic)
- Aggregates as interfaces (not structs)
- Repository interfaces (no implementation)
- Domain events and value objects
- Enum types with validation

### Infrastructure Layer (External Integrations)
- Repository implementations (empty structs)
- Database access via composables
- Query constants at file top
- Use `pkg/repo` for SQL building - NEVER concatenate strings

### Service Layer (Business Orchestration)
- Service structs with repository interfaces
- Transaction management
- Business workflow coordination
- Validation before persistence

### Presentation Layer (HTTP/UI)
- Controllers with `di.H` injection
- ViewModels for templates
- HTMX response handling
- Use `pkg/htmx` for all HTMX operations

## Critical API Usage

### SQL Query Building (pkg/repo)
**MUST USE - NEVER concatenate strings**:
```go
const queryName = `SELECT * FROM table WHERE id = $1 AND organization_id = $2`
query := repo.Join(base, repo.JoinWhere(conditions...))
```

### HTMX Operations (pkg/htmx)
**NEVER access Hx-* headers directly**:
```go
htmx.IsHxRequest(r)
htmx.Redirect(w, "/path")
htmx.SetTrigger(w, "event", data)
```

### Composables for Context
```go
composables.UseForm[T](defaults, r)
composables.UsePageCtx(ctx)
composables.GetOrgID(ctx)
composables.UseTx(ctx)
```

## Testing with ITF Framework

### Modern Test Structure
```go
func TestController(t *testing.T) {
    suite := itf.NewSuiteBuilder(t).
        WithModules(modules.BuiltInModules...).
        AsAdmin().
        Build()
    
    controller := NewController()
    suite.Register(controller)
    
    suite.POST("/path").
        FormString("name", "Test").
        HTMX().
        Assert(t).
        ExpectOK().
        ExpectBodyContains("success")
}
```

### Test Execution
- Quick: `make test`
- Coverage: `make test coverage`
- Single: `go test -run ^TestName$ ./path`
- **PostgreSQL DB name ≤ 63 chars** - keep test names short

## Build Commands
- After .templ changes: `templ generate && make css`
- After Go changes: `go vet ./...` (NOT `go build`)
- Run tests: `make test` or `go test -v ./...`
- Check translations: `make check tr`

## Security Requirements
- **Template URLs**: Always use `templ.URL(dynamicURL)`
- **SQL**: Always parameterized queries with `$1, $2` placeholders
- **CSRF**: Include tokens in forms
- **Multi-tenant**: Always filter by organization_id/tenant_id

## Code Style
- Use `go fmt` for formatting
- Follow Go v1.23.2 idioms
- Table-driven tests with descriptive names
- Error handling with `pkg/serrors`
- NEVER read `*_templ.go` files (they're generated)
