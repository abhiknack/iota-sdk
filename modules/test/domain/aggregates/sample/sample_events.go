package sample

import (
	"context"

	"github.com/iota-uz/iota-sdk/pkg/composables"
	"github.com/iota-uz/iota-sdk/pkg/event"
)

const (
	SampleCreatedEvent = "sample.created"
	SampleUpdatedEvent = "sample.updated"
	SampleDeletedEvent = "sample.deleted"
)

func NewSampleCreatedEvent(ctx context.Context, entity Sample) (*event.Event, error) {
	userID, _ := composables.UseUserID(ctx)
	tenantID, _ := composables.UseTenantID(ctx)

	return &event.Event{
		Type:     SampleCreatedEvent,
		TenantID: tenantID,
		UserID:   userID,
		Payload:  entity,
	}, nil
}

func NewSampleUpdatedEvent(ctx context.Context, entity Sample) (*event.Event, error) {
	userID, _ := composables.UseUserID(ctx)
	tenantID, _ := composables.UseTenantID(ctx)

	return &event.Event{
		Type:     SampleUpdatedEvent,
		TenantID: tenantID,
		UserID:   userID,
		Payload:  entity,
	}, nil
}

func NewSampleDeletedEvent(ctx context.Context) (*event.Event, error) {
	userID, _ := composables.UseUserID(ctx)
	tenantID, _ := composables.UseTenantID(ctx)

	return &event.Event{
		Type:     SampleDeletedEvent,
		TenantID: tenantID,
		UserID:   userID,
	}, nil
}
