---
inclusion: always
---

# Configuration Management for IOTA SDK

## Scope

### In Scope
- Environment files: `.env*`, `e2e/.env.e2e`
- Docker: `compose*.yml`
- Build: `Makefile`, `tailwind.config.js`, `tsconfig*.json`
- Documentation: `README.MD`, `docs/`, `pkg/*/README.md`

### Out of Scope (use other steering)
- Go code changes
- Database migrations
- UI/template changes

## Configuration Types

### Environment Files
- `.env.example` - Template for environment variables
- `e2e/.env.e2e` - E2E testing environment
- Sync across dev/test/prod environments
- Never commit secrets

### Docker Configuration
- `compose.dev.yml` - Development services
- `compose.yml` - Production configuration
- `compose.testing.yml` - Testing environment
- Maintain service dependencies and networks

### Build Configuration
- `Makefile` - Build automation targets
- `tailwind.config.js` - UI framework config
- `tsconfig*.json` - TypeScript configuration
- `e2e/playwright.config.ts` - E2E testing config

## Update Triggers

### When to Update Environment
- New service integrations (Stripe, OAuth)
- Database configuration changes
- Port or host changes
- Feature flags
- Authentication updates

### When to Update Docker
- New services in architecture
- Database version upgrades
- Service dependency changes
- Port conflicts
- Volume mount requirements

### When to Update Makefile
- New build steps (CSS, assets)
- Database migration commands
- Testing infrastructure updates
- Code generation steps
- Linting/formatting tools

## Validation
```bash
make check tr              # Translation consistency
git status                 # Uncommitted changes
git diff                   # Review modifications
docker compose -f compose.dev.yml config  # Validate Docker
```
