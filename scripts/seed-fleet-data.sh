#!/bin/bash

# Fleet Module Seed Data Script
# This script seeds the database with sample fleet data for testing

set -e

echo "🚗 Fleet Module - Seed Data Script"
echo "=================================="
echo ""

# Check if running in Docker
if [ -f /.dockerenv ]; then
    echo "Running inside Docker container..."
    DB_HOST="db"
    DB_PORT="5432"
    DB_USER="postgres"
    DB_NAME="iota_erp"
    PGPASSWORD="postgres"
else
    echo "Running on host machine..."
    # Check if we should use Docker
    if command -v docker &> /dev/null && docker compose ps | grep -q "db"; then
        echo "Using Docker Compose database..."
        docker compose -f compose.dev.yml exec -T db psql -U postgres -d iota_erp -f - < modules/fleet/infrastructure/persistence/seed_data.sql
        exit 0
    else
        # Use local PostgreSQL
        DB_HOST="${DB_HOST:-localhost}"
        DB_PORT="${DB_PORT:-5432}"
        DB_USER="${DB_USER:-postgres}"
        DB_NAME="${DB_NAME:-iota_erp}"
        PGPASSWORD="${PGPASSWORD:-postgres}"
    fi
fi

export PGPASSWORD

echo "Database: $DB_NAME"
echo "Host: $DB_HOST:$DB_PORT"
echo "User: $DB_USER"
echo ""

# Check if database is accessible
if ! psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1" > /dev/null 2>&1; then
    echo "❌ Error: Cannot connect to database"
    echo "Please ensure PostgreSQL is running and credentials are correct"
    exit 1
fi

echo "✅ Database connection successful"
echo ""

# Run the seed script
echo "📝 Running seed script..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f modules/fleet/infrastructure/persistence/seed_data.sql

echo ""
echo "✅ Seed data created successfully!"
echo ""
echo "You can now:"
echo "  - Access the fleet dashboard at http://localhost:3200/fleet/dashboard"
echo "  - View vehicles at http://localhost:3200/fleet/vehicles"
echo "  - View drivers at http://localhost:3200/fleet/drivers"
echo "  - View trips at http://localhost:3200/fleet/trips"
echo ""
