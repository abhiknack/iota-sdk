package enums

import "fmt"

type DriverStatus int

const (
	DriverStatusActive DriverStatus = iota
	DriverStatusInactive
	DriverStatusOnLeave
	DriverStatusTerminated
)

func (d DriverStatus) String() string {
	switch d {
	case DriverStatusActive:
		return "Active"
	case DriverStatusInactive:
		return "Inactive"
	case DriverStatusOnLeave:
		return "OnLeave"
	case DriverStatusTerminated:
		return "Terminated"
	default:
		return "Unknown"
	}
}

func (d DriverStatus) IsValid() bool {
	return d >= DriverStatusActive && d <= DriverStatusTerminated
}

func ParseDriverStatus(s string) (DriverStatus, error) {
	switch s {
	case "Active":
		return DriverStatusActive, nil
	case "Inactive":
		return DriverStatusInactive, nil
	case "OnLeave":
		return DriverStatusOnLeave, nil
	case "Terminated":
		return DriverStatusTerminated, nil
	default:
		return 0, fmt.Errorf("invalid driver status: %s", s)
	}
}
