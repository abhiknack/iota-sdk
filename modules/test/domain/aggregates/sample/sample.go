package sample

import (
	"time"

	"github.com/google/uuid"
)

type Sample interface {
	ID() uuid.UUID
	TenantID() uuid.UUID
	Name() string
	Age() int
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

type SampleOption func(*sample)

func NewSample(
	id uuid.UUID,
	tenantID uuid.UUID,
	name string,
	age int,
	opts ...SampleOption,
) Sample {
	e := &sample{
		id:        id,
		tenantID:  tenantID,
		name: name,
		age: age,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func WithTimestamps(createdAt, updatedAt time.Time) SampleOption {
	return func(e *sample) {
		e.createdAt = createdAt
		e.updatedAt = updatedAt
	}
}

type sample struct {
	id        uuid.UUID
	tenantID  uuid.UUID
	name string
	age int
	createdAt time.Time
	updatedAt time.Time
}

func (e *sample) ID() uuid.UUID {
	return e.id
}

func (e *sample) TenantID() uuid.UUID {
	return e.tenantID
}

func (e *sample) Name() string {
	return e.name
}

func (e *sample) Age() int {
	return e.age
}

func (e *sample) CreatedAt() time.Time {
	return e.createdAt
}

func (e *sample) UpdatedAt() time.Time {
	return e.updatedAt
}
