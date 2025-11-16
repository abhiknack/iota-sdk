# Fleet Management Services

## Notification System

The fleet module includes a notification system that monitors and alerts on important fleet events.

### NotificationService

The `NotificationService` handles the creation and sending of notifications for various fleet events:

- **License Expiry**: Notifies when driver licenses are expiring within a specified number of days
- **Registration Expiry**: Notifies when vehicle registrations are expiring
- **Insurance Expiry**: Notifies when vehicle insurance is expiring
- **Maintenance Due**: Notifies when scheduled maintenance is due
- **Fuel Anomaly**: Notifies when unusual fuel efficiency is detected

### SchedulerService

The `SchedulerService` runs background jobs on a schedule to check for notification conditions.

#### Usage

To start the scheduler, retrieve the service from the application and call `Start()`:

```go
// Get the scheduler service
schedulerSvc := app.Service(&services.SchedulerService{}).(*services.SchedulerService)

// Set tenant IDs to monitor (can be retrieved from database)
tenantIDs := []uuid.UUID{
    uuid.MustParse("tenant-id-1"),
    uuid.MustParse("tenant-id-2"),
}
schedulerSvc.SetTenantIDs(tenantIDs)

// Start the scheduler (runs in background)
schedulerSvc.Start()

// To stop the scheduler when shutting down
defer schedulerSvc.Stop()
```

#### Default Jobs

The scheduler automatically registers the following daily jobs:

1. **check_expiring_licenses**: Checks for driver licenses expiring within 30 days
2. **check_expiring_registrations**: Checks for vehicle registrations expiring within 30 days
3. **check_expiring_insurance**: Checks for vehicle insurance expiring within 30 days
4. **check_due_maintenance**: Checks for maintenance that is due

#### Manual Execution

You can also run all jobs once without starting the scheduler:

```go
schedulerSvc.RunOnce()
```

### Notification Events

All notifications are published to the event bus and can be subscribed to by other modules:

```go
app.EventPublisher().Subscribe(func(notification *services.Notification) {
    // Handle notification
    fmt.Printf("Notification: %s - %s\n", notification.Title, notification.Message)
})
```

### Logging

All notifications are logged with structured fields including:
- notification_id
- tenant_id
- notification_type
- title
- message
- data (additional context)

Check application logs for notification history.
