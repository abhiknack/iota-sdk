# Requirements Document

## Introduction

This document defines the requirements for a comprehensive Supermarket/Retail Management Module for the IOTA SDK application. The module will provide Point of Sale (POS) functionality, inventory management, sales order processing, and multi-company support, similar to Frappe's retail capabilities. The system will enable retail businesses to manage their daily operations including sales transactions, inventory tracking, cash management, and reporting across multiple company entities.

## Glossary

- **Retail Module**: The supermarket/retail management system component
- **POS System**: Point of Sale transaction processing system
- **POS Profile**: Configuration settings for a specific point of sale terminal
- **POS Invoice**: Sales transaction document generated at point of sale
- **POS Opening Entry**: Record of cash and payment methods at shift start
- **POS Closing Entry**: Record of cash and payment methods at shift end
- **Sales Order**: Customer order document that may be fulfilled immediately or later
- **Item Master**: Central repository of product/item information
- **Inventory System**: Stock tracking and management system
- **Company Entity**: Separate business entity in multi-company setup
- **Stock Ledger**: Chronological record of all inventory movements
- **Payment Method**: Mode of payment (cash, card, mobile payment, etc.)
- **Shift**: A working period for a POS terminal with opening and closing entries

## Requirements

### Requirement 1: POS Profile Management

**User Story:** As a retail manager, I want to configure POS profiles for different terminals, so that each point of sale operates with appropriate settings and permissions.

#### Acceptance Criteria

1. WHEN a retail manager creates a POS profile, THE Retail Module SHALL store the profile name, company association, warehouse assignment, price list, and payment method configurations
2. WHEN a POS profile is assigned to a terminal, THE Retail Module SHALL enforce the configured item groups, customer groups, and pricing rules for that terminal
3. WHERE multi-company mode is enabled, THE Retail Module SHALL restrict POS profile access to users authorized for the associated company
4. THE Retail Module SHALL validate that each POS profile has at least one enabled payment method before activation
5. WHEN a user updates a POS profile, THE Retail Module SHALL apply changes to new transactions while preserving historical transaction settings

### Requirement 2: Item and Inventory Management

**User Story:** As a store manager, I want to manage item master data and track inventory levels, so that I can maintain accurate stock information and prevent stockouts.

#### Acceptance Criteria

1. WHEN a user creates an item, THE Retail Module SHALL store item code, name, description, unit of measure, item group, and pricing information
2. THE Retail Module SHALL maintain real-time inventory balances for each item across all warehouses
3. WHEN an item quantity falls below the reorder level, THE Retail Module SHALL generate a low stock alert
4. THE Retail Module SHALL record all inventory movements in the Stock Ledger with timestamp, quantity, warehouse, and transaction reference
5. WHERE an item has variants (size, color, etc.), THE Retail Module SHALL manage each variant as a separate stock keeping unit with shared master attributes

### Requirement 3: POS Opening Entry

**User Story:** As a cashier, I want to record opening cash and payment method balances at shift start, so that I can reconcile transactions at shift end.

#### Acceptance Criteria

1. WHEN a cashier starts a shift, THE Retail Module SHALL create a POS Opening Entry with the POS profile, user, date, and time
2. THE Retail Module SHALL require the cashier to enter opening balances for each enabled payment method in the POS profile
3. THE Retail Module SHALL prevent creating multiple active opening entries for the same POS profile and user combination
4. WHEN an opening entry is submitted, THE Retail Module SHALL set the shift status to "Open" and enable transaction processing
5. THE Retail Module SHALL associate all POS invoices created during the shift with the corresponding opening entry

### Requirement 4: POS Invoice Processing

**User Story:** As a cashier, I want to process customer sales transactions quickly, so that I can serve customers efficiently and maintain accurate sales records.

#### Acceptance Criteria

1. WHEN a cashier adds items to a POS invoice, THE Retail Module SHALL calculate line totals, apply taxes, and compute the grand total in real-time
2. THE Retail Module SHALL validate item availability against current stock levels before allowing invoice submission
3. WHEN a POS invoice is submitted, THE Retail Module SHALL reduce inventory quantities, record the sale, and update the Stock Ledger
4. THE Retail Module SHALL support multiple payment methods on a single invoice with split payment functionality
5. WHEN a customer requests a return, THE Retail Module SHALL create a return invoice that reverses the original transaction and restores inventory

