package driver

import (
	"time"

	"github.com/google/uuid"
	"github.com/iota-uz/iota-sdk/modules/fleet/domain/enums"
)

type Driver interface {
	ID() uuid.UUID
	TenantID() uuid.UUID
	UserID() *int64
	FirstName() string
	LastName() string
	LicenseNumber() string
	LicenseExpiry() time.Time
	Phone() string
	Email() string
	Status() enums.DriverStatus
	CreatedAt() time.Time
	UpdatedAt() time.Time

	UpdateLicense(number string, expiry time.Time) Driver
	UpdateContact(phone, email string) Driver
	UpdateStatus(status enums.DriverStatus) Driver
}

type DriverOption func(*driver)

func NewDriver(
	id uuid.UUID,
	tenantID uuid.UUID,
	firstName string,
	lastName string,
	licenseNumber string,
	licenseExpiry time.Time,
	opts ...DriverOption,
) Driver {
	d := &driver{
		id:            id,
		tenantID:      tenantID,
		userID:        nil,
		firstName:     firstName,
		lastName:      lastName,
		licenseNumber: licenseNumber,
		licenseExpiry: licenseExpiry,
		phone:         "",
		email:         "",
		status:        enums.DriverStatusActive,
		createdAt:     time.Now(),
		updatedAt:     time.Now(),
	}

	for _, opt := range opts {
		opt(d)
	}

	return d
}

func WithUserID(userID int64) DriverOption {
	return func(d *driver) {
		d.userID = &userID
	}
}

func WithPhone(phone string) DriverOption {
	return func(d *driver) {
		d.phone = phone
	}
}

func WithEmail(email string) DriverOption {
	return func(d *driver) {
		d.email = email
	}
}

func WithDriverStatus(status enums.DriverStatus) DriverOption {
	return func(d *driver) {
		d.status = status
	}
}

func WithDriverTimestamps(createdAt, updatedAt time.Time) DriverOption {
	return func(d *driver) {
		d.createdAt = createdAt
		d.updatedAt = updatedAt
	}
}

type driver struct {
	id            uuid.UUID
	tenantID      uuid.UUID
	userID        *int64
	firstName     string
	lastName      string
	licenseNumber string
	licenseExpiry time.Time
	phone         string
	email         string
	status        enums.DriverStatus
	createdAt     time.Time
	updatedAt     time.Time
}

func (d *driver) ID() uuid.UUID {
	return d.id
}

func (d *driver) TenantID() uuid.UUID {
	return d.tenantID
}

func (d *driver) UserID() *int64 {
	return d.userID
}

func (d *driver) FirstName() string {
	return d.firstName
}

func (d *driver) LastName() string {
	return d.lastName
}

func (d *driver) LicenseNumber() string {
	return d.licenseNumber
}

func (d *driver) LicenseExpiry() time.Time {
	return d.licenseExpiry
}

func (d *driver) Phone() string {
	return d.phone
}

func (d *driver) Email() string {
	return d.email
}

func (d *driver) Status() enums.DriverStatus {
	return d.status
}

func (d *driver) CreatedAt() time.Time {
	return d.createdAt
}

func (d *driver) UpdatedAt() time.Time {
	return d.updatedAt
}

func (d *driver) UpdateLicense(number string, expiry time.Time) Driver {
	return &driver{
		id:            d.id,
		tenantID:      d.tenantID,
		userID:        d.userID,
		firstName:     d.firstName,
		lastName:      d.lastName,
		licenseNumber: number,
		licenseExpiry: expiry,
		phone:         d.phone,
		email:         d.email,
		status:        d.status,
		createdAt:     d.createdAt,
		updatedAt:     time.Now(),
	}
}

func (d *driver) UpdateContact(phone, email string) Driver {
	return &driver{
		id:            d.id,
		tenantID:      d.tenantID,
		userID:        d.userID,
		firstName:     d.firstName,
		lastName:      d.lastName,
		licenseNumber: d.licenseNumber,
		licenseExpiry: d.licenseExpiry,
		phone:         phone,
		email:         email,
		status:        d.status,
		createdAt:     d.createdAt,
		updatedAt:     time.Now(),
	}
}

func (d *driver) UpdateStatus(status enums.DriverStatus) Driver {
	return &driver{
		id:            d.id,
		tenantID:      d.tenantID,
		userID:        d.userID,
		firstName:     d.firstName,
		lastName:      d.lastName,
		licenseNumber: d.licenseNumber,
		licenseExpiry: d.licenseExpiry,
		phone:         d.phone,
		email:         d.email,
		status:        status,
		createdAt:     d.createdAt,
		updatedAt:     time.Now(),
	}
}
