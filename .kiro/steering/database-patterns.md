---
inclusion: always
---

# Database Patterns for IOTA SDK

## Critical Rules
1. **NEVER edit existing migration files** - immutable once created
2. **ALWAYS include tenant_id** for multi-tenant isolation (except system tables)
3. **ALWAYS provide Down migrations** that fully reverse Up changes
4. **ALWAYS use Unix timestamp** in filename: `migrations/changes-{timestamp}.sql`
5. **NEVER use raw SQL in application code** - all schema changes via migrations
6. **NEVER use anonymous code blocks (DO $ ... $)** in migrations
7. **NEVER use BEGIN/COMMIT/ROLLBACK** in migrations - handled by migration tool

## Migration Template
```sql
-- Migration: [Brief description]
-- Date: YYYY-MM-DD
-- Purpose: [Detailed explanation]

-- +migrate Up
CREATE TABLE entities (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_entities_tenant_id ON entities(tenant_id);

-- +migrate Down
DROP TABLE IF EXISTS entities;
```

## Multi-Tenant Patterns
```sql
-- All queries must filter by tenant_id
SELECT * FROM users WHERE tenant_id = $1;
SELECT * FROM products WHERE tenant_id = $1 AND id = $2;
```

## Standard Table Structure
- `id uuid PRIMARY KEY DEFAULT gen_random_uuid()`
- `tenant_id uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE`
- Business fields
- Audit fields: `created_at`, `updated_at`, `created_by`, `updated_by`, `deleted_at`
- Indexes on `tenant_id`, `deleted_at`, and common query fields

## Connection Management
```bash
# Local development
PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d iota_erp

# Staging (Railway)
PGPASSWORD=A6E4g1d2ae43Bebg2F65CEc3e56aa25g psql -h shuttle.proxy.rlwy.net -U postgres -p 31150 -d railway
```

## Migration Commands
- Apply: `make db migrate up`
- Rollback: `make db migrate down`
- Status: `make db migrate status`
- Create: Generate timestamp with `date +%s`

## ITF Testing
- `itf.Setup(t)` provides isolated test database
- Each test gets clean database state
- Use repository methods, not raw SQL
- Automatic cleanup after test completion
