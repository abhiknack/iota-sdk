package enums

import "fmt"

type VehicleStatus int

const (
	VehicleStatusAvailable VehicleStatus = iota
	VehicleStatusInUse
	VehicleStatusMaintenance
	VehicleStatusOutOfService
	VehicleStatusRetired
)

func (s VehicleStatus) String() string {
	switch s {
	case VehicleStatusAvailable:
		return "Available"
	case VehicleStatusInUse:
		return "InUse"
	case VehicleStatusMaintenance:
		return "Maintenance"
	case VehicleStatusOutOfService:
		return "OutOfService"
	case VehicleStatusRetired:
		return "Retired"
	default:
		return "Unknown"
	}
}

func (s VehicleStatus) IsValid() bool {
	return s >= VehicleStatusAvailable && s <= VehicleStatusRetired
}

func ParseVehicleStatus(s string) (VehicleStatus, error) {
	switch s {
	case "Available":
		return VehicleStatusAvailable, nil
	case "InUse":
		return VehicleStatusInUse, nil
	case "Maintenance":
		return VehicleStatusMaintenance, nil
	case "OutOfService":
		return VehicleStatusOutOfService, nil
	case "Retired":
		return VehicleStatusRetired, nil
	default:
		return 0, fmt.Errorf("invalid vehicle status: %s", s)
	}
}

func (s VehicleStatus) CanTransitionTo(target VehicleStatus) bool {
	if !s.IsValid() || !target.IsValid() {
		return false
	}

	switch s {
	case VehicleStatusAvailable:
		return target == VehicleStatusInUse || target == VehicleStatusMaintenance ||
			target == VehicleStatusOutOfService || target == VehicleStatusRetired
	case VehicleStatusInUse:
		return target == VehicleStatusAvailable || target == VehicleStatusMaintenance ||
			target == VehicleStatusOutOfService
	case VehicleStatusMaintenance:
		return target == VehicleStatusAvailable || target == VehicleStatusOutOfService ||
			target == VehicleStatusRetired
	case VehicleStatusOutOfService:
		return target == VehicleStatusAvailable || target == VehicleStatusMaintenance ||
			target == VehicleStatusRetired
	case VehicleStatusRetired:
		return false
	default:
		return false
	}
}
