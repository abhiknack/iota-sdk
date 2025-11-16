# Code Generator Examples

Real-world examples of using the IOTA SDK code generator.

## Example 1: Fleet Management - Vehicle Entity

### Command

```bash
./scripts/generate.sh crud \
  -m fleet \
  -e Vehicle \
  -f "Make:string:required,Model:string:required,Year:int:min=1900,VIN:string:len=17,LicensePlate:string:required,CurrentOdometer:int:min=0"
```

### Generated Files

```
modules/fleet/
├── domain/aggregates/vehicle/
│   ├── vehicle.go
│   ├── vehicle_repository.go
│   └── vehicle_events.go
├── infrastructure/persistence/
│   └── vehicle_repository.go
├── services/
│   └── vehicle_service.go
└── presentation/controllers/
    ├── dtos/vehicle_dto.go
    └── vehicle_controller.go
```

### Post-Generation Steps

1. **Create migration:**
```bash
./scripts/generate.sh migration
```

Edit `migrations/changes-{timestamp}.sql`:
```sql
-- +migrate Up
CREATE TABLE fleet_vehicles (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    make VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INT NOT NULL,
    vin VARCHAR(17) UNIQUE NOT NULL,
    license_plate VARCHAR(20) NOT NULL,
    current_odometer INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_fleet_vehicles_tenant_id ON fleet_vehicles(tenant_id);
CREATE INDEX idx_fleet_vehicles_deleted_at ON fleet_vehicles(deleted_at);
CREATE INDEX idx_fleet_vehicles_vin ON fleet_vehicles(vin);

-- +migrate Down
DROP TABLE IF EXISTS fleet_vehicles;
```

2. **Register in module.go:**
```go
func (m *Module) Register(app application.Application) error {
    vehicleRepo := persistence.NewVehicleRepository()
    
    app.RegisterServices(
        services.NewVehicleService(vehicleRepo, app.EventPublisher()),
    )
    
    app.RegisterControllers(
        controllers.NewVehicleController(app),
    )
    
    return nil
}
```

3. **Add permissions:**
```go
// modules/fleet/permissions/constants.go
const (
    VehicleRead   = "fleet.vehicle.read"
    VehicleCreate = "fleet.vehicle.create"
    VehicleUpdate = "fleet.vehicle.update"
    VehicleDelete = "fleet.vehicle.delete"
)
```

4. **Build:**
```bash
make db migrate up
templ generate && make css
go vet ./...
```

## Example 2: CRM - Contact Entity

### Command

```bash
./scripts/generate.sh crud \
  -m crm \
  -e Contact \
  -f "FirstName:string:required,LastName:string:required,Email:string:email,Phone:string:max=20,Company:string:max=100"
```

### Use Case

Customer relationship management system needs to track contacts with their basic information.

### Customization After Generation

Add custom method to domain aggregate:
```go
// modules/crm/domain/aggregates/contact/contact.go

func (c *contact) FullName() string {
    return c.firstName + " " + c.lastName
}
```

Add custom repository query:
```go
// modules/crm/infrastructure/persistence/contact_repository.go

func (r *ContactRepository) GetByEmail(ctx context.Context, email string) (contact.Contact, error) {
    tenantID, err := composables.UseTenantID(ctx)
    if err != nil {
        return nil, err
    }
    
    query := repo.Join(selectContactQuery, "WHERE email = $1 AND tenant_id = $2 AND deleted_at IS NULL")
    contacts, err := r.queryContacts(ctx, query, email, tenantID)
    if err != nil {
        return nil, err
    }
    if len(contacts) == 0 {
        return nil, ErrContactNotFound
    }
    return contacts[0], nil
}
```

## Example 3: Warehouse - Product Entity

### Command

```bash
./scripts/generate.sh crud \
  -m warehouse \
  -e Product \
  -f "Name:string:required,SKU:string:required,Price:float64:min=0,Quantity:int:min=0,Description:string:max=500,CategoryID:uuid.UUID"
```

### Use Case

Inventory management system needs to track products with pricing and stock levels.

### Customization After Generation

Add business logic to service:
```go
// modules/warehouse/services/product_service.go

func (s *ProductService) AdjustStock(ctx context.Context, productID uuid.UUID, adjustment int) error {
    if err := composables.CanUser(ctx, permissions.ProductUpdate); err != nil {
        return err
    }
    
    product, err := s.repo.GetByID(ctx, productID)
    if err != nil {
        return err
    }
    
    newQuantity := product.Quantity() + adjustment
    if newQuantity < 0 {
        return fmt.Errorf("insufficient stock: current=%d, adjustment=%d", product.Quantity(), adjustment)
    }
    
    // Update product with new quantity
    // ... implementation
    
    return nil
}
```

## Example 4: HRM - Employee Entity

### Command

```bash
./scripts/generate.sh crud \
  -m hrm \
  -e Employee \
  -f "FirstName:string:required,LastName:string:required,Email:string:email,HireDate:time.Time:required,Salary:float64:min=0,DepartmentID:uuid.UUID:required"
```

### Use Case

Human resources management system needs to track employees with their employment details.

### Migration with Additional Constraints

