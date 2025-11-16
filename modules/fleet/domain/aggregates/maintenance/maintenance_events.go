package maintenance

import (
	"context"

	"github.com/iota-uz/iota-sdk/modules/core/domain/aggregates/user"
	"github.com/iota-uz/iota-sdk/modules/core/domain/entities/session"
	"github.com/iota-uz/iota-sdk/pkg/composables"
)

func NewMaintenanceCreatedEvent(ctx context.Context, data Maintenance) (*MaintenanceCreatedEvent, error) {
	sender, err := composables.UseUser(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := composables.UseSession(ctx)
	if err != nil {
		return nil, err
	}
	return &MaintenanceCreatedEvent{
		Sender:  sender,
		Session: *sess,
		Data:    data,
	}, nil
}

func NewMaintenanceUpdatedEvent(ctx context.Context, data Maintenance) (*MaintenanceUpdatedEvent, error) {
	sender, err := composables.UseUser(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := composables.UseSession(ctx)
	if err != nil {
		return nil, err
	}
	return &MaintenanceUpdatedEvent{
		Sender:  sender,
		Session: *sess,
		Data:    data,
	}, nil
}

func NewMaintenanceDeletedEvent(ctx context.Context) (*MaintenanceDeletedEvent, error) {
	sender, err := composables.UseUser(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := composables.UseSession(ctx)
	if err != nil {
		return nil, err
	}
	return &MaintenanceDeletedEvent{
		Sender:  sender,
		Session: *sess,
	}, nil
}

type MaintenanceCreatedEvent struct {
	Sender  user.User
	Session session.Session
	Data    Maintenance
	Result  Maintenance
}

type MaintenanceUpdatedEvent struct {
	Sender  user.User
	Session session.Session
	Data    Maintenance
	Result  Maintenance
}

type MaintenanceDeletedEvent struct {
	Sender  user.User
	Session session.Session
	Result  Maintenance
}
