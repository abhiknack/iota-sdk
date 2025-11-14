# Design Document

## Overview

The Supermarket/Retail Module is a comprehensive point-of-sale and inventory management system built following Domain-Driven Design (DDD) principles within the IOTA SDK architecture. The module enables retail businesses to manage sales transactions, inventory, cash reconciliation, and multi-company operations through a cohesive set of domain aggregates and services.

The design follows the established IOTA SDK patterns with clear separation between domain logic, infrastructure concerns, and presentation layers. The module integrates with the existing core modules (authentication, permissions, file uploads) and provides event-driven communication for cross-module coordination.

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ POS Terminal │  │ Inventory UI │  │ Reports UI   │      │
│  │  (HTMX)      │  │   (HTMX)     │  │   (HTMX)     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────────┐
│                      Service Layer                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ POS Service  │  │ Item Service │  │ Sales Service│      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────────┐
│                      Domain Layer                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ POS Invoice  │  │ Item Master  │  │ Stock Ledger │      │
│  │  Aggregate   │  │  Aggregate   │  │  Aggregate   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────────┐
│                  Infrastructure Layer                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ PostgreSQL   │  │ Event Bus    │  │ File Storage │      │
│  │ Repositories │  │              │  │              │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

### Module Structure


Following the IOTA SDK DDD structure:

```
modules/retail/
├── domain/
│   ├── aggregates/
│   │   ├── pos_profile/
│   │   ├── pos_invoice/
│   │   ├── pos_shift/
│   │   ├── item/
│   │   ├── sales_order/
│   │   ├── stock_ledger/
│   │   ├── price_list/
│   │   └── payment_method/
│   ├── entities/
│   │   ├── invoice_item/
│   │   ├── order_item/
│   │   └── stock_entry/
│   └── value_objects/
│       ├── money.go
│       ├── quantity.go
│       └── item_code.go
├── infrastructure/
│   └── persistence/
│       ├── models/models.go
│       ├── schema/retail-schema.sql
│       └── {entity}_repository.go (for each aggregate)
├── services/
│   ├── pos_service.go
│   ├── item_service.go
│   ├── inventory_service.go
│   ├── sales_order_service.go
│   └── pricing_service.go
├── presentation/
│   ├── controllers/
│   ├── templates/
│   ├── viewmodels/
│   └── locales/
├── module.go
├── links.go
└── permissions/constants.go
```

## Components and Interfaces

### Domain Aggregates

#### 1. POS Profile Aggregate

**Purpose**: Configuration for point-of-sale terminals

**Key Entities**:
- POSProfile (root)
- POSPaymentMethod (entity)

**Repository Interface**:
```go
type POSProfileRepository interface {
    Create(profile *POSProfile) error
    Update(profile *POSProfile) error
    GetByID(id int) (*POSProfile, error)
    GetByCompany(companyID int) ([]*POSProfile, error)
    Delete(id int) error
}
```

**Domain Events**:
- POSProfileCreated
- POSProfileUpdated
- POSProfileActivated
- POSProfileDeactivated

**Business Rules**:
- Must have at least one enabled payment method
- Must be associated with a company and warehouse
- Cannot be deleted if active shifts exist


#### 2. Item Aggregate

**Purpose**: Master data for products/items

**Key Entities**:
- Item (root)
- ItemVariant (entity)
- ItemPrice (entity)

**Repository Interface**:
```go
type ItemRepository interface {
    Create(item *Item) error
    Update(item *Item) error
    GetByID(id int) (*Item, error)
    GetByCode(code string) (*Item, error)
    List(filters ItemFilters) ([]*Item, error)
    Delete(id int) error
    GetVariants(itemID int) ([]*ItemVariant, error)
}
```

**Domain Events**:
- ItemCreated
- ItemUpdated
- ItemPriceChanged
- ItemDiscontinued

**Business Rules**:
- Item code must be unique within company
- Cannot delete items with existing stock
- Variants inherit base item properties
- Must have at least one unit of measure

#### 3. POS Invoice Aggregate

**Purpose**: Sales transaction at point of sale

**Key Entities**:
- POSInvoice (root)
- InvoiceItem (entity)
- InvoicePayment (entity)

