package module_definition

import (
	"time"

	"github.com/google/uuid"
)

type ModuleDefinitionCreatedEvent struct {
	ModuleID    uuid.UUID
	TenantID    uuid.UUID
	Name        string
	DisplayName string
	OccurredAt  time.Time
}

type ModuleDefinitionUpdatedEvent struct {
	ModuleID    uuid.UUID
	TenantID    uuid.UUID
	Name        string
	DisplayName string
	OccurredAt  time.Time
}

type ModuleDefinitionGeneratedEvent struct {
	ModuleID   uuid.UUID
	TenantID   uuid.UUID
	Name       string
	OccurredAt time.Time
}

type ModuleDefinitionDeletedEvent struct {
	ModuleID   uuid.UUID
	TenantID   uuid.UUID
	OccurredAt time.Time
}
