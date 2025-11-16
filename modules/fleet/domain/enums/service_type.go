package enums

import "fmt"

type ServiceType int

const (
	ServiceTypeOilChange ServiceType = iota
	ServiceTypeTireRotation
	ServiceTypeBrakeService
	ServiceTypeInspection
	ServiceTypeRepair
	ServiceTypeOther
)

func (s ServiceType) String() string {
	switch s {
	case ServiceTypeOilChange:
		return "OilChange"
	case ServiceTypeTireRotation:
		return "TireRotation"
	case ServiceTypeBrakeService:
		return "BrakeService"
	case ServiceTypeInspection:
		return "Inspection"
	case ServiceTypeRepair:
		return "Repair"
	case ServiceTypeOther:
		return "OtherService"
	default:
		return "Unknown"
	}
}

func (s ServiceType) IsValid() bool {
	return s >= ServiceTypeOilChange && s <= ServiceTypeOther
}

func ParseServiceType(str string) (ServiceType, error) {
	switch str {
	case "OilChange":
		return ServiceTypeOilChange, nil
	case "TireRotation":
		return ServiceTypeTireRotation, nil
	case "BrakeService":
		return ServiceTypeBrakeService, nil
	case "Inspection":
		return ServiceTypeInspection, nil
	case "Repair":
		return ServiceTypeRepair, nil
	case "OtherService", "Other":
		return ServiceTypeOther, nil
	default:
		return 0, fmt.Errorf("invalid service type: %s", str)
	}
}
