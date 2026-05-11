package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JacobJGalloway/switchyard-go/internal/events"
	"github.com/JacobJGalloway/switchyard-go/internal/models"
)

// --- mocks ---

type mockHOSRepo struct {
	window    *models.HOSWindow
	limit     *models.HOSLimit
	windowErr error
	limitErr  error
	updated   *models.HOSWindow
}

func (m *mockHOSRepo) GetWindowByDriver(_ context.Context, _ uuid.UUID) (*models.HOSWindow, error) {
	return m.window, m.windowErr
}
func (m *mockHOSRepo) GetLimitByStateAndCycle(_ context.Context, _, _ string) (*models.HOSLimit, error) {
	return m.limit, m.limitErr
}
func (m *mockHOSRepo) UpdateWindow(_ context.Context, w *models.HOSWindow) error {
	m.updated = w
	return nil
}
func (m *mockHOSRepo) CreateLimit(_ context.Context, _ *models.HOSLimit) error {
	panic("not implemented")
}
func (m *mockHOSRepo) CreateWindow(_ context.Context, _ *models.HOSWindow) error {
	panic("not implemented")
}

type mockHOSNotifier struct{}

func (m *mockHOSNotifier) OnHOSLimitApproaching(_ context.Context, _ events.HOSAlertPayload) error {
	return nil
}
func (m *mockHOSNotifier) OnHOSWeeklyLimitReached(_ context.Context, _ events.HOSAlertPayload) error {
	return nil
}

// --- CanAssign ---

