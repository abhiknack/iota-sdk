package services

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type JobFunc func(ctx context.Context, tenantID uuid.UUID) error

type Job struct {
	Name     string
	Interval time.Duration
	Fn       JobFunc
}

type SchedulerService struct {
	jobs            []Job
	notificationSvc *NotificationService
	logger          *logrus.Logger
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	tenantIDs       []uuid.UUID
}

func NewSchedulerService(
	notificationSvc *NotificationService,
	logger *logrus.Logger,
) *SchedulerService {
	ctx, cancel := context.WithCancel(context.Background())
	return &SchedulerService{
		jobs:            []Job{},
		notificationSvc: notificationSvc,
		logger:          logger,
		ctx:             ctx,
		cancel:          cancel,
		tenantIDs:       []uuid.UUID{},
	}
}

func (s *SchedulerService) RegisterJob(job Job) {
	s.jobs = append(s.jobs, job)
}

func (s *SchedulerService) SetTenantIDs(tenantIDs []uuid.UUID) {
	s.tenantIDs = tenantIDs
}

func (s *SchedulerService) Start() {
	s.logger.Info("Starting fleet notification scheduler")

	s.RegisterJob(Job{
		Name:     "check_expiring_licenses",
		Interval: 24 * time.Hour,
		Fn: func(ctx context.Context, tenantID uuid.UUID) error {
			return s.notificationSvc.CheckExpiringLicenses(ctx, tenantID, 30)
		},
	})

	s.RegisterJob(Job{
		Name:     "check_expiring_registrations",
		Interval: 24 * time.Hour,
		Fn: func(ctx context.Context, tenantID uuid.UUID) error {
			return s.notificationSvc.CheckExpiringRegistrations(ctx, tenantID, 30)
		},
	})

	s.RegisterJob(Job{
		Name:     "check_expiring_insurance",
		Interval: 24 * time.Hour,
		Fn: func(ctx context.Context, tenantID uuid.UUID) error {
			return s.notificationSvc.CheckExpiringInsurance(ctx, tenantID, 30)
		},
	})

	s.RegisterJob(Job{
		Name:     "check_due_maintenance",
		Interval: 24 * time.Hour,
		Fn: func(ctx context.Context, tenantID uuid.UUID) error {
			return s.notificationSvc.CheckDueMaintenance(ctx, tenantID)
		},
	})

	for _, job := range s.jobs {
		s.wg.Add(1)
		go s.runJob(job)
	}
}

func (s *SchedulerService) Stop() {
	s.logger.Info("Stopping fleet notification scheduler")
	s.cancel()
	s.wg.Wait()
}

func (s *SchedulerService) runJob(job Job) {
	defer s.wg.Done()

	ticker := time.NewTicker(job.Interval)
	defer ticker.Stop()

	s.executeJob(job)

	for {
		select {
		case <-s.ctx.Done():
			s.logger.WithField("job", job.Name).Info("Job stopped")
			return
		case <-ticker.C:
			s.executeJob(job)
		}
	}
}

func (s *SchedulerService) executeJob(job Job) {
	s.logger.WithField("job", job.Name).Info("Executing scheduled job")

	for _, tenantID := range s.tenantIDs {
		ctx := context.WithValue(s.ctx, "tenant_id", tenantID)

		if err := job.Fn(ctx, tenantID); err != nil {
			s.logger.WithError(err).WithFields(logrus.Fields{
				"job":       job.Name,
				"tenant_id": tenantID,
			}).Error("Job execution failed")
		}
	}
}

func (s *SchedulerService) RunOnce() {
	s.logger.Info("Running all fleet notification jobs once")

	for _, job := range s.jobs {
		s.executeJob(job)
	}
}
