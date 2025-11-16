package dtos

import "github.com/google/uuid"

type ModuleDefinitionCreateDTO struct {
	Name        string `form:"Name" validate:"required,min=2,max=50"`
	DisplayName string `form:"DisplayName" validate:"required,min=2,max=100"`
	Description string `form:"Description" validate:"max=500"`
	Icon        string `form:"Icon" validate:"required"`
}

type ModuleDefinitionUpdateDTO struct {
	DisplayName string `form:"DisplayName" validate:"required,min=2,max=100"`
	Description string `form:"Description" validate:"max=500"`
	Icon        string `form:"Icon" validate:"required"`
}

type EntityDefinitionDTO struct {
	ID          uuid.UUID            `json:"id"`
	Name        string               `json:"name" validate:"required,min=2,max=50"`
	DisplayName string               `json:"display_name" validate:"required,min=2,max=100"`
	Fields      []FieldDefinitionDTO `json:"fields"`
}

type FieldDefinitionDTO struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name" validate:"required,min=2,max=50"`
	Type       string    `json:"type" validate:"required,oneof=string int int64 float64 bool time.Time uuid.UUID"`
	Required   bool      `json:"required"`
	Validation string    `json:"validation"`
	Order      int       `json:"order"`
}

type GenerateModuleDTO struct {
	ModuleID uuid.UUID `form:"ModuleID" validate:"required"`
}
