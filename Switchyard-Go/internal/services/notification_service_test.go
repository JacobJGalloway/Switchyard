package services

import (
	"context"
	"net/smtp"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JacobJGalloway/switchyard-go/internal/events"
)

// capturedMail holds the fields parsed from the raw SMTP message written by send().
type capturedMail struct {
	addr    string
	from    string
	to      []string
	subject string
	body    string
}

// noopMailer returns a mailerFunc that captures the outgoing message into *sent
// instead of opening an SMTP connection.
func noopMailer(sent *capturedMail) mailerFunc {
	return func(addr string, _ smtp.Auth, from string, to []string, msg []byte) error {
		raw := string(msg)
		lines := strings.Split(raw, "\r\n")
		var subject, body string
		for i, l := range lines {
			if strings.HasPrefix(l, "Subject: ") {
				subject = strings.TrimPrefix(l, "Subject: ")
			}
			if l == "" && i < len(lines)-1 {
				body = strings.Join(lines[i+1:], "\r\n")
				break
			}
		}
		*sent = capturedMail{addr: addr, from: from, to: to, subject: subject, body: body}
		return nil
	}
}

func newTestNotifSvc() (*NotificationService, *capturedMail) {
	cfg := NotificationConfig{
		SMTPHost:      "mail.test.invalid",
		SMTPPort:      "587",
		SMTPUser:      "from@test.invalid",
		SMTPPass:      "pass",
		DispatchEmail: "dispatch@test.invalid",
	}
	svc := NewNotificationService(cfg)
	sent := &capturedMail{}
	svc.mailer = noopMailer(sent)
	return svc, sent
}

func TestNotification_OnHOSLimitApproaching(t *testing.T) {
	svc, sent := newTestNotifSvc()
	driverID := uuid.New()
	err := svc.OnHOSLimitApproaching(context.Background(), events.HOSAlertPayload{
		DriverID:        driverID,
		StateCode:       "IL",
		CycleLabel:      "60h/7d",
		DailyHoursUsed:  9.5,
		WeeklyHoursUsed: 52.0,
	})
	require.NoError(t, err)
	assert.Contains(t, sent.subject, "HOS Warning")
	assert.Contains(t, sent.subject, "IL")
	assert.Contains(t, sent.body, driverID.String())
	assert.Equal(t, []string{"dispatch@test.invalid"}, sent.to)
}

func TestNotification_OnHOSWeeklyLimitReached(t *testing.T) {
	svc, sent := newTestNotifSvc()
	driverID := uuid.New()
	err := svc.OnHOSWeeklyLimitReached(context.Background(), events.HOSAlertPayload{
		DriverID:        driverID,
		StateCode:       "OH",
		CycleLabel:      "60h/7d",
		WeeklyHoursUsed: 60.0,
	})
	require.NoError(t, err)
	assert.Contains(t, sent.subject, "HOS Limit Reached")
	assert.Contains(t, sent.subject, driverID.String())
	assert.Contains(t, sent.body, "cannot be assigned")
	assert.Equal(t, []string{"dispatch@test.invalid"}, sent.to)
}

func TestNotification_OnBOLWorkflowCompleted(t *testing.T) {
	svc, sent := newTestNotifSvc()
	assignID := uuid.New()
	bolID := uuid.New()
	driverID := uuid.New()
	err := svc.OnBOLWorkflowCompleted(context.Background(), events.BOLCompletedPayload{
		AssignmentID: assignID,
		PlanBOLID:    bolID,
		DriverID:     driverID,
	})
	require.NoError(t, err)
	assert.Contains(t, sent.subject, "BOL Completed")
	assert.Contains(t, sent.subject, assignID.String())
	assert.Contains(t, sent.body, bolID.String())
	assert.Contains(t, sent.body, "dead-head")
}

func TestNotification_OnDeadheadWindowExpiring(t *testing.T) {
	svc, sent := newTestNotifSvc()
	pairingID := uuid.New()
	bolID := uuid.New()
	expiresAt := time.Now().UTC().Add(45 * time.Minute)
	err := svc.OnDeadheadWindowExpiring(context.Background(), events.DeadheadExpiryPayload{
		PairingID:   pairingID,
		ActiveBOLID: bolID,
		ExpiresAt:   expiresAt,
	})
	require.NoError(t, err)
	assert.Contains(t, sent.subject, "Dead-Head Window Expiring")
	assert.Contains(t, sent.subject, pairingID.String())
	assert.Contains(t, sent.body, bolID.String())
	assert.Contains(t, sent.body, "Confirm the dead-head pairing")
}

func TestNotification_OnRoadsideBreakdownWithLoad_WithDriver(t *testing.T) {
	svc, sent := newTestNotifSvc()
	equipID := uuid.New()
	breakID := uuid.New()
	driverID := uuid.New()
	loc := "I-90 mm 42"
	err := svc.OnRoadsideBreakdownWithLoad(context.Background(), events.EquipmentBreakdownPayload{
		EquipmentID:   equipID,
		BreakdownID:   breakID,
		BreakdownType: "roadside",
		LocationDesc:  &loc,
		DriverID:      &driverID,
		LoadAttached:  true,
		ReportedAt:    time.Now().UTC(),
	})
	require.NoError(t, err)
	assert.Contains(t, sent.subject, "URGENT")
	assert.Contains(t, sent.subject, equipID.String())
	assert.Contains(t, sent.body, "rescue dispatch required")
	assert.Contains(t, sent.body, driverID.String())
	assert.Contains(t, sent.body, loc)
}

func TestNotification_OnRoadsideBreakdownWithLoad_NilDriverNilLocation(t *testing.T) {
	svc, sent := newTestNotifSvc()
	equipID := uuid.New()
	breakID := uuid.New()
	err := svc.OnRoadsideBreakdownWithLoad(context.Background(), events.EquipmentBreakdownPayload{
		EquipmentID:   equipID,
		BreakdownID:   breakID,
		BreakdownType: "roadside",
		LocationDesc:  nil,
		DriverID:      nil,
		LoadAttached:  true,
		ReportedAt:    time.Now().UTC(),
	})
	require.NoError(t, err)
	assert.Contains(t, sent.subject, "URGENT")
	assert.NotContains(t, sent.body, "Driver:")
	assert.NotContains(t, sent.body, "Location:")
}