**Repository Interface**:
```go
type POSInvoiceRepository interface {
    Create(invoice *POSInvoice) error
    Update(invoice *POSInvoice) error
    GetByID(id int) (*POSInvoice, error)
    GetByShift(shiftID int) ([]*POSInvoice, error)
    GetByDateRange(start, end time.Time, companyID int) ([]*POSInvoice, error)
    Cancel(id int) error
}
```

**Domain Events**:
- InvoiceCreated
- InvoiceSubmitted
- InvoiceCancelled
- InvoiceReturned
- PaymentReceived

**Business Rules**:
- Must have at least one item
- Total payments must equal grand total
- Cannot modify submitted invoices
- Returns must reference original invoice
- Stock must be available for all items


#### 4. POS Shift Aggregate

**Purpose**: Cash reconciliation for POS sessions

**Key Entities**:
- POSShift (root)
- OpeningEntry (entity)
- ClosingEntry (entity)
- ShiftPaymentBalance (entity)

**Repository Interface**:
```go
type POSShiftRepository interface {
    Create(shift *POSShift) error
    Update(shift *POSShift) error
    GetByID(id int) (*POSShift, error)
    GetActiveShift(profileID, userID int) (*POSShift, error)
    Close(id int, closingEntry *ClosingEntry) error
    GetShiftSummary(id int) (*ShiftSummary, error)
}
```

**Domain Events**:
- ShiftOpened
- ShiftClosed
- VarianceDetected

**Business Rules**:
- Only one active shift per user per POS profile
- Cannot close shift with pending transactions
- Closing balances must be entered for all payment methods
- Variance calculation: actual - (opening + transactions)

#### 5. Sales Order Aggregate

**Purpose**: Customer orders for future fulfillment

**Key Entities**:
- SalesOrder (root)
- OrderItem (entity)

**Repository Interface**:
```go
type SalesOrderRepository interface {
    Create(order *SalesOrder) error
    Update(order *SalesOrder) error
    GetByID(id int) (*SalesOrder, error)
    GetByCustomer(customerID int) ([]*SalesOrder, error)
    GetPendingOrders(companyID int) ([]*SalesOrder, error)
    UpdateStatus(id int, status OrderStatus) error
}
```

**Domain Events**:
- OrderCreated
- OrderConfirmed
- OrderFulfilled
- OrderCancelled
- OrderPartiallyFulfilled

**Business Rules**:
- Cannot fulfill more than ordered quantity
- Cancelled orders release reserved stock
- Confirmed orders reserve stock
- Can convert to POS invoice when ready


#### 6. Stock Ledger Aggregate

**Purpose**: Chronological record of inventory movements

**Key Entities**:
- StockLedgerEntry (root)

**Repository Interface**:
```go
type StockLedgerRepository interface {
    CreateEntry(entry *StockLedgerEntry) error
    GetByItem(itemID int, warehouseID int) ([]*StockLedgerEntry, error)
    GetCurrentBalance(itemID, warehouseID int) (float64, error)
    GetMovementHistory(filters StockFilters) ([]*StockLedgerEntry, error)
}
```

**Domain Events**:
- StockIncreased
- StockDecreased
- LowStockAlert

**Business Rules**:
- Entries are immutable once created
- Balance cannot go negative
- Each entry must reference a transaction
- Running balance maintained per item per warehouse

#### 7. Price List Aggregate

**Purpose**: Pricing configuration and rules

**Key Entities**:
- PriceList (root)
- PriceListItem (entity)
- PricingRule (entity)

**Repository Interface**:
```go
type PriceListRepository interface {
    Create(priceList *PriceList) error
    Update(priceList *PriceList) error
    GetByID(id int) (*PriceList, error)
    GetActiveForCompany(companyID int) ([]*PriceList, error)
    GetItemPrice(itemID int, priceListID int) (*PriceListItem, error)
}
```

**Domain Events**:
- PriceListCreated
- PriceListUpdated
- PriceChanged
- PricingRuleApplied

**Business Rules**:
- Price list must have validity dates
- Expired price lists revert to base pricing
- Multiple rules: highest priority wins
- Promotional rules can override base prices


#### 8. Payment Method Aggregate

**Purpose**: Configuration for payment types

**Key Entities**:
- PaymentMethod (root)

