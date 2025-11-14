---
inclusion: always
---

# E2E Testing Standards for IOTA SDK

## E2E Infrastructure

### Database Configuration
- **E2E Database**: `iota_erp_e2e` (separate from dev `iota_erp`)
- **Config**: `/e2e/.env.e2e`
- **Reset**: `make e2e reset` for clean state
- **Seed**: `make e2e seed` for test data

### Commands
```bash
make e2e test              # Run all E2E tests
make e2e reset             # Reset database
cd e2e && npx playwright test --ui    # Interactive debugging
cd e2e && npx playwright test --debug # Debug mode
```

## Critical Testing Patterns

### Database & State Management
- **Always start with `make e2e reset`** for clean state
- Use database fixtures in `beforeEach`
- Verify isolation (E2E tests don't affect main DB)
- Clear sessions in `afterEach`

### Timing & Alpine.js
- Wait for Alpine initialization after navigation
- Handle async HTMX/Alpine.js requests
- Use explicit waits for dynamic content
- Playwright auto-waits but may need explicit waits

### Form & Component Testing
- Use Playwright's file upload methods
- Intercept POST requests to verify FormData
- Validate hidden inputs and form association
- Test Alpine.js component state and reactivity

## Common Failure Patterns

### Database Issues
- Test pollution → Use `make e2e reset`
- Wrong database → Verify `.env.e2e`
- Missing data → Check fixtures

### Timing/Race Conditions
- Alpine.js not initialized → Add explicit waits
- HTMX incomplete → Use route interception
- Elements not visible → Use `waitFor({ state: 'visible' })`

### Infrastructure
- Port conflicts → `lsof -i :3201` (E2E server), `:5438` (E2E DB)
- Server not running → `make e2e dev`
- DB connection → Check PostgreSQL and E2E database
