package viewmodels

type MaintenanceListViewModel struct {
	ID                  string
	VehicleID           string
	VehicleName         string
	ServiceType         string
	ServiceDate         string
	Odometer            int
	Cost                string
	ServiceProvider     string
	NextServiceDue      string
	NextServiceOdometer int
	IsDue               bool
}

type MaintenanceDetailViewModel struct {
	ID                  string
	TenantID            string
	VehicleID           string
	VehicleName         string
	ServiceType         string
	ServiceDate         string
	Odometer            int
	Cost                string
	ServiceProvider     string
	Description         string
	NextServiceDue      string
	NextServiceOdometer int
	CreatedAt           string
	UpdatedAt           string
}

type MaintenanceFormViewModel struct {
	ID                  string
	VehicleID           string
	ServiceType         string
	ServiceDate         string
	Odometer            int
	Cost                string
	ServiceProvider     string
	Description         string
	NextServiceDue      string
	NextServiceOdometer string
	IsEdit              bool
	Errors              map[string]string
}