**Repository Interface**:
```go
type PaymentMethodRepository interface {
    Create(method *PaymentMethod) error
    Update(method *PaymentMethod) error
    GetByID(id int) (*PaymentMethod, error)
    GetEnabled(companyID int) ([]*PaymentMethod, error)
    Delete(id int) error
}
```

**Domain Events**:
- PaymentMethodCreated
- PaymentMethodEnabled
- PaymentMethodDisabled

**Business Rules**:
- Must have unique name per company
- Cannot delete if used in transactions
- Integration config required for electronic payments

### Service Layer

#### POSService

**Responsibilities**:
- Orchestrate POS invoice creation and submission
- Validate stock availability
- Apply pricing rules
- Process payments
- Generate receipts

**Key Methods**:
```go
type POSService interface {
    CreateInvoice(dto CreateInvoiceDTO) (*POSInvoice, error)
    AddItem(invoiceID, itemID int, quantity float64) error
    ApplyDiscount(invoiceID int, discount Discount) error
    ProcessPayment(invoiceID int, payments []Payment) error
    SubmitInvoice(invoiceID int) error
    CancelInvoice(invoiceID int, reason string) error
    CreateReturn(originalInvoiceID int, items []ReturnItem) error
}
```

**Dependencies**:
- POSInvoiceRepository
- ItemRepository
- StockLedgerRepository
- PricingService
- EventPublisher


#### ItemService

**Responsibilities**:
- Manage item master data
- Handle item variants
- Validate item codes
- Manage item groups

**Key Methods**:
```go
type ItemService interface {
    CreateItem(dto CreateItemDTO) (*Item, error)
    UpdateItem(id int, dto UpdateItemDTO) error
    GetItem(id int) (*Item, error)
    ListItems(filters ItemFilters) ([]*Item, error)
    CreateVariant(itemID int, dto VariantDTO) (*ItemVariant, error)
    DiscontinueItem(id int) error
}
```

**Dependencies**:
- ItemRepository
- EventPublisher

#### InventoryService

**Responsibilities**:
- Track stock movements
- Maintain stock balances
- Generate stock alerts
- Handle stock adjustments

**Key Methods**:
```go
type InventoryService interface {
    RecordStockEntry(entry StockEntryDTO) error
    GetStockBalance(itemID, warehouseID int) (float64, error)
    GetLowStockItems(warehouseID int) ([]*Item, error)
    AdjustStock(itemID, warehouseID int, quantity float64, reason string) error
    GetMovementHistory(filters StockFilters) ([]*StockLedgerEntry, error)
}
```

**Dependencies**:
- StockLedgerRepository
- ItemRepository
- EventPublisher

#### ShiftService

**Responsibilities**:
- Manage POS shift lifecycle
- Handle opening/closing entries
- Calculate variances
- Generate shift reports

**Key Methods**:
```go
type ShiftService interface {
    OpenShift(dto OpenShiftDTO) (*POSShift, error)
    CloseShift(shiftID int, dto CloseShiftDTO) error
    GetActiveShift(profileID, userID int) (*POSShift, error)
    GetShiftSummary(shiftID int) (*ShiftSummary, error)
    GetShiftTransactions(shiftID int) ([]*POSInvoice, error)
}
```

**Dependencies**:
- POSShiftRepository
- POSInvoiceRepository
- EventPublisher


#### SalesOrderService

**Responsibilities**:
- Manage sales orders
- Handle order fulfillment
- Convert orders to invoices
- Track order status

**Key Methods**:
```go
type SalesOrderService interface {
    CreateOrder(dto CreateOrderDTO) (*SalesOrder, error)
    UpdateOrder(id int, dto UpdateOrderDTO) error
    ConfirmOrder(id int) error
    ConvertToInvoice(orderID int) (*POSInvoice, error)
    CancelOrder(id int, reason string) error
    GetPendingOrders(companyID int) ([]*SalesOrder, error)
}
```

**Dependencies**:
- SalesOrderRepository
- POSService
- InventoryService
- EventPublisher

#### PricingService

**Responsibilities**:
- Calculate item prices
- Apply pricing rules
- Handle promotions
- Manage price lists

**Key Methods**:
```go
type PricingService interface {
    GetItemPrice(itemID int, priceListID int, quantity float64) (Money, error)
    ApplyPricingRules(items []InvoiceItem, context PricingContext) error
    CreatePriceList(dto PriceListDTO) (*PriceList, error)
    UpdatePriceList(id int, dto PriceListDTO) error
    GetActivePriceList(companyID int) (*PriceList, error)
}
```