### Requirement 5: POS Closing Entry

**User Story:** As a cashier, I want to reconcile my shift transactions and close my POS session, so that I can account for all cash and payments received.

#### Acceptance Criteria

1. WHEN a cashier ends a shift, THE Retail Module SHALL create a POS Closing Entry linked to the corresponding opening entry
2. THE Retail Module SHALL calculate expected closing balances by summing opening balances and transaction amounts for each payment method
3. THE Retail Module SHALL require the cashier to enter actual closing balances for reconciliation
4. WHEN actual and expected balances differ, THE Retail Module SHALL calculate and display the variance for each payment method
5. WHEN a closing entry is submitted, THE Retail Module SHALL set the shift status to "Closed" and prevent further transactions on that shift

### Requirement 6: Sales Order Management

**User Story:** As a sales associate, I want to create sales orders for customers, so that I can process advance orders and scheduled deliveries.

#### Acceptance Criteria

1. WHEN a user creates a sales order, THE Retail Module SHALL store customer information, delivery date, items, quantities, and pricing
2. THE Retail Module SHALL track the fulfillment status of each sales order (Draft, Confirmed, Partially Delivered, Delivered, Cancelled)
3. WHEN a sales order is converted to a POS invoice, THE Retail Module SHALL transfer all line items and update the order status
4. THE Retail Module SHALL allow partial fulfillment of sales orders with multiple invoices
5. WHEN a sales order is cancelled, THE Retail Module SHALL release any reserved inventory and update the order status

### Requirement 7: Multi-Company Support

**User Story:** As a business owner with multiple retail entities, I want to manage separate companies within one system, so that I can maintain isolated financial and operational data.

#### Acceptance Criteria

1. WHEN multi-company mode is enabled, THE Retail Module SHALL require company selection for all transactions and master data
2. THE Retail Module SHALL isolate inventory, sales, and financial data by company entity
3. THE Retail Module SHALL restrict user access to companies based on assigned permissions
4. WHEN generating reports, THE Retail Module SHALL filter data by the selected company or provide consolidated views where authorized
5. THE Retail Module SHALL maintain separate number series for documents (invoices, orders) per company

### Requirement 8: Reporting and Analytics

**User Story:** As a retail manager, I want to view sales reports and analytics, so that I can make informed business decisions.

#### Acceptance Criteria

1. THE Retail Module SHALL generate daily sales summary reports showing total sales, payment method breakdown, and top-selling items
2. THE Retail Module SHALL provide inventory reports displaying current stock levels, stock value, and movement history
3. WHEN a manager requests a shift report, THE Retail Module SHALL display all transactions between opening and closing entries with reconciliation details
4. THE Retail Module SHALL calculate and display key performance indicators including average transaction value, items per transaction, and sales per hour
5. WHERE multi-company mode is enabled, THE Retail Module SHALL provide company-wise comparative reports

### Requirement 9: Payment Method Configuration

**User Story:** As a retail manager, I want to configure available payment methods, so that cashiers can accept various forms of payment.

#### Acceptance Criteria

1. WHEN a manager creates a payment method, THE Retail Module SHALL store the method name, type (cash, card, mobile, etc.), and account mapping
2. THE Retail Module SHALL allow enabling or disabling payment methods per POS profile
3. WHEN a payment method requires external integration (card terminal, mobile payment), THE Retail Module SHALL provide configuration fields for integration parameters
4. THE Retail Module SHALL validate that at least one payment method is enabled before allowing POS transactions
5. THE Retail Module SHALL track payment method usage statistics for reporting purposes

### Requirement 10: Price List and Pricing Rules

**User Story:** As a pricing manager, I want to configure price lists and promotional rules, so that correct prices are applied at point of sale.

#### Acceptance Criteria

1. WHEN a user creates a price list, THE Retail Module SHALL store item-specific prices with validity dates and currency
2. THE Retail Module SHALL apply the POS profile's default price list to all transactions unless overridden
3. WHEN multiple pricing rules apply to an item, THE Retail Module SHALL apply the rule with the highest priority
4. THE Retail Module SHALL support promotional pricing rules based on quantity, customer group, or time period
5. WHEN a price list or rule expires, THE Retail Module SHALL automatically revert to the base price list
