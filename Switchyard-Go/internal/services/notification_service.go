package services

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/JacobJGalloway/switchyard-go/internal/events"
)

// NotificationConfig holds the SMTP and dispatch routing settings.
// All fields map directly to env vars (see ARCHITECTURE.md §11).
type NotificationConfig struct {
	SMTPHost      string
	SMTPPort      string
	SMTPUser      string
	SMTPPass      string
	DispatchEmail string
}

// mailerFunc is the send signature of smtp.SendMail, injected so tests can
// capture outgoing mail without opening an SMTP connection.
type mailerFunc func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

// NotificationService sends operational email alerts to the dispatch team.
// It satisfies both events.NotificationService (5 trigger methods wired through
// the event handler) and the hosNotifier interface in hos_service.go.
// Notification failures are non-fatal — callers discard the error so that
// a failed email never blocks a driver workflow event.
type NotificationService struct {
	cfg    NotificationConfig
	mailer mailerFunc
}

func NewNotificationService(cfg NotificationConfig) *NotificationService {
	return &NotificationService{cfg: cfg, mailer: smtp.SendMail}
}

// OnHOSLimitApproaching fires when a driver's remaining daily hours fall within
// the warning threshold (HOS_WARNING_THRESHOLD_HOURS env var).
func (s *NotificationService) OnHOSLimitApproaching(_ context.Context, e events.HOSAlertPayload) error {
	subject := fmt.Sprintf("HOS Warning — driver approaching limit (%s / %s)", e.StateCode, e.CycleLabel)
	body := fmt.Sprintf(
		"Driver %s is approaching their HOS driving limit.\r\n\r\nDaily hours used:  %.1f\r\nWeekly hours used: %.1f\r\nState: %s  Cycle: %s",
		e.DriverID, e.DailyHoursUsed, e.WeeklyHoursUsed, e.StateCode, e.CycleLabel,
	)
	return s.send(subject, body)
}

// OnHOSWeeklyLimitReached fires when a driver hits their weekly HOS cap.
// The driver cannot be assigned to new runs until their window resets.
func (s *NotificationService) OnHOSWeeklyLimitReached(_ context.Context, e events.HOSAlertPayload) error {
	subject := fmt.Sprintf("HOS Limit Reached — driver %s unavailable", e.DriverID)
	body := fmt.Sprintf(
		"Driver %s has reached their weekly HOS limit and cannot be assigned to new runs.\r\n\r\nWeekly hours used: %.1f\r\nState: %s  Cycle: %s\r\n\r\nDriver will appear in the HOS Limited column on the dispatch board.",
		e.DriverID, e.WeeklyHoursUsed, e.StateCode, e.CycleLabel,
	)
	return s.send(subject, body)
}

// OnBOLWorkflowCompleted fires when all stops on an active BOL are confirmed.
// The dead-head return window opens at this point.
func (s *NotificationService) OnBOLWorkflowCompleted(_ context.Context, e events.BOLCompletedPayload) error {
	subject := fmt.Sprintf("BOL Completed — assignment %s", e.AssignmentID)
	body := fmt.Sprintf(
		"All stops have been confirmed for BOL %s.\r\n\r\nDriver:     %s\r\nAssignment: %s\r\n\r\nThe dead-head return window is now open. Arrange a return run before the timer expires.",
		e.PlanBOLID, e.DriverID, e.AssignmentID,
	)
	return s.send(subject, body)
}

// OnDeadheadWindowExpiring fires when the dead-head search window drops below
// one hour. Dispatcher must confirm a pairing or the driver goes idle.
func (s *NotificationService) OnDeadheadWindowExpiring(_ context.Context, e events.DeadheadExpiryPayload) error {
	subject := fmt.Sprintf("Dead-Head Window Expiring — pairing %s", e.PairingID)
	body := fmt.Sprintf(
		"The dead-head return window for active BOL %s is expiring.\r\n\r\nPairing ID: %s\r\nExpires at: %s\r\n\r\nConfirm the dead-head pairing now to avoid the driver going idle.",
		e.ActiveBOLID, e.PairingID, e.ExpiresAt.Format("2006-01-02 15:04 MST"),
	)
	return s.send(subject, body)
}

// OnRoadsideBreakdownWithLoad fires when equipment breaks down in the field
// with a load still attached. This is the only critical-severity notification —
// rescue dispatch is required immediately.
func (s *NotificationService) OnRoadsideBreakdownWithLoad(_ context.Context, e events.EquipmentBreakdownPayload) error {
	subject := fmt.Sprintf("URGENT: Roadside Breakdown with Load — equipment %s", e.EquipmentID)

	var lines []string
	lines = append(lines, fmt.Sprintf("Roadside breakdown reported with load still attached — rescue dispatch required."))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Equipment:    %s", e.EquipmentID))
	lines = append(lines, fmt.Sprintf("Breakdown ID: %s", e.BreakdownID))
	if e.DriverID != nil {
		lines = append(lines, fmt.Sprintf("Driver:       %s", *e.DriverID))
	}
	if e.LocationDesc != nil {
		lines = append(lines, fmt.Sprintf("Location:     %s", *e.LocationDesc))
	}
	lines = append(lines, fmt.Sprintf("Reported at:  %s", e.ReportedAt.Format("2006-01-02 15:04 MST")))

	return s.send(subject, strings.Join(lines, "\r\n"))
}

// send delivers a plain-text email to the configured dispatch address via SMTP.
// Uses STARTTLS negotiation (port 587) with PLAIN auth.
func (s *NotificationService) send(subject, body string) error {
	addr := s.cfg.SMTPHost + ":" + s.cfg.SMTPPort
	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)

	msg := strings.Join([]string{
		"From: " + s.cfg.SMTPUser,
		"To: " + s.cfg.DispatchEmail,
		"Subject: [Switchyard] " + subject,
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	return s.mailer(addr, auth, s.cfg.SMTPUser, []string{s.cfg.DispatchEmail}, []byte(msg))
}
