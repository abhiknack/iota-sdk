package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const migrationTemplate = `-- Migration: [Brief description]
-- Date: %s
-- Purpose: [Detailed explanation]

-- +migrate Up
CREATE TABLE IF NOT EXISTS table_name (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_table_name_tenant_id ON table_name(tenant_id);
CREATE INDEX idx_table_name_deleted_at ON table_name(deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS table_name;
`

func GenerateMigration() error {
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("changes-%d.sql", timestamp)
	filepath := filepath.Join("migrations", filename)

	content := fmt.Sprintf(migrationTemplate, time.Now().Format("2006-01-02"))

	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}

	fmt.Printf("Created migration file: %s\n", filepath)
	return nil
}
