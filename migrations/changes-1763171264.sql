-- Migration: Create fleet management tables
-- Date: 2024-11-15
-- Purpose: Create tables for fleet management module including vehicles, drivers, trips, maintenance, and fuel entries

-- +migrate Up

-- Vehicles table
CREATE TABLE fleet_vehicles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    make VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INT NOT NULL,
    vin VARCHAR(17) NOT NULL,
    license_plate VARCHAR(20) NOT NULL,
    status INT NOT NULL DEFAULT 0,
    current_odometer INT NOT NULL DEFAULT 0,
    registration_expiry TIMESTAMPTZ NOT NULL,
    insurance_expiry TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fleet_vehicles_tenant_vin_unique UNIQUE (tenant_id, vin)
);

CREATE INDEX idx_fleet_vehicles_tenant_id ON fleet_vehicles(tenant_id);
CREATE INDEX idx_fleet_vehicles_status ON fleet_vehicles(status);
CREATE INDEX idx_fleet_vehicles_deleted_at ON fleet_vehicles(deleted_at);
CREATE INDEX idx_fleet_vehicles_registration_expiry ON fleet_vehicles(registration_expiry);
CREATE INDEX idx_fleet_vehicles_insurance_expiry ON fleet_vehicles(insurance_expiry);

-- Drivers table
CREATE TABLE fleet_drivers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id INT8 REFERENCES users(id) ON DELETE SET NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    license_number VARCHAR(50) NOT NULL,
    license_expiry TIMESTAMPTZ NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255),
    status INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fleet_drivers_tenant_license_unique UNIQUE (tenant_id, license_number)
);

CREATE INDEX idx_fleet_drivers_tenant_id ON fleet_drivers(tenant_id);
CREATE INDEX idx_fleet_drivers_user_id ON fleet_drivers(user_id);
CREATE INDEX idx_fleet_drivers_status ON fleet_drivers(status);
CREATE INDEX idx_fleet_drivers_deleted_at ON fleet_drivers(deleted_at);
CREATE INDEX idx_fleet_drivers_license_expiry ON fleet_drivers(license_expiry);

-- Trips table
CREATE TABLE fleet_trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    vehicle_id UUID NOT NULL REFERENCES fleet_vehicles(id) ON DELETE CASCADE,
    driver_id UUID NOT NULL REFERENCES fleet_drivers(id) ON DELETE CASCADE,
    origin VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    purpose TEXT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    start_odometer INT NOT NULL,
    end_odometer INT,
    status INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_fleet_trips_tenant_id ON fleet_trips(tenant_id);
CREATE INDEX idx_fleet_trips_vehicle_id ON fleet_trips(vehicle_id);
CREATE INDEX idx_fleet_trips_driver_id ON fleet_trips(driver_id);
CREATE INDEX idx_fleet_trips_status ON fleet_trips(status);
CREATE INDEX idx_fleet_trips_start_time ON fleet_trips(start_time);
CREATE INDEX idx_fleet_trips_deleted_at ON fleet_trips(deleted_at);

-- Maintenance records table
CREATE TABLE fleet_maintenance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    vehicle_id UUID NOT NULL REFERENCES fleet_vehicles(id) ON DELETE CASCADE,
    service_type INT NOT NULL,
    service_date TIMESTAMPTZ NOT NULL,
    odometer INT NOT NULL,
    cost DECIMAL(10,2) NOT NULL DEFAULT 0,
    service_provider VARCHAR(255),
    description TEXT,
    next_service_due TIMESTAMPTZ,
    next_service_odometer INT,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_fleet_maintenance_tenant_id ON fleet_maintenance(tenant_id);
CREATE INDEX idx_fleet_maintenance_vehicle_id ON fleet_maintenance(vehicle_id);
CREATE INDEX idx_fleet_maintenance_service_date ON fleet_maintenance(service_date);
CREATE INDEX idx_fleet_maintenance_next_service_due ON fleet_maintenance(next_service_due);
CREATE INDEX idx_fleet_maintenance_deleted_at ON fleet_maintenance(deleted_at);

-- Fuel entries table
CREATE TABLE fleet_fuel_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    vehicle_id UUID NOT NULL REFERENCES fleet_vehicles(id) ON DELETE CASCADE,
    driver_id UUID REFERENCES fleet_drivers(id) ON DELETE SET NULL,
    date TIMESTAMPTZ NOT NULL,
    quantity DECIMAL(10,2) NOT NULL,
    cost DECIMAL(10,2) NOT NULL,
    odometer INT NOT NULL,
    fuel_type INT NOT NULL,
    location VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_fleet_fuel_entries_tenant_id ON fleet_fuel_entries(tenant_id);
CREATE INDEX idx_fleet_fuel_entries_vehicle_id ON fleet_fuel_entries(vehicle_id);
CREATE INDEX idx_fleet_fuel_entries_driver_id ON fleet_fuel_entries(driver_id);
CREATE INDEX idx_fleet_fuel_entries_date ON fleet_fuel_entries(date);
CREATE INDEX idx_fleet_fuel_entries_deleted_at ON fleet_fuel_entries(deleted_at);

-- +migrate Down

DROP TABLE IF EXISTS fleet_fuel_entries;
DROP TABLE IF EXISTS fleet_maintenance;
DROP TABLE IF EXISTS fleet_trips;
DROP TABLE IF EXISTS fleet_drivers;
DROP TABLE IF EXISTS fleet_vehicles;
