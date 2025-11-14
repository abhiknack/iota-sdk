---
inclusion: always
---

# Debugging Standards for IOTA SDK

## Debugging Workflow

### Phase 1: Triage
1. Thoroughly analyze error messages, logs, and context
2. Make three hypotheses on root cause
3. Find relevant code/tests/commands to confirm or refute each

### Phase 2: Analysis
1. Run "imaginary interpreter" to simulate dataflow
2. Trace each code path from start to finish
   - Controller → Service → Repository → SQL

### Phase 3: Solution
Determine most likely root cause and provide detailed fix with:
- Location (file:line)
- Technical explanation
- Required changes (FROM/TO)
- Verification command

## Critical Patterns

### Tenant Safety
```go
// WRONG: SELECT * FROM loads WHERE id = $1
// RIGHT: repo.NewQuery().Select("*").From("loads")
//        .Where("id = ?", id).Where("organization_id = ?", orgID)
```

### Common Issues
- Missing organization_id in queries
- Unchecked type assertions
- Missing error handling
- Panic in request handlers

## Commands
```bash
go test -v ./path -run TestName [-race]
go vet ./...
make db migrate up
git diff HEAD~1
git blame -L start,end file
```
