package enums

import "fmt"

type FuelType int

const (
	FuelTypeGasoline FuelType = iota
	FuelTypeDiesel
	FuelTypeElectric
	FuelTypeHybrid
	FuelTypeCNG
)

func (f FuelType) String() string {
	switch f {
	case FuelTypeGasoline:
		return "Gasoline"
	case FuelTypeDiesel:
		return "Diesel"
	case FuelTypeElectric:
		return "Electric"
	case FuelTypeHybrid:
		return "Hybrid"
	case FuelTypeCNG:
		return "CNG"
	default:
		return "Unknown"
	}
}

func (f FuelType) IsValid() bool {
	return f >= FuelTypeGasoline && f <= FuelTypeCNG
}

func ParseFuelType(s string) (FuelType, error) {
	switch s {
	case "Gasoline":
		return FuelTypeGasoline, nil
	case "Diesel":
		return FuelTypeDiesel, nil
	case "Electric":
		return FuelTypeElectric, nil
	case "Hybrid":
		return FuelTypeHybrid, nil
	case "CNG":
		return FuelTypeCNG, nil
	default:
		return 0, fmt.Errorf("invalid fuel type: %s", s)
	}
}