```sql
-- +migrate Up
CREATE TABLE hrm_employees (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    hire_date TIMESTAMPTZ NOT NULL,
    salary DECIMAL(10,2) NOT NULL CHECK (salary >= 0),
    department_id uuid NOT NULL REFERENCES hrm_departments(id),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_hrm_employees_tenant_id ON hrm_employees(tenant_id);
CREATE INDEX idx_hrm_employees_department_id ON hrm_employees(department_id);
CREATE INDEX idx_hrm_employees_email ON hrm_employees(email);
CREATE INDEX idx_hrm_employees_deleted_at ON hrm_employees(deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS hrm_employees;
```

## Example 5: Finance - Invoice Entity

### Command

```bash
./scripts/generate.sh crud \
  -m finance \
  -e Invoice \
  -f "Number:string:required,Amount:float64:min=0,DueDate:time.Time:required,Status:int:required,CustomerID:uuid.UUID:required"
```

### Use Case

Financial management system needs to track invoices with payment status.

### Add Enum for Status

Create enum file:
```go
// modules/finance/domain/enums/invoice_status.go
package enums

type InvoiceStatus int

const (
    InvoiceStatusDraft InvoiceStatus = iota
    InvoiceStatusSent
    InvoiceStatusPaid
    InvoiceStatusOverdue
    InvoiceStatusCancelled
)

func (s InvoiceStatus) String() string {
    switch s {
    case InvoiceStatusDraft:
        return "Draft"
    case InvoiceStatusSent:
        return "Sent"
    case InvoiceStatusPaid:
        return "Paid"
    case InvoiceStatusOverdue:
        return "Overdue"
    case InvoiceStatusCancelled:
        return "Cancelled"
    default:
        return "Unknown"
    }
}
```

Update domain aggregate to use enum:
```go
// modules/finance/domain/aggregates/invoice/invoice.go

import "github.com/iota-uz/iota-sdk/modules/finance/domain/enums"

type Invoice interface {
    // ... other methods
    Status() enums.InvoiceStatus
}
```

## Example 6: Projects - Task Entity

### Command

```bash
./scripts/generate.sh crud \
  -m projects \
  -e Task \
  -f "Title:string:required,Description:string:max=1000,Priority:int:min=1,max=5,DueDate:time.Time,ProjectID:uuid.UUID:required,AssigneeID:uuid.UUID"
```

### Use Case

Project management system needs to track tasks with priorities and assignments.

### Add Custom Controller Endpoint

```go
// modules/projects/presentation/controllers/task_controller.go

func (c *TaskController) Assign(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    vars := mux.Vars(r)
    taskID, err := uuid.Parse(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    var dto struct {
        AssigneeID uuid.UUID `form:"AssigneeID" validate:"required"`
    }
    
    if err := composables.UseForm(&dto, r); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Custom business logic for task assignment
    // ...
}
```

Register custom route:
```go
func (c *TaskController) Register(r *mux.Router) {
    // ... existing routes
    
    setRouter.HandleFunc("/{id}/assign", c.Assign).Methods(http.MethodPost)
}
```

## Tips for Effective Usage

### 1. Start Simple

Generate basic entity first, then add complexity:
```bash
# Start with basic fields
./scripts/generate.sh crud -m module -e Entity -f "Name:string:required"

# Add custom methods and logic after generation
```

### 2. Use Descriptive Field Names

```bash
# Good
-f "FirstName:string:required,LastName:string:required,EmailAddress:string:email"

# Avoid
-f "fname:string:required,lname:string:required,email:string:email"
```

### 3. Plan Your Validations

```bash
# Include appropriate validations from the start
-f "Age:int:required,min=0,max=150,Email:string:required,email,Phone:string:max=20"
```

### 4. Consider Relationships

```bash
# Include foreign key fields
-f "Name:string:required,ParentID:uuid.UUID,CategoryID:uuid.UUID:required"
```

### 5. Review and Customize

Always review generated code and customize:
- Add custom business logic
- Add custom queries
- Add custom endpoints
- Add validation rules
- Add domain events

## Common Patterns

### Pattern 1: Master-Detail Relationship

```bash
# Generate master entity
./scripts/generate.sh crud -m sales -e Order -f "Number:string:required,Date:time.Time:required,CustomerID:uuid.UUID:required"

# Generate detail entity
./scripts/generate.sh crud -m sales -e OrderItem -f "OrderID:uuid.UUID:required,ProductID:uuid.UUID:required,Quantity:int:min=1,Price:float64:min=0"
```

### Pattern 2: Hierarchical Data

```bash
# Generate entity with self-reference
./scripts/generate.sh crud -m organization -e Department -f "Name:string:required,ParentID:uuid.UUID,ManagerID:uuid.UUID"
```

### Pattern 3: Audit Trail

```bash
# Generate entity with audit fields
./scripts/generate.sh crud -m audit -e AuditLog -f "EntityType:string:required,EntityID:uuid.UUID:required,Action:string:required,Changes:string,UserID:uuid.UUID:required"
```

## Troubleshooting

### Issue: Import Errors After Generation

**Solution:**
```bash
go mod tidy
```

### Issue: Template Compilation Errors

**Solution:**
```bash
templ generate
```

### Issue: Database Errors

**Solution:**
Check migration was applied:
```bash
make db migrate status
make db migrate up
```

### Issue: Permission Denied Errors

**Solution:**
Ensure permissions are registered in seed data and assigned to roles.

## Next Steps

After generating code:

1. Review generated files
2. Customize business logic
3. Create templates
4. Add translations
5. Write tests
6. Document custom features
