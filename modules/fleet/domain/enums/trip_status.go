package enums

import "fmt"

type TripStatus int

const (
	TripStatusScheduled TripStatus = iota
	TripStatusInProgress
	TripStatusCompleted
	TripStatusCancelled
)

func (t TripStatus) String() string {
	switch t {
	case TripStatusScheduled:
		return "Scheduled"
	case TripStatusInProgress:
		return "InProgress"
	case TripStatusCompleted:
		return "Completed"
	case TripStatusCancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}

func (t TripStatus) IsValid() bool {
	return t >= TripStatusScheduled && t <= TripStatusCancelled
}

func ParseTripStatus(s string) (TripStatus, error) {
	switch s {
	case "Scheduled":
		return TripStatusScheduled, nil
	case "InProgress":
		return TripStatusInProgress, nil
	case "Completed":
		return TripStatusCompleted, nil
	case "Cancelled":
		return TripStatusCancelled, nil
	default:
		return 0, fmt.Errorf("invalid trip status: %s", s)
	}
}
