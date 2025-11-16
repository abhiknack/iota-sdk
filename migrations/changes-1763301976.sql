-- Migration: Create Studio Module Tables
-- Date: 2025-01-15
-- Purpose: Add tables for Module Studio feature to enable visual module and entity creation

-- +migrate Up
CREATE TABLE studio_module_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(10) NOT NULL,
    status INT NOT NULL DEFAULT 0,
    entities JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ,
    UNIQUE(tenant_id, name)
);

CREATE INDEX idx_studio_module_definitions_tenant_id ON studio_module_definitions(tenant_id);
CREATE INDEX idx_studio_module_definitions_status ON studio_module_definitions(status);
CREATE INDEX idx_studio_module_definitions_deleted_at ON studio_module_definitions(deleted_at);
CREATE INDEX idx_studio_module_definitions_name ON studio_module_definitions(tenant_id, name) WHERE deleted_at IS NULL;

-- +migrate Down
DROP TABLE IF EXISTS studio_module_definitions;
