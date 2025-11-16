package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ModuleDefinition struct {
	ID          uuid.UUID    `db:"id"`
	TenantID    uuid.UUID    `db:"tenant_id"`
	Name        string       `db:"name"`
	DisplayName string       `db:"display_name"`
	Description string       `db:"description"`
	Icon        string       `db:"icon"`
	Status      int          `db:"status"`
	Entities    EntitiesJSON `db:"entities"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
	DeletedAt   *time.Time   `db:"deleted_at"`
}

type EntityDefinitionJSON struct {
	ID          uuid.UUID             `json:"id"`
	Name        string                `json:"name"`
	DisplayName string                `json:"display_name"`
	Fields      []FieldDefinitionJSON `json:"fields"`
	CreatedAt   time.Time             `json:"created_at"`
}

type FieldDefinitionJSON struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Required   bool      `json:"required"`
	Validation string    `json:"validation"`
	Order      int       `json:"order"`
}

type EntitiesJSON []EntityDefinitionJSON

func (e EntitiesJSON) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *EntitiesJSON) Scan(value interface{}) error {
	if value == nil {
		*e = EntitiesJSON{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, e)
}