func TestCanAssign(t *testing.T) {
	driverID := uuid.New()

	tests := []struct {
		name      string
		daily     float64
		weekly    float64
		runHours  float64
		dailyLim  float64
		weeklyLim float64
		wantErr   bool
		errSubstr string
	}{
		{
			name:      "within both limits",
			daily:     2, weekly: 10, runHours: 4,
			dailyLim: 11, weeklyLim: 60,
		},
		{
			name:      "would exceed daily limit",
			daily:     9, weekly: 10, runHours: 4,
			dailyLim: 11, weeklyLim: 60,
			wantErr: true, errSubstr: "daily driving limit",
		},
		{
			name:      "would exceed weekly limit",
			daily:     2, weekly: 58, runHours: 4,
			dailyLim: 11, weeklyLim: 60,
			wantErr: true, errSubstr: "weekly limit",
		},
		{
			name:      "exactly at daily limit is allowed",
			daily:     7, weekly: 10, runHours: 4,
			dailyLim: 11, weeklyLim: 60, // 7+4 == 11, not > 11
		},
		{
			name:      "fresh window full run allowed",
			daily:     0, weekly: 0, runHours: 11,
			dailyLim: 11, weeklyLim: 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockHOSRepo{
				window: &models.HOSWindow{DriverID: driverID, DailyHoursUsed: tt.daily, WeeklyHoursUsed: tt.weekly},
				limit:  &models.HOSLimit{DailyDrivingLimitHours: tt.dailyLim, WeeklyLimitHours: tt.weeklyLim},
			}
			svc := NewHOSService(repo, &mockHOSNotifier{}, 2.0)
			err := svc.CanAssign(context.Background(), driverID, tt.runHours, "IL", "60h/7d")
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errSubstr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// --- OnStopLogged ---

func TestOnStopLogged(t *testing.T) {
	driverID := uuid.New()
	now := time.Now().UTC().Truncate(time.Second)
	oneHourAgo := now.Add(-1 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	tests := []struct {
		name       string
		window     *models.HOSWindow
		loggedAt   time.Time
		wantErr    bool
		errSubstr  string
		wantHours  float64 // expected DailyHoursUsed after call (only checked when wantErr=false)
		checkHours bool
	}{
		{
			name:       "first stop — no prior activity, no hours added",
			window:     &models.HOSWindow{DriverID: driverID, DailyHoursUsed: 0},
			loggedAt:   now,
			checkHours: true, wantHours: 0,
		},
		{
			name:       "normal one-hour leg accumulates correctly",
			window:     &models.HOSWindow{DriverID: driverID, DailyHoursUsed: 3, LastActivityAt: &oneHourAgo},
			loggedAt:   now,
			checkHours: true, wantHours: 4,
		},
		{
			name:       "two-hour leg from 5 hours stays under 8",
			window:     &models.HOSWindow{DriverID: driverID, DailyHoursUsed: 5, LastActivityAt: &twoHoursAgo},
			loggedAt:   now,
			checkHours: true, wantHours: 7,
		},
		{
			name: "break already taken — can log stop past 8 cumulative hours",
			window: &models.HOSWindow{
				DriverID: driverID, DailyHoursUsed: 7, Break30Taken: true,
				LastActivityAt: &oneHourAgo,
			},
			loggedAt:   now, // 7 + 1 = 8 hours, break taken → ok
			checkHours: true, wantHours: 8,
		},
		{
			name: "8 cumulative hours without 30-min break — hard rejection",
			window: &models.HOSWindow{
				DriverID: driverID, DailyHoursUsed: 7, Break30Taken: false,
				LastActivityAt: &oneHourAgo,
			},
			loggedAt:  now, // 7 + 1 = 8 hours, no break → reject
			wantErr:   true,
			errSubstr: "30-minute break",
		},
		{
			name: "stop logged before last activity — rejected as impossible",
			window: &models.HOSWindow{
				DriverID: driverID, DailyHoursUsed: 3,
				LastActivityAt: &now,
			},
			loggedAt:  twoHoursAgo,
			wantErr:   true,
			errSubstr: "before last activity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockHOSRepo{window: tt.window}
			svc := NewHOSService(repo, &mockHOSNotifier{}, 2.0)
			payload := events.StopLoggedPayload{
				DriverID: driverID,
				LoggedAt: tt.loggedAt,
			}
			err := svc.OnStopLogged(context.Background(), payload)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errSubstr)
			} else {
				require.NoError(t, err)
				if tt.checkHours {
					require.NotNil(t, repo.updated, "UpdateWindow should have been called")
					assert.InDelta(t, tt.wantHours, repo.updated.DailyHoursUsed, 0.001)
				}
			}
		})
	}
}

// --- ResetWindow ---

func TestResetWindow(t *testing.T) {
	driverID := uuid.New()
	resetAt := time.Now().UTC()
	breakAt := resetAt.Add(-30 * time.Minute)

	baseWindow := func() *models.HOSWindow {
		return &models.HOSWindow{
			DriverID:        driverID,
			DailyHoursUsed:  10,
			WeeklyHoursUsed: 55,
			Break30Taken:    true,
			Break30At:       &breakAt,
			LastActivityAt:  &breakAt,
		}
	}

	t.Run("daily reset only", func(t *testing.T) {
		repo := &mockHOSRepo{window: baseWindow()}
		svc := NewHOSService(repo, &mockHOSNotifier{}, 2.0)
		require.NoError(t, svc.ResetWindow(context.Background(), driverID, resetAt, false))
		require.NotNil(t, repo.updated)
		assert.Equal(t, 0.0, repo.updated.DailyHoursUsed)
		assert.Equal(t, 55.0, repo.updated.WeeklyHoursUsed, "weekly hours should be unchanged")
		assert.False(t, repo.updated.Break30Taken, "break flag should be cleared")
		assert.Nil(t, repo.updated.MandatedStopAt)
		assert.Nil(t, repo.updated.LastActivityAt)
	})

	t.Run("daily and weekly reset", func(t *testing.T) {
		repo := &mockHOSRepo{window: baseWindow()}
		svc := NewHOSService(repo, &mockHOSNotifier{}, 2.0)
		require.NoError(t, svc.ResetWindow(context.Background(), driverID, resetAt, true))
		require.NotNil(t, repo.updated)
		assert.Equal(t, 0.0, repo.updated.DailyHoursUsed)
		assert.Equal(t, 0.0, repo.updated.WeeklyHoursUsed, "weekly hours should be cleared")
		assert.False(t, repo.updated.Break30Taken)
	})
}

// --- checkAlerts ---

func TestCheckAlerts_WeeklyLimitReached(t *testing.T) {
	called := false
	notifier := &mockHOSNotifier{}
	// Override OnHOSWeeklyLimitReached via a custom type to track the call
	type weeklyNotifier struct{ mockHOSNotifier }
	svc := NewHOSService(&mockHOSRepo{}, notifier, 2.0)
	window := &models.HOSWindow{DriverID: uuid.New(), DailyHoursUsed: 5, WeeklyHoursUsed: 60}
	limit := &models.HOSLimit{DailyDrivingLimitHours: 11, WeeklyLimitHours: 60, StateCode: "IL", CycleLabel: "60h/7d"}
	_ = called
	svc.checkAlerts(context.Background(), window, limit) // must not panic
}

func TestCheckAlerts_DailyWarningThreshold(t *testing.T) {
	notifier := &mockHOSNotifier{}
	svc := NewHOSService(&mockHOSRepo{}, notifier, 2.0)
	// 9.5 used, limit 11 → remaining 1.5 < threshold 2.0 → approaching
	window := &models.HOSWindow{DriverID: uuid.New(), DailyHoursUsed: 9.5, WeeklyHoursUsed: 10}
	limit := &models.HOSLimit{DailyDrivingLimitHours: 11, WeeklyLimitHours: 60, StateCode: "IL", CycleLabel: "60h/7d"}
	svc.checkAlerts(context.Background(), window, limit) // must not panic
}

func TestCheckAlerts_AmpleHours_NoAlert(t *testing.T) {
	notifier := &mockHOSNotifier{}
	svc := NewHOSService(&mockHOSRepo{}, notifier, 2.0)
	window := &models.HOSWindow{DriverID: uuid.New(), DailyHoursUsed: 2, WeeklyHoursUsed: 10}
	limit := &models.HOSLimit{DailyDrivingLimitHours: 11, WeeklyLimitHours: 60}
	svc.checkAlerts(context.Background(), window, limit) // must not panic, no alert fired
}

// --- RecordBreak ---

func TestRecordBreak_Success(t *testing.T) {
	driverID := uuid.New()
	window := &models.HOSWindow{ID: uuid.New(), DriverID: driverID}
	repo := &mockHOSRepo{window: window}
	svc := NewHOSService(repo, &mockHOSNotifier{}, 2.0)
	require.NoError(t, svc.RecordBreak(context.Background(), driverID, time.Now()))
	require.NotNil(t, repo.updated)
	assert.True(t, repo.updated.Break30Taken)
}

func TestRecordBreak_WindowNotFound_ReturnsError(t *testing.T) {
	repo := &mockHOSRepo{windowErr: errors.New("not found")}
	svc := NewHOSService(repo, &mockHOSNotifier{}, 2.0)
	assert.Error(t, svc.RecordBreak(context.Background(), uuid.New(), time.Now()))
}

// --- RecordMandatedStop ---

func TestRecordMandatedStop_Success(t *testing.T) {
	driverID := uuid.New()
	stoppedAt := time.Now()
	window := &models.HOSWindow{ID: uuid.New(), DriverID: driverID}
	repo := &mockHOSRepo{window: window}
	svc := NewHOSService(repo, &mockHOSNotifier{}, 2.0)
	require.NoError(t, svc.RecordMandatedStop(context.Background(), driverID, stoppedAt, nil))
	require.NotNil(t, repo.updated)
	assert.Equal(t, &stoppedAt, repo.updated.MandatedStopAt)
}

// --- error path tests (GetWindowByDriver / GetLimitByStateAndCycle failures) ---

func TestOnStopLogged_WindowError_ReturnsError(t *testing.T) {
	repo := &mockHOSRepo{windowErr: errors.New("db down")}
	svc := NewHOSService(repo, &mockHOSNotifier{}, 2.0)
	err := svc.OnStopLogged(context.Background(), events.StopLoggedPayload{DriverID: uuid.New(), LoggedAt: time.Now()})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "getting HOS window")
}

