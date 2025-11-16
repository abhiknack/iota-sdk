package module_definition

import (
	"time"

	"github.com/google/uuid"
)

type ModuleDefinition interface {
	ID() uuid.UUID
	TenantID() uuid.UUID
	Name() string
	DisplayName() string
	Description() string
	Icon() string
	Status() Status
	Entities() []EntityDefinition
	CreatedAt() time.Time
	UpdatedAt() time.Time

	UpdateDetails(displayName, description, icon string) ModuleDefinition
	UpdateStatus(status Status) ModuleDefinition
	AddEntity(entity EntityDefinition) ModuleDefinition
	RemoveEntity(entityID uuid.UUID) ModuleDefinition
}

type EntityDefinition struct {
	ID          uuid.UUID
	Name        string
	DisplayName string
	Fields      []FieldDefinition
	CreatedAt   time.Time
}

type FieldDefinition struct {
	ID         uuid.UUID
	Name       string
	Type       string
	Required   bool
	Validation string
	Order      int
}

type Status int

const (
	StatusDraft Status = iota
	StatusGenerated
	StatusActive
	StatusArchived
)

func (s Status) String() string {
	switch s {
	case StatusDraft:
		return "draft"
	case StatusGenerated:
		return "generated"
	case StatusActive:
		return "active"
	case StatusArchived:
		return "archived"
	default:
		return "unknown"
	}
}

type moduleDefinition struct {
	id          uuid.UUID
	tenantID    uuid.UUID
	name        string
	displayName string
	description string
	icon        string
	status      Status
	entities    []EntityDefinition
	createdAt   time.Time
	updatedAt   time.Time
}

func New(
	id uuid.UUID,
	tenantID uuid.UUID,
	name string,
	displayName string,
	description string,
	icon string,
) ModuleDefinition {
	return &moduleDefinition{
		id:          id,
		tenantID:    tenantID,
		name:        name,
		displayName: displayName,
		description: description,
		icon:        icon,
		status:      StatusDraft,
		entities:    []EntityDefinition{},
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}
}

func (m *moduleDefinition) ID() uuid.UUID                { return m.id }
func (m *moduleDefinition) TenantID() uuid.UUID          { return m.tenantID }
func (m *moduleDefinition) Name() string                 { return m.name }
func (m *moduleDefinition) DisplayName() string          { return m.displayName }
func (m *moduleDefinition) Description() string          { return m.description }
func (m *moduleDefinition) Icon() string                 { return m.icon }
func (m *moduleDefinition) Status() Status               { return m.status }
func (m *moduleDefinition) Entities() []EntityDefinition { return m.entities }
func (m *moduleDefinition) CreatedAt() time.Time         { return m.createdAt }
func (m *moduleDefinition) UpdatedAt() time.Time         { return m.updatedAt }

func (m *moduleDefinition) UpdateDetails(displayName, description, icon string) ModuleDefinition {
	m.displayName = displayName
	m.description = description
	m.icon = icon
	m.updatedAt = time.Now()
	return m
}

func (m *moduleDefinition) UpdateStatus(status Status) ModuleDefinition {
	m.status = status
	m.updatedAt = time.Now()
	return m
}

func (m *moduleDefinition) AddEntity(entity EntityDefinition) ModuleDefinition {
	m.entities = append(m.entities, entity)
	m.updatedAt = time.Now()
	return m
}

func (m *moduleDefinition) RemoveEntity(entityID uuid.UUID) ModuleDefinition {
	filtered := make([]EntityDefinition, 0)
	for _, e := range m.entities {
		if e.ID != entityID {
			filtered = append(filtered, e)
		}
	}
	m.entities = filtered
	m.updatedAt = time.Now()
	return m
}
