package trip

import (
	"context"

	"github.com/iota-uz/iota-sdk/modules/core/domain/aggregates/user"
	"github.com/iota-uz/iota-sdk/modules/core/domain/entities/session"
	"github.com/iota-uz/iota-sdk/pkg/composables"
)

func NewTripCreatedEvent(ctx context.Context, data Trip) (*TripCreatedEvent, error) {
	sender, err := composables.UseUser(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := composables.UseSession(ctx)
	if err != nil {
		return nil, err
	}
	return &TripCreatedEvent{
		Sender:  sender,
		Session: *sess,
		Data:    data,
	}, nil
}

func NewTripCompletedEvent(ctx context.Context, trip Trip) (*TripCompletedEvent, error) {
	sender, err := composables.UseUser(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := composables.UseSession(ctx)
	if err != nil {
		return nil, err
	}
	return &TripCompletedEvent{
		Sender:  sender,
		Session: *sess,
		Trip:    trip,
	}, nil
}

func NewTripCancelledEvent(ctx context.Context, trip Trip, reason string) (*TripCancelledEvent, error) {
	sender, err := composables.UseUser(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := composables.UseSession(ctx)
	if err != nil {
		return nil, err
	}
	return &TripCancelledEvent{
		Sender:  sender,
		Session: *sess,
		Trip:    trip,
		Reason:  reason,
	}, nil
}

type TripCreatedEvent struct {
	Sender  user.User
	Session session.Session
	Data    Trip
	Result  Trip
}

type TripCompletedEvent struct {
	Sender  user.User
	Session session.Session
	Trip    Trip
}

type TripCancelledEvent struct {
	Sender  user.User
	Session session.Session
	Trip    Trip
	Reason  string
}