**Dependencies**:
- PriceListRepository
- ItemRepository
- EventPublisher

## Data Models

### Database Schema Overview

#### Core Tables

**items**
- id (PK)
- company_id (FK)
- item_code (unique per company)
- item_name
- description
- item_group_id (FK)
- unit_of_measure
- has_variants (boolean)
- is_stock_item (boolean)
- reorder_level
- standard_rate
- created_at, updated_at

**item_variants**
- id (PK)
- item_id (FK)
- variant_code
- attribute_values (JSONB)
- created_at


**pos_profiles**
- id (PK)
- company_id (FK)
- profile_name
- warehouse_id (FK)
- price_list_id (FK)
- is_active (boolean)
- allow_negative_stock (boolean)
- created_at, updated_at

**pos_profile_payment_methods**
- id (PK)
- pos_profile_id (FK)
- payment_method_id (FK)
- is_default (boolean)

**pos_shifts**
- id (PK)
- pos_profile_id (FK)
- user_id (FK)
- company_id (FK)
- shift_start (timestamp)
- shift_end (timestamp)
- status (enum: open, closed)
- created_at, updated_at

**pos_shift_balances**
- id (PK)
- pos_shift_id (FK)
- payment_method_id (FK)
- opening_amount (decimal)
- closing_amount (decimal)
- expected_amount (decimal)
- variance (decimal)
- entry_type (enum: opening, closing)

**pos_invoices**
- id (PK)
- invoice_number (unique)
- pos_shift_id (FK)
- company_id (FK)
- customer_id (FK, nullable)
- invoice_date (timestamp)
- status (enum: draft, submitted, cancelled, returned)
- subtotal (decimal)
- tax_amount (decimal)
- discount_amount (decimal)
- grand_total (decimal)
- return_against (FK to pos_invoices, nullable)
- created_at, updated_at

**pos_invoice_items**
- id (PK)
- pos_invoice_id (FK)
- item_id (FK)
- item_code
- item_name
- quantity (decimal)
- rate (decimal)
- amount (decimal)
- warehouse_id (FK)

**pos_invoice_payments**
- id (PK)
- pos_invoice_id (FK)
- payment_method_id (FK)
- amount (decimal)
- reference_number (nullable)


**sales_orders**
- id (PK)
- order_number (unique)
- company_id (FK)
- customer_id (FK)
- order_date (timestamp)
- delivery_date (date)
- status (enum: draft, confirmed, partially_fulfilled, fulfilled, cancelled)
- total_amount (decimal)
- created_at, updated_at

**sales_order_items**
- id (PK)
- sales_order_id (FK)
- item_id (FK)
- quantity (decimal)
- rate (decimal)
- amount (decimal)
- delivered_quantity (decimal)

**stock_ledger_entries**
- id (PK)
- item_id (FK)
- warehouse_id (FK)
- company_id (FK)
- posting_date (timestamp)
- voucher_type (varchar)
- voucher_no (varchar)
- actual_qty (decimal)
- qty_after_transaction (decimal)
- stock_value (decimal)
- created_at

**price_lists**
- id (PK)
- company_id (FK)
- price_list_name
- currency
- valid_from (date)
- valid_to (date, nullable)
- is_active (boolean)
- created_at, updated_at

**price_list_items**
- id (PK)
- price_list_id (FK)
- item_id (FK)
- rate (decimal)

**pricing_rules**
- id (PK)
- company_id (FK)
- rule_name
- priority (integer)
- item_group_id (FK, nullable)
- customer_group_id (FK, nullable)
- min_quantity (decimal, nullable)
- discount_percentage (decimal, nullable)
- valid_from (date)
- valid_to (date, nullable)
- is_active (boolean)

**payment_methods**
- id (PK)
- company_id (FK)
- method_name
- method_type (enum: cash, card, mobile, bank_transfer)
- account_id (FK, nullable)
- integration_config (JSONB, nullable)
- is_active (boolean)
- created_at, updated_at


### Multi-Company Data Isolation

All tables include `company_id` foreign key to ensure data isolation:

```sql
-- Row-level security policy example
CREATE POLICY company_isolation ON pos_invoices
    USING (company_id = current_setting('app.current_company_id')::int);
```

