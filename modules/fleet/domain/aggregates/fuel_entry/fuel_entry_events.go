package fuel_entry

import (
	"context"

	"github.com/iota-uz/iota-sdk/modules/core/domain/aggregates/user"
	"github.com/iota-uz/iota-sdk/modules/core/domain/entities/session"
	"github.com/iota-uz/iota-sdk/pkg/composables"
)

func NewFuelEntryCreatedEvent(ctx context.Context, data FuelEntry) (*FuelEntryCreatedEvent, error) {
	sender, err := composables.UseUser(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := composables.UseSession(ctx)
	if err != nil {
		return nil, err
	}
	return &FuelEntryCreatedEvent{
		Sender:  sender,
		Session: *sess,
		Data:    data,
	}, nil
}

func NewFuelEntryUpdatedEvent(ctx context.Context, data FuelEntry) (*FuelEntryUpdatedEvent, error) {
	sender, err := composables.UseUser(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := composables.UseSession(ctx)
	if err != nil {
		return nil, err
	}
	return &FuelEntryUpdatedEvent{
		Sender:  sender,
		Session: *sess,
		Data:    data,
	}, nil
}

func NewFuelEntryDeletedEvent(ctx context.Context) (*FuelEntryDeletedEvent, error) {
	sender, err := composables.UseUser(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := composables.UseSession(ctx)
	if err != nil {
		return nil, err
	}
	return &FuelEntryDeletedEvent{
		Sender:  sender,
		Session: *sess,
	}, nil
}

type FuelEntryCreatedEvent struct {
	Sender  user.User
	Session session.Session
	Data    FuelEntry
	Result  FuelEntry
}

type FuelEntryUpdatedEvent struct {
	Sender  user.User
	Session session.Session
	Data    FuelEntry
	Result  FuelEntry
}

type FuelEntryDeletedEvent struct {
	Sender  user.User
	Session session.Session
	Result  FuelEntry
}