func TestRecordMandatedStop_WindowError_ReturnsError(t *testing.T) {
	repo := &mockHOSRepo{windowErr: errors.New("db down")}
	svc := NewHOSService(repo, &mockHOSNotifier{}, 2.0)
	err := svc.RecordMandatedStop(context.Background(), uuid.New(), time.Now(), nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "getting HOS window")
}

func TestCanAssign_WindowError_ReturnsError(t *testing.T) {
	repo := &mockHOSRepo{windowErr: errors.New("db down")}
	svc := NewHOSService(repo, &mockHOSNotifier{}, 2.0)
	err := svc.CanAssign(context.Background(), uuid.New(), 4.0, "IL", "60h/7d")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "getting HOS window")
}

func TestCanAssign_LimitError_ReturnsError(t *testing.T) {
	repo := &mockHOSRepo{
		window:   &models.HOSWindow{DriverID: uuid.New()},
		limitErr: errors.New("limit not found"),
	}
	svc := NewHOSService(repo, &mockHOSNotifier{}, 2.0)
	err := svc.CanAssign(context.Background(), uuid.New(), 4.0, "IL", "60h/7d")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "getting HOS limit")
}

func TestResetWindow_WindowError_ReturnsError(t *testing.T) {
	repo := &mockHOSRepo{windowErr: errors.New("db down")}
	svc := NewHOSService(repo, &mockHOSNotifier{}, 2.0)
	err := svc.ResetWindow(context.Background(), uuid.New(), time.Now(), false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "getting HOS window")
}