Alternatively, application-level filtering in repositories:

```go
func (r *POSInvoiceRepository) List(companyID int) ([]*POSInvoice, error) {
    query := `SELECT * FROM pos_invoices WHERE company_id = $1`
    // ...
}
```

## Error Handling

### Domain Errors

Custom error types for business rule violations:

```go
type DomainError struct {
    Code    string
    Message string
    Field   string
}

// Examples
var (
    ErrInsufficientStock = &DomainError{
        Code:    "INSUFFICIENT_STOCK",
        Message: "Item quantity exceeds available stock",
    }
    ErrShiftAlreadyOpen = &DomainError{
        Code:    "SHIFT_ALREADY_OPEN",
        Message: "User already has an active shift",
    }
    ErrInvalidPaymentTotal = &DomainError{
        Code:    "INVALID_PAYMENT_TOTAL",
        Message: "Payment total does not match invoice total",
    }
)
```

### Error Propagation

Services return domain errors, controllers map to HTTP responses:

```go
// Service layer
func (s *POSService) SubmitInvoice(id int) error {
    invoice, err := s.repo.GetByID(id)
    if err != nil {
        return err
    }
    
    if !s.inventoryService.CheckStock(invoice.Items) {
        return ErrInsufficientStock
    }
    
    // Business logic...
    return nil
}

// Controller layer
func (c *POSController) SubmitInvoice(w http.ResponseWriter, r *http.Request) {
    err := c.posService.SubmitInvoice(invoiceID)
    if errors.Is(err, ErrInsufficientStock) {
        http.Error(w, "Insufficient stock", http.StatusBadRequest)
        return
    }
    // ...
}
```


## Testing Strategy

### Unit Tests

**Domain Layer**:
- Test business rules in aggregates
- Validate state transitions
- Test invariant enforcement

```go
func TestPOSInvoice_AddItem_ValidatesStock(t *testing.T) {
    invoice := NewPOSInvoice(...)
    item := &InvoiceItem{Quantity: 100}
    
    err := invoice.AddItem(item, availableStock: 50)
    
    assert.Error(t, err)
    assert.Equal(t, ErrInsufficientStock, err)
}
```

**Service Layer**:
- Test orchestration logic
- Mock repository dependencies
- Verify event publishing

```go
func TestPOSService_SubmitInvoice_PublishesEvent(t *testing.T) {
    mockRepo := &MockPOSInvoiceRepository{}
    mockPublisher := &MockEventPublisher{}
    service := NewPOSService(mockRepo, mockPublisher)
    
    err := service.SubmitInvoice(1)
    
    assert.NoError(t, err)
    assert.True(t, mockPublisher.WasPublished("InvoiceSubmitted"))
}
```

### Integration Tests

**Repository Layer**:
- Test database operations
- Verify transaction handling
- Test complex queries

```go
func TestPOSInvoiceRepository_GetByShift_ReturnsAllInvoices(t *testing.T) {
    db := setupTestDB(t)
    repo := NewPOSInvoiceRepository(db)
    
    invoices, err := repo.GetByShift(shiftID)
    
    assert.NoError(t, err)
    assert.Len(t, invoices, 3)
}
```

### Controller Tests

**HTTP Handlers**:
- Test request/response handling
- Verify HTMX interactions
- Test authorization

```go
func TestPOSController_CreateInvoice_RequiresAuth(t *testing.T) {
    controller := setupController(t)
    req := httptest.NewRequest("POST", "/retail/pos/invoices", nil)
    
    w := httptest.NewRecorder()
    controller.CreateInvoice(w, req)
    
    assert.Equal(t, http.StatusUnauthorized, w.Code)
}
```


## User Interface Design

### POS Terminal Interface

**Layout**:
```
┌─────────────────────────────────────────────────────────┐
│  POS Terminal - Profile: Main Counter                   │
├─────────────────────────────────────────────────────────┤
│  Item Search: [_______________] [Scan Barcode]          │
├──────────────────────┬──────────────────────────────────┤
│  Cart Items          │  Item Catalog                    │
│  ┌────────────────┐  │  ┌────────────────────────────┐ │
│  │ Item A  x2     │  │  │ [Item 1] [Item 2] [Item 3] │ │
│  │ $10.00         │  │  │ [Item 4] [Item 5] [Item 6] │ │
│  │ [Remove]       │  │  └────────────────────────────┘ │
│  └────────────────┘  │                                  │
│                      │  Categories: [All] [Food] [...]  │
├──────────────────────┴──────────────────────────────────┤
│  Subtotal: $20.00    Tax: $2.00    Total: $22.00       │
│  [Cash] [Card] [Mobile]  [Complete Sale]               │
└─────────────────────────────────────────────────────────┘
```

