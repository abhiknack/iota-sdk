package viewmodels

type FuelEntryListViewModel struct {
	ID          string
	VehicleID   string
	VehicleName string
	DriverName  string
	Date        string
	Quantity    string
	Cost        string
	Odometer    int
	FuelType    string
	Location    string
	Efficiency  string
}

type FuelEntryDetailViewModel struct {
	ID          string
	TenantID    string
	VehicleID   string
	VehicleName string
	DriverID    string
	DriverName  string
	Date        string
	Quantity    string
	Cost        string
	Odometer    int
	FuelType    string
	Location    string
	Efficiency  string
	CreatedAt   string
	UpdatedAt   string
}

type FuelEntryFormViewModel struct {
	ID        string
	VehicleID string
	DriverID  string
	Date      string
	Quantity  string
	Cost      string
	Odometer  int
	FuelType  string
	Location  string
	IsEdit    bool
	Errors    map[string]string
}
