-- Migration: Add Fleet Management Permissions
-- Date: 2025-11-15
-- Purpose: Insert permissions for fleet management module

-- +migrate Up

-- Insert Fleet Management Permissions
INSERT INTO permissions (id, name, resource, action, modifier, description) VALUES
-- Vehicle Permissions
('a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d', 'Vehicle.Create', 'vehicle', 'create', 'all', 'Create new vehicles'),
('b2c3d4e5-f6a7-5b6c-9d0e-1f2a3b4c5d6e', 'Vehicle.Read', 'vehicle', 'read', 'all', 'View vehicle information'),
('c3d4e5f6-a7b8-6c7d-0e1f-2a3b4c5d6e7f', 'Vehicle.Update', 'vehicle', 'update', 'all', 'Update vehicle information'),
('d4e5f6a7-b8c9-7d8e-1f2a-3b4c5d6e7f8a', 'Vehicle.Delete', 'vehicle', 'delete', 'all', 'Delete vehicles'),

-- Driver Permissions
('e5f6a7b8-c9d0-8e9f-2a3b-4c5d6e7f8a9b', 'Driver.Create', 'driver', 'create', 'all', 'Create new drivers'),
('f6a7b8c9-d0e1-9f0a-3b4c-5d6e7f8a9b0c', 'Driver.Read', 'driver', 'read', 'all', 'View driver information'),
('a7b8c9d0-e1f2-0a1b-4c5d-6e7f8a9b0c1d', 'Driver.Update', 'driver', 'update', 'all', 'Update driver information'),
('b8c9d0e1-f2a3-1b2c-5d6e-7f8a9b0c1d2e', 'Driver.Delete', 'driver', 'delete', 'all', 'Delete drivers'),

-- Trip Permissions
('c9d0e1f2-a3b4-2c3d-6e7f-8a9b0c1d2e3f', 'Trip.Create', 'trip', 'create', 'all', 'Create new trips'),
('d0e1f2a3-b4c5-3d4e-7f8a-9b0c1d2e3f4a', 'Trip.Read', 'trip', 'read', 'all', 'View trip information'),
('e1f2a3b4-c5d6-4e5f-8a9b-0c1d2e3f4a5b', 'Trip.Update', 'trip', 'update', 'all', 'Update trip information'),
('f2a3b4c5-d6e7-5f6a-9b0c-1d2e3f4a5b6c', 'Trip.Delete', 'trip', 'delete', 'all', 'Delete trips'),

-- Maintenance Permissions
('a3b4c5d6-e7f8-6a7b-0c1d-2e3f4a5b6c7d', 'Maintenance.Create', 'maintenance', 'create', 'all', 'Create maintenance records'),
('b4c5d6e7-f8a9-7b8c-1d2e-3f4a5b6c7d8e', 'Maintenance.Read', 'maintenance', 'read', 'all', 'View maintenance records'),
('c5d6e7f8-a9b0-8c9d-2e3f-4a5b6c7d8e9f', 'Maintenance.Update', 'maintenance', 'update', 'all', 'Update maintenance records'),
('d6e7f8a9-b0c1-9d0e-3f4a-5b6c7d8e9f0a', 'Maintenance.Delete', 'maintenance', 'delete', 'all', 'Delete maintenance records'),

-- Fuel Entry Permissions
('e7f8a9b0-c1d2-0e1f-4a5b-6c7d8e9f0a1b', 'FuelEntry.Create', 'fuel_entry', 'create', 'all', 'Create fuel entries'),
('f8a9b0c1-d2e3-1f2a-5b6c-7d8e9f0a1b2c', 'FuelEntry.Read', 'fuel_entry', 'read', 'all', 'View fuel entries'),
('a9b0c1d2-e3f4-2a3b-6c7d-8e9f0a1b2c3d', 'FuelEntry.Update', 'fuel_entry', 'update', 'all', 'Update fuel entries'),
('b0c1d2e3-f4a5-3b4c-7d8e-9f0a1b2c3d4e', 'FuelEntry.Delete', 'fuel_entry', 'delete', 'all', 'Delete fuel entries')
ON CONFLICT (id) DO NOTHING;

-- +migrate Down

-- Remove Fleet Management Permissions
DELETE FROM permissions WHERE id IN (
    'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d',
    'b2c3d4e5-f6a7-5b6c-9d0e-1f2a3b4c5d6e',
    'c3d4e5f6-a7b8-6c7d-0e1f-2a3b4c5d6e7f',
    'd4e5f6a7-b8c9-7d8e-1f2a-3b4c5d6e7f8a',
    'e5f6a7b8-c9d0-8e9f-2a3b-4c5d6e7f8a9b',
    'f6a7b8c9-d0e1-9f0a-3b4c-5d6e7f8a9b0c',
    'a7b8c9d0-e1f2-0a1b-4c5d-6e7f8a9b0c1d',
    'b8c9d0e1-f2a3-1b2c-5d6e-7f8a9b0c1d2e',
    'c9d0e1f2-a3b4-2c3d-6e7f-8a9b0c1d2e3f',
    'd0e1f2a3-b4c5-3d4e-7f8a-9b0c1d2e3f4a',
    'e1f2a3b4-c5d6-4e5f-8a9b-0c1d2e3f4a5b',
    'f2a3b4c5-d6e7-5f6a-9b0c-1d2e3f4a5b6c',
    'a3b4c5d6-e7f8-6a7b-0c1d-2e3f4a5b6c7d',
    'b4c5d6e7-f8a9-7b8c-1d2e-3f4a5b6c7d8e',
    'c5d6e7f8-a9b0-8c9d-2e3f-4a5b6c7d8e9f',
    'd6e7f8a9-b0c1-9d0e-3f4a-5b6c7d8e9f0a',
    'e7f8a9b0-c1d2-0e1f-4a5b-6c7d8e9f0a1b',
    'f8a9b0c1-d2e3-1f2a-5b6c-7d8e9f0a1b2c',
    'a9b0c1d2-e3f4-2a3b-6c7d-8e9f0a1b2c3d',
    'b0c1d2e3-f4a5-3b4c-7d8e-9f0a1b2c3d4e'
);