**HTMX Interactions**:
- Add item: `hx-post="/retail/pos/cart/add"` → Updates cart section
- Remove item: `hx-delete="/retail/pos/cart/remove/{id}"` → Updates cart
- Complete sale: `hx-post="/retail/pos/invoice/submit"` → Shows receipt modal

### Shift Management Interface

**Opening Entry**:
```
┌─────────────────────────────────────────┐
│  Open Shift - Main Counter              │
├─────────────────────────────────────────┤
│  Date: 2025-11-11  Time: 09:00          │
│                                         │
│  Opening Balances:                      │
│  Cash:          [_______]               │
│  Credit Card:   [_______]               │
│  Mobile Pay:    [_______]               │
│                                         │
│  [Cancel]  [Open Shift]                 │
└─────────────────────────────────────────┘
```

**Closing Entry**:
```
┌─────────────────────────────────────────┐
│  Close Shift - Main Counter             │
├─────────────────────────────────────────┤
│  Shift Summary:                         │
│  Opened: 09:00  Duration: 8h 30m        │
│  Transactions: 45                       │
│                                         │
│  Payment Method  Expected   Actual  Var │
│  Cash           $1,250.00  [_____]  --- │
│  Credit Card    $2,340.00  [_____]  --- │
│  Mobile Pay     $  890.00  [_____]  --- │
│                                         │
│  [Cancel]  [Close Shift]                │
└─────────────────────────────────────────┘
```

### Inventory Management Interface

**Item List**:
- Filterable by item group, stock status
- Sortable by name, code, stock level
- Bulk actions: adjust stock, update prices
- Quick actions: edit, view history, reorder

**Item Form**:
- Basic info: code, name, description
- Pricing: standard rate, price lists
- Inventory: UOM, reorder level, warehouse
- Variants: size, color, etc. (if applicable)


### Sales Order Interface

**Order List**:
- Filter by status, customer, date range
- Quick convert to invoice action
- Status badges: Draft, Confirmed, Fulfilled

**Order Form**:
- Customer selection (with quick add)
- Delivery date picker
- Item selection with quantity
- Real-time total calculation
- Save as draft or confirm

### Reports Interface

**Dashboard**:
- Today's sales summary
- Active shifts count
- Low stock alerts
- Top selling items

**Sales Reports**:
- Daily/weekly/monthly sales
- Payment method breakdown
- Cashier performance
- Item-wise sales analysis

**Inventory Reports**:
- Current stock levels
- Stock movement history
- Reorder recommendations
- Stock valuation

## Integration Points

### Core Module Integration

**Authentication & Authorization**:
- Use existing user authentication
- Leverage RBAC for permissions
- Company-based access control

**File Uploads**:
- Item images via core upload repository
- Receipt attachments
- Invoice documents

**Event System**:
- Publish domain events to event bus
- Subscribe to relevant events from other modules
- Enable cross-module workflows

### External Integrations

**Payment Gateways** (Future):
- Card terminal integration
- Mobile payment APIs (e.g., Stripe, Square)
- Bank transfer verification

**Barcode Scanners**:
- USB/Bluetooth scanner support
- Barcode generation for items
- Quick item lookup

**Receipt Printers**:
- Thermal printer support
- Receipt template customization
- Print queue management

**Accounting Integration** (Future):
- Sync with finance module
- Journal entry generation
- Account reconciliation


## Security Considerations

### Data Access Control

**Company Isolation**:
- All queries filtered by company_id
- User permissions scoped to assigned companies
- Cross-company data access prevented

