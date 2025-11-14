---
inclusion: always
---

# Refactoring Standards for IOTA SDK

## Issue Priority Classification

### Critical Issues ❌ (Fix Immediately)
- **SQL Injection**: String concatenation in SQL queries
- **XSS Vulnerabilities**: Missing `templ.URL()` for dynamic URLs
- **Raw HTML**: Using `@templ.Raw()` with user input
- **Runtime Failures**: `panic()` in handlers, unchecked type assertions
- **Data Integrity**: Raw SQL in tests, test names > 63 chars, missing organization_id

### Minor Issues 🟡 (Important)
- **API Misuse**: Direct HTMX headers, manual ID parsing, string literals for HTTP methods
- **DRY Violations**: Repeated business logic, copy-pasted validation, duplicate error handling
- **Hard-coded Values**: Status strings, magic numbers, repeated error messages
- **Code Smells**: Unused variables/functions/imports, long functions

### Style Issues 🟢 (Best Practices)
- Missing import aliases with `sdk` prefix
- Missing `t.Parallel()` in tests
- Inconsistent error wrapping

## Pattern Application

### SQL Query Management
```go
// BEFORE (Wrong)
query := "SELECT * FROM users WHERE org_id = " + orgID

// AFTER (Correct)
const userListQuery = `SELECT * FROM users WHERE organization_id = $1`
query := repo.Join(userListQuery, repo.JoinWhere(conditions...))
```

### HTMX Handling
```go
// BEFORE (Wrong)
if r.Header.Get("Hx-Request") == "true" {
    w.Header().Add("Hx-Redirect", "/path")
}

// AFTER (Correct)
if htmx.IsHxRequest(r) {
    htmx.Redirect(w, "/path")
}
```

### DRY Violation Fix
```go
// BEFORE (Duplicated)
if status == "Active" || status == "Pending" {
    // logic in 3 places
}

// AFTER (Extracted)
func isProcessableStatus(status LoadStatus) bool {
    return status == LOAD_STATUS_ACTIVE || status == LOAD_STATUS_PENDING
}
```

## Validation After Fixes

```bash
go vet ./...                    # After Go changes
make check tr                   # After translations
git diff --check                # Check whitespace
```

## Refactoring Workflow

1. **Identify scope**: Changed files or specified targets
2. **Scan for issues**: Categorize by priority
3. **Fix systematically**: Critical → Minor → Style
4. **Validate**: Run checks after each batch
5. **Report**: Detailed summary of changes

## Security Checklist

- No raw SQL concatenation
- All queries parameterized
- Tenant isolation verified
- No credentials in code
- SQL injection prevention verified
- XSS protection in templates
