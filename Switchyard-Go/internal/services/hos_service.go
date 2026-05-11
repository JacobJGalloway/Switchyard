package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/JacobJGalloway/switchyard-go/internal/events"
	"github.com/JacobJGalloway/switchyard-go/internal/models"
	"github.com/JacobJGalloway/switchyard-go/internal/repository"
)

// hosNotifier is the minimal alert interface HOSService needs.
// The concrete notification service satisfies this structurally.
type hosNotifier interface {
	OnHOSLimitApproaching(ctx context.Context, e events.HOSAlertPayload) error
	OnHOSWeeklyLimitReached(ctx context.Context, e events.HOSAlertPayload) error
}

type HOSService struct {
	hosRepo               repository.HOSRepository
	notifier              hosNotifier
	warningThresholdHours float64
}

func NewHOSService(
	hosRepo repository.HOSRepository,
	notifier hosNotifier,
	warningThresholdHours float64,
) *HOSService {
	return &HOSService{
		hosRepo:               hosRepo,
		notifier:              notifier,
		warningThresholdHours: warningThresholdHours,
	}
}

// OnStopLogged implements events.HOSService.
// Calculates leg driving hours, updates the window, and enforces all HOS rules.
// Returns an error if any hard constraint is violated — callers must not allow
// the driver to continue until the violation is resolved.
func (s *HOSService) OnStopLogged(ctx context.Context, e events.StopLoggedPayload) error {
	window, err := s.hosRepo.GetWindowByDriver(ctx, e.DriverID)
	if err != nil {
		return fmt.Errorf("getting HOS window for driver %s: %w", e.DriverID, err)
	}

	// Calculate driving hours for this leg.
	// LastActivityAt is set when the driver departs a stop (logged externally).
	// If nil (first stop of the day), no leg hours are added — stop 1 is
	// auto-processed at BOL creation and has no prior departure.
	if window.LastActivityAt != nil {
		legHours := e.LoggedAt.Sub(*window.LastActivityAt).Hours()
		if legHours < 0 {
			return fmt.Errorf("stop logged at %s is before last activity at %s",
				e.LoggedAt.Format(time.RFC3339), window.LastActivityAt.Format(time.RFC3339))
		}
		window.DailyHoursUsed += legHours
		window.WeeklyHoursUsed += legHours
	}

	window.LastActivityAt = &e.LoggedAt

	// Hard constraint: 30-minute break required at 8 cumulative driving hours.
	// The break must be taken before any further driving — reject the stop log
	// if the driver has hit 8 hours without taking it.
	if !window.Break30Taken && window.DailyHoursUsed >= 8.0 {
		return fmt.Errorf("driver %s must take a 30-minute break before continuing (%.1f cumulative driving hours)",
			e.DriverID, window.DailyHoursUsed)
	}

	if err := s.hosRepo.UpdateWindow(ctx, window); err != nil {
		return fmt.Errorf("updating HOS window for driver %s: %w", e.DriverID, err)
	}

	// Fetch limit for alert threshold checks.
	// Driver's state and cycle are needed — fetched here from the window's driver.
	// TODO: extend StopLoggedPayload or fetch driver record to get state + cycle_label.
	// For now, alert checks are skipped — wire them in when driver lookup is available.

	return nil
}

// RecordBreak records a completed 30-minute break for the driver.
// Called when the driver logs their rest period, not when they arrive at a stop.
func (s *HOSService) RecordBreak(ctx context.Context, driverID uuid.UUID, takenAt time.Time) error {
	window, err := s.hosRepo.GetWindowByDriver(ctx, driverID)
	if err != nil {
		return fmt.Errorf("getting HOS window: %w", err)
	}

	window.Break30Taken = true
	window.Break30At = &takenAt
	window.LastActivityAt = &takenAt

	return s.hosRepo.UpdateWindow(ctx, window)
}

// RecordMandatedStop records a required rest stop and pauses the dead-head timer.
func (s *HOSService) RecordMandatedStop(ctx context.Context, driverID uuid.UUID, stoppedAt time.Time, eldRef *string) error {
	window, err := s.hosRepo.GetWindowByDriver(ctx, driverID)
	if err != nil {
		return fmt.Errorf("getting HOS window: %w", err)
	}

	window.MandatedStopAt = &stoppedAt
	window.ELDStopRef = eldRef

	return s.hosRepo.UpdateWindow(ctx, window)
}

// CanAssign checks whether a driver can legally complete a planned run without
// exceeding their daily or weekly HOS limits. Hard constraint — returns an error
// if the driver cannot be assigned (see ARCHITECTURE.md §4.3).
func (s *HOSService) CanAssign(ctx context.Context, driverID uuid.UUID, estimatedRunHours float64, stateCode, cycleLabel string) error {
	window, err := s.hosRepo.GetWindowByDriver(ctx, driverID)
	if err != nil {
		return fmt.Errorf("getting HOS window: %w", err)
	}

	limit, err := s.hosRepo.GetLimitByStateAndCycle(ctx, stateCode, cycleLabel)
	if err != nil {
		return fmt.Errorf("getting HOS limit for %s/%s: %w", stateCode, cycleLabel, err)
	}

	projectedDaily := window.DailyHoursUsed + estimatedRunHours
	projectedWeekly := window.WeeklyHoursUsed + estimatedRunHours

	if projectedDaily > limit.DailyDrivingLimitHours {
		return fmt.Errorf("driver %s would exceed daily driving limit: %.1f projected vs %.1f allowed",
			driverID, projectedDaily, limit.DailyDrivingLimitHours)
	}

	if projectedWeekly > limit.WeeklyLimitHours {
		return fmt.Errorf("driver %s would exceed weekly limit: %.1f projected vs %.1f allowed",
			driverID, projectedWeekly, limit.WeeklyLimitHours)
	}

	return nil
}

// ResetWindow resets the daily and weekly clocks after a mandated rest period.
// Called when the required rest hours have been confirmed as completed.
func (s *HOSService) ResetWindow(ctx context.Context, driverID uuid.UUID, resetAt time.Time, resetWeekly bool) error {
	window, err := s.hosRepo.GetWindowByDriver(ctx, driverID)
	if err != nil {
		return fmt.Errorf("getting HOS window: %w", err)
	}

	window.WindowStart = resetAt
	window.DailyHoursUsed = 0
	window.Break30Taken = false
	window.Break30At = nil
	window.MandatedStopAt = nil
	window.LastActivityAt = nil

	if resetWeekly {
		window.WeeklyHoursUsed = 0
	}

	return s.hosRepo.UpdateWindow(ctx, window)
}

// checkAlerts fires notification events if the driver is approaching or has
// reached HOS limits. Non-fatal — a notification failure does not block the stop log.
func (s *HOSService) checkAlerts(ctx context.Context, window *models.HOSWindow, limit *models.HOSLimit) {
	remaining := limit.DailyDrivingLimitHours - window.DailyHoursUsed
	payload := events.HOSAlertPayload{
		DriverID:        window.DriverID,
		DailyHoursUsed:  window.DailyHoursUsed,
		WeeklyHoursUsed: window.WeeklyHoursUsed,
		StateCode:       limit.StateCode,
		CycleLabel:      limit.CycleLabel,
	}

	if window.WeeklyHoursUsed >= limit.WeeklyLimitHours {
		// Non-fatal: log if notification fails
		_ = s.notifier.OnHOSWeeklyLimitReached(ctx, payload)
		return
	}

	if remaining <= s.warningThresholdHours {
		_ = s.notifier.OnHOSLimitApproaching(ctx, payload)
	}
}