**Role-Based Permissions**:
```go
const (
    PermissionPOSCreate     = "retail.pos.create"
    PermissionPOSView       = "retail.pos.view"
    PermissionPOSCancel     = "retail.pos.cancel"
    PermissionShiftOpen     = "retail.shift.open"
    PermissionShiftClose    = "retail.shift.close"
    PermissionItemCreate    = "retail.item.create"
    PermissionItemUpdate    = "retail.item.update"
    PermissionStockAdjust   = "retail.stock.adjust"
    PermissionPriceUpdate   = "retail.price.update"
    PermissionReportsView   = "retail.reports.view"
)
```

### Audit Trail

**Transaction Logging**:
- All invoice submissions logged
- Stock adjustments tracked with user and reason
- Shift opening/closing recorded
- Price changes audited

**Change History**:
- Maintain modification history for critical entities
- Track who, when, what changed
- Immutable audit log

### Data Validation

**Input Sanitization**:
- Validate all user inputs
- Prevent SQL injection via parameterized queries
- Sanitize HTML in descriptions

**Business Rule Enforcement**:
- Server-side validation for all operations
- Cannot bypass stock checks
- Payment totals must match
- Shift constraints enforced

## Performance Considerations

### Database Optimization

**Indexing Strategy**:
```sql
-- Frequently queried columns
CREATE INDEX idx_pos_invoices_shift ON pos_invoices(pos_shift_id);
CREATE INDEX idx_pos_invoices_date ON pos_invoices(invoice_date);
CREATE INDEX idx_pos_invoices_company ON pos_invoices(company_id);
CREATE INDEX idx_stock_ledger_item_warehouse ON stock_ledger_entries(item_id, warehouse_id);
CREATE INDEX idx_items_code ON items(item_code, company_id);
```

**Query Optimization**:
- Use appropriate JOINs for related data
- Implement pagination for large lists
- Cache frequently accessed data (price lists, item catalog)


### Caching Strategy

**Application-Level Cache**:
- Item catalog (with TTL)
- Active price lists
- POS profile configurations
- Payment method settings

**Cache Invalidation**:
- On item updates
- On price list changes
- On configuration modifications

### Concurrency Handling

**Optimistic Locking**:
```go
type POSInvoice struct {
    ID      int
    Version int  // Incremented on each update
    // ... other fields
}

func (r *Repository) Update(invoice *POSInvoice) error {
    result := db.Exec(`
        UPDATE pos_invoices 
        SET ..., version = version + 1
        WHERE id = $1 AND version = $2
    `, invoice.ID, invoice.Version)
    
    if result.RowsAffected == 0 {
        return ErrConcurrentModification
    }
    return nil
}
```

**Stock Reservation**:
- Lock stock during invoice processing
- Release on cancellation
- Prevent overselling

## Deployment Considerations

### Database Migration

**Schema Versioning**:
- Use migration files in `infrastructure/persistence/schema/`
- Sequential numbering: `001_initial_schema.sql`, `002_add_variants.sql`
- Rollback scripts for each migration

**Data Seeding**:
- Default payment methods (Cash, Card)
- Sample item groups
- Base price list

### Configuration

**Environment Variables**:
```
RETAIL_DEFAULT_WAREHOUSE_ID=1
RETAIL_ENABLE_BARCODE_SCANNER=true
RETAIL_RECEIPT_PRINTER_URL=http://localhost:9100
RETAIL_LOW_STOCK_THRESHOLD=10
```

**Module Registration**:
```go
// In main application
app.RegisterModule(retail.NewModule())
```

## Future Enhancements

### Phase 2 Features

1. **Customer Loyalty Program**
   - Points accumulation
   - Reward redemption
   - Tier-based benefits

2. **Advanced Inventory**
   - Batch/serial number tracking
   - Expiry date management
   - Multi-location transfers

3. **E-commerce Integration**
   - Online order sync
   - Unified inventory
   - Click-and-collect

4. **Analytics & BI**
   - Predictive analytics
   - Demand forecasting
   - Customer behavior analysis

5. **Mobile POS**
   - Tablet/mobile app
   - Offline mode
   - Sync when online

6. **Supplier Management**
   - Purchase orders
   - Goods receipt
   - Supplier payments

## Conclusion

This design provides a comprehensive foundation for a retail/supermarket module that follows IOTA SDK's DDD architecture. The modular structure allows for incremental implementation, starting with core POS functionality and expanding to advanced features. The design emphasizes data integrity, multi-company support, and seamless integration with existing modules while maintaining clear separation of concerns across all layers.
