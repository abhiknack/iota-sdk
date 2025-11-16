@echo off
REM Fleet Module Seed Data Script for Windows
REM This script seeds the database with sample fleet data for testing

echo.
echo Fleet Module - Seed Data Script
echo ==================================
echo.

REM Check if Docker is available
docker compose ps >nul 2>&1
if %errorlevel% equ 0 (
    echo Using Docker Compose database...
    docker compose -f compose.dev.yml exec -T db psql -U postgres -d iota_erp < modules\fleet\infrastructure\persistence\seed_data.sql
    if %errorlevel% equ 0 (
        echo.
        echo Seed data created successfully!
        echo.
        echo You can now:
        echo   - Access the fleet dashboard at http://localhost:3200/fleet/dashboard
        echo   - View vehicles at http://localhost:3200/fleet/vehicles
        echo   - View drivers at http://localhost:3200/fleet/drivers
        echo   - View trips at http://localhost:3200/fleet/trips
        echo.
    ) else (
        echo.
        echo Error: Failed to seed database
        exit /b 1
    )
) else (
    echo Error: Docker Compose is not running
    echo Please start Docker Compose first with: docker compose -f compose.dev.yml up -d
    exit /b 1
)
