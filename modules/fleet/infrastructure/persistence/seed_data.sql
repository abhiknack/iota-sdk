-- Fleet Module Seed Data
-- This script creates sample data for testing the fleet management module
-- Run with: psql -h localhost -U postgres -d iota_erp -f modules/fleet/infrastructure/persistence/seed_data.sql

-- Note: Replace the tenant_id and user_id values with actual IDs from your database
-- You can find tenant_id with: SELECT id FROM tenants LIMIT 1;
-- You can find user_id with: SELECT id FROM users LIMIT 1;

DO $$
DECLARE
    v_tenant_id UUID;
    v_user_id INT8;
    v_vehicle1_id UUID;
    v_vehicle2_id UUID;
    v_vehicle3_id UUID;
    v_driver1_id UUID;
    v_driver2_id UUID;
    v_trip1_id UUID;
    v_trip2_id UUID;
BEGIN
    -- Get first tenant and user
    SELECT id INTO v_tenant_id FROM tenants LIMIT 1;
    SELECT id INTO v_user_id FROM users LIMIT 1;
    
    IF v_tenant_id IS NULL THEN
        RAISE EXCEPTION 'No tenant found. Please create a tenant first.';
    END IF;
    
    IF v_user_id IS NULL THEN
        RAISE EXCEPTION 'No user found. Please create a user first.';
    END IF;
    
    RAISE NOTICE 'Using tenant_id: %', v_tenant_id;
    RAISE NOTICE 'Using user_id: %', v_user_id;
    
    -- Insert sample vehicles one by one to capture IDs
    INSERT INTO fleet_vehicles (id, tenant_id, make, model, year, vin, license_plate, status, current_odometer, registration_expiry, insurance_expiry, created_at, updated_at)
    VALUES (gen_random_uuid(), v_tenant_id, 'Toyota', 'Camry', 2022, '1HGBH41JXMN109186', 'ABC-123', 0, 15000, NOW() + INTERVAL '6 months', NOW() + INTERVAL '8 months', NOW(), NOW())
    RETURNING id INTO v_vehicle1_id;
    
    INSERT INTO fleet_vehicles (id, tenant_id, make, model, year, vin, license_plate, status, current_odometer, registration_expiry, insurance_expiry, created_at, updated_at)
    VALUES (gen_random_uuid(), v_tenant_id, 'Honda', 'Accord', 2021, '2HGBH41JXMN109187', 'XYZ-456', 0, 22000, NOW() + INTERVAL '4 months', NOW() + INTERVAL '5 months', NOW(), NOW())
    RETURNING id INTO v_vehicle2_id;
    
    INSERT INTO fleet_vehicles (id, tenant_id, make, model, year, vin, license_plate, status, current_odometer, registration_expiry, insurance_expiry, created_at, updated_at)
    VALUES (gen_random_uuid(), v_tenant_id, 'Ford', 'F-150', 2023, '3HGBH41JXMN109188', 'DEF-789', 0, 8000, NOW() + INTERVAL '10 months', NOW() + INTERVAL '11 months', NOW(), NOW())
    RETURNING id INTO v_vehicle3_id;
    
    RAISE NOTICE 'Created 3 vehicles';
    
    -- Insert sample drivers one by one to capture IDs
    INSERT INTO fleet_drivers (id, tenant_id, user_id, first_name, last_name, license_number, license_expiry, phone, email, status, created_at, updated_at)
    VALUES (gen_random_uuid(), v_tenant_id, v_user_id, 'John', 'Doe', 'DL123456', NOW() + INTERVAL '1 year', '+1234567890', 'john.doe@example.com', 0, NOW(), NOW())
    RETURNING id INTO v_driver1_id;
    
    INSERT INTO fleet_drivers (id, tenant_id, user_id, first_name, last_name, license_number, license_expiry, phone, email, status, created_at, updated_at)
    VALUES (gen_random_uuid(), v_tenant_id, NULL, 'Jane', 'Smith', 'DL789012', NOW() + INTERVAL '2 years', '+1234567891', 'jane.smith@example.com', 0, NOW(), NOW())
    RETURNING id INTO v_driver2_id;
    
    RAISE NOTICE 'Created 2 drivers';
    
    -- Insert sample trips one by one
    INSERT INTO fleet_trips (id, tenant_id, vehicle_id, driver_id, origin, destination, purpose, start_time, end_time, start_odometer, end_odometer, status, created_at, updated_at)
    VALUES (gen_random_uuid(), v_tenant_id, v_vehicle1_id, v_driver1_id, 'Office', 'Client Site A', 'Client Meeting', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days' + INTERVAL '3 hours', 15000, 15120, 2, NOW(), NOW())
    RETURNING id INTO v_trip1_id;
    
    INSERT INTO fleet_trips (id, tenant_id, vehicle_id, driver_id, origin, destination, purpose, start_time, end_time, start_odometer, end_odometer, status, created_at, updated_at)
    VALUES (gen_random_uuid(), v_tenant_id, v_vehicle2_id, v_driver2_id, 'Warehouse', 'Delivery Point B', 'Delivery', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days' + INTERVAL '2 hours', 22000, 22080, 2, NOW(), NOW())
    RETURNING id INTO v_trip2_id;
    
    INSERT INTO fleet_trips (id, tenant_id, vehicle_id, driver_id, origin, destination, purpose, start_time, end_time, start_odometer, end_odometer, status, created_at, updated_at)
    VALUES (gen_random_uuid(), v_tenant_id, v_vehicle1_id, v_driver1_id, 'Office', 'Airport', 'Airport Transfer', NOW() - INTERVAL '1 day', NULL, 15120, NULL, 1, NOW(), NOW());
    
    RAISE NOTICE 'Created 3 trips';
    
    -- Insert sample maintenance records
    INSERT INTO fleet_maintenance (tenant_id, vehicle_id, service_type, service_date, odometer, cost, service_provider, description, next_service_due, next_service_odometer, created_at, updated_at)
    VALUES 
        (v_tenant_id, v_vehicle1_id, 0, NOW() - INTERVAL '30 days', 14500, 45.00, 'Quick Lube', 'Regular oil change', NOW() + INTERVAL '90 days', 18500, NOW(), NOW()),
        (v_tenant_id, v_vehicle2_id, 1, NOW() - INTERVAL '60 days', 21000, 80.00, 'Tire Shop', 'Tire rotation and balance', NOW() + INTERVAL '120 days', 27000, NOW(), NOW()),
        (v_tenant_id, v_vehicle3_id, 3, NOW() - INTERVAL '15 days', 7800, 150.00, 'State Inspection', 'Annual safety inspection', NOW() + INTERVAL '350 days', NULL, NOW(), NOW());
    
    RAISE NOTICE 'Created 3 maintenance records';
    
    -- Insert sample fuel entries
    INSERT INTO fleet_fuel_entries (tenant_id, vehicle_id, driver_id, date, quantity, cost, odometer, fuel_type, location, created_at, updated_at)
    VALUES 
        (v_tenant_id, v_vehicle1_id, v_driver1_id, NOW() - INTERVAL '7 days', 45.5, 150.00, 14800, 0, 'Shell Station - Main St', NOW(), NOW()),
        (v_tenant_id, v_vehicle1_id, v_driver1_id, NOW() - INTERVAL '3 days', 42.0, 138.00, 15050, 0, 'BP Station - Highway 101', NOW(), NOW()),
        (v_tenant_id, v_vehicle2_id, v_driver2_id, NOW() - INTERVAL '5 days', 50.0, 165.00, 21800, 0, 'Chevron - Downtown', NOW(), NOW()),
        (v_tenant_id, v_vehicle2_id, v_driver2_id, NOW() - INTERVAL '2 days', 48.5, 159.00, 22050, 0, 'Shell Station - Main St', NOW(), NOW()),
        (v_tenant_id, v_vehicle3_id, NULL, NOW() - INTERVAL '4 days', 65.0, 214.00, 7900, 0, 'Costco Gas', NOW(), NOW());
    
    RAISE NOTICE 'Created 5 fuel entries';
    
    RAISE NOTICE 'Seed data created successfully!';
    RAISE NOTICE 'You can now access the fleet dashboard at /fleet/dashboard';
END $$;
