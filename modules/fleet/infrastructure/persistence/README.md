# Fleet Module - Database Seed Data

This directory contains seed data for testing the Fleet Management module.

## Quick Start

### Using Docker (Recommended)

If you're running the application with Docker Compose:

**Windows:**
```cmd
scripts\seed-fleet-data.bat
```

**Linux/Mac:**
```bash
chmod +x scripts/seed-fleet-data.sh
./scripts/seed-fleet-data.sh
```

### Manual Method

If you prefer to run the SQL directly:

```bash
# Using Docker
docker compose -f compose.dev.yml exec -T db psql -U postgres -d iota_erp < modules/fleet/infrastructure/persistence/seed_data.sql

# Using local PostgreSQL
PGPASSWORD=postgres psql -h localhost -U postgres -d iota_erp -f modules/fleet/infrastructure/persistence/seed_data.sql
```

## What Gets Created

The seed script creates sample data for testing:

### Vehicles (3)
- Toyota Camry 2022 (ABC-123)
- Honda Accord 2021 (XYZ-456)
- Ford F-150 2023 (DEF-789)

### Drivers (2)
- John Doe (License: DL123456)
- Jane Smith (License: DL789012)

### Trips (3)
- 2 completed trips
- 1 active trip

### Maintenance Records (3)
- Oil change
- Tire rotation
- Safety inspection

### Fuel Entries (5)
- Multiple fuel entries across different vehicles
- Includes cost and efficiency data

## Requirements

- The database must have at least one tenant
- The database must have at least one user
- All fleet tables must be created (run migrations first)

## Troubleshooting

### "No tenant found" Error
Run migrations first to create the tenants table:
```bash
make db migrate up
```

### "No user found" Error
Create a user account through the application first, or ensure the users table has at least one record.

### Connection Errors
Verify your database is running:
```bash
docker compose -f compose.dev.yml ps
```

## After Seeding

Once the seed data is created, you can:

1. Access the fleet dashboard: http://localhost:3200/fleet/dashboard
2. View vehicles: http://localhost:3200/fleet/vehicles
3. View drivers: http://localhost:3200/fleet/drivers
4. View trips: http://localhost:3200/fleet/trips
5. View maintenance: http://localhost:3200/fleet/maintenance
6. View fuel entries: http://localhost:3200/fleet/fuel

## Cleaning Up

To remove all seed data:

```sql
-- Run this in psql
DELETE FROM fleet_fuel_entries WHERE tenant_id IN (SELECT id FROM tenants LIMIT 1);
DELETE FROM fleet_maintenance WHERE tenant_id IN (SELECT id FROM tenants LIMIT 1);
DELETE FROM fleet_trips WHERE tenant_id IN (SELECT id FROM tenants LIMIT 1);
DELETE FROM fleet_drivers WHERE tenant_id IN (SELECT id FROM tenants LIMIT 1);
DELETE FROM fleet_vehicles WHERE tenant_id IN (SELECT id FROM tenants LIMIT 1);
```

Or simply reset the entire database:
```bash
make db migrate down
make db migrate up
```
