package driver

import (
	"context"

	"github.com/iota-uz/iota-sdk/modules/core/domain/aggregates/user"
	"github.com/iota-uz/iota-sdk/modules/core/domain/entities/session"
	"github.com/iota-uz/iota-sdk/pkg/composables"
)

func NewDriverCreatedEvent(ctx context.Context, data Driver) (*DriverCreatedEvent, error) {
	sender, err := composables.UseUser(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := composables.UseSession(ctx)
	if err != nil {
		return nil, err
	}
	return &DriverCreatedEvent{
		Sender:  sender,
		Session: *sess,
		Data:    data,
	}, nil
}

func NewDriverUpdatedEvent(ctx context.Context, data Driver) (*DriverUpdatedEvent, error) {
	sender, err := composables.UseUser(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := composables.UseSession(ctx)
	if err != nil {
		return nil, err
	}
	return &DriverUpdatedEvent{
		Sender:  sender,
		Session: *sess,
		Data:    data,
	}, nil
}

func NewDriverDeletedEvent(ctx context.Context) (*DriverDeletedEvent, error) {
	sender, err := composables.UseUser(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := composables.UseSession(ctx)
	if err != nil {
		return nil, err
	}
	return &DriverDeletedEvent{
		Sender:  sender,
		Session: *sess,
	}, nil
}

type DriverCreatedEvent struct {
	Sender  user.User
	Session session.Session
	Data    Driver
	Result  Driver
}

type DriverUpdatedEvent struct {
	Sender  user.User
	Session session.Session
	Data    Driver
	Result  Driver
}

type DriverDeletedEvent struct {
	Sender  user.User
	Session session.Session
	Result  Driver
}
