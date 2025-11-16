package vehicle

import (
	"context"

	"github.com/iota-uz/iota-sdk/modules/core/domain/aggregates/user"
	"github.com/iota-uz/iota-sdk/modules/core/domain/entities/session"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
	"github.com/iota-uz/iota-sdk/pkg/composables"
)

func NewVehicleCreatedEvent(ctx context.Context, data Vehicle) (*VehicleCreatedEvent, error) {
	sender, err := composables.UseUser(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := composables.UseSession(ctx)
	if err != nil {
		return nil, err
	}
	return &VehicleCreatedEvent{
		Sender:  sender,
		Session: *sess,
		Data:    data,
	}, nil
}

func NewVehicleUpdatedEvent(ctx context.Context, data Vehicle) (*VehicleUpdatedEvent, error) {
	sender, err := composables.UseUser(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := composables.UseSession(ctx)
	if err != nil {
		return nil, err
	}
	return &VehicleUpdatedEvent{
		Sender:  sender,
		Session: *sess,
		Data:    data,
	}, nil
}

func NewVehicleStatusChangedEvent(ctx context.Context, vehicle Vehicle, oldStatus, newStatus enums.VehicleStatus) (*VehicleStatusChangedEvent, error) {
	sender, err := composables.UseUser(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := composables.UseSession(ctx)
	if err != nil {
		return nil, err
	}
	return &VehicleStatusChangedEvent{
		Sender:    sender,
		Session:   *sess,
		Vehicle:   vehicle,
		OldStatus: oldStatus,
		NewStatus: newStatus,
	}, nil
}

func NewVehicleDeletedEvent(ctx context.Context) (*VehicleDeletedEvent, error) {
	sender, err := composables.UseUser(ctx)
	if err != nil {
		return nil, err
	}
	sess, err := composables.UseSession(ctx)
	if err != nil {
		return nil, err
	}
	return &VehicleDeletedEvent{
		Sender:  sender,
		Session: *sess,
	}, nil
}

type VehicleCreatedEvent struct {
	Sender  user.User
	Session session.Session
	Data    Vehicle
	Result  Vehicle
}

type VehicleUpdatedEvent struct {
	Sender  user.User
	Session session.Session
	Data    Vehicle
	Result  Vehicle
}

type VehicleStatusChangedEvent struct {
	Sender    user.User
	Session   session.Session
	Vehicle   Vehicle
	OldStatus enums.VehicleStatus
	NewStatus enums.VehicleStatus
}

type VehicleDeletedEvent struct {
	Sender  user.User
	Session session.Session
	Result  Vehicle
}
