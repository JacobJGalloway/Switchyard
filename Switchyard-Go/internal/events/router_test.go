package events

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- minimal mocks ---

type mockHOS struct{ called bool }

func (m *mockHOS) OnStopLogged(_ context.Context, _ StopLoggedPayload) error {
	m.called = true
	return nil
}

type mockWB struct {
	departed  bool
	fulfilled bool
	deadhead  bool
	mandated  bool
	breakdown bool
	resolved  bool
}

func (m *mockWB) OnAssignmentDeparted(_ context.Context, _ AssignmentPayload) error {
	m.departed = true
	return nil
}
func (m *mockWB) OnAssignmentFulfilled(_ context.Context, _ AssignmentPayload) error {
	m.fulfilled = true
	return nil
}
func (m *mockWB) OnDeadheadConfirmed(_ context.Context, _ AssignmentPayload) error {
	m.deadhead = true
	return nil
}
func (m *mockWB) OnMandatedStop(_ context.Context, _ MandatedStopPayload) error {
	m.mandated = true
	return nil
}
func (m *mockWB) OnEquipmentBreakdown(_ context.Context, _ EquipmentBreakdownPayload) error {
	m.breakdown = true
	return nil
}
func (m *mockWB) OnEquipmentResolved(_ context.Context, _ EquipmentResolvedPayload) error {
	m.resolved = true
	return nil
}

type mockNotify struct {
	hosApproaching bool
	hosWeekly      bool
	bolCompleted   bool
	deadheadExpiry bool
	roadsideBreak  bool
}

func (m *mockNotify) OnHOSLimitApproaching(_ context.Context, _ HOSAlertPayload) error {
	m.hosApproaching = true
	return nil
}
func (m *mockNotify) OnHOSWeeklyLimitReached(_ context.Context, _ HOSAlertPayload) error {
	m.hosWeekly = true
	return nil
}
func (m *mockNotify) OnBOLWorkflowCompleted(_ context.Context, _ BOLCompletedPayload) error {
	m.bolCompleted = true
	return nil
}
func (m *mockNotify) OnDeadheadWindowExpiring(_ context.Context, _ DeadheadExpiryPayload) error {
	m.deadheadExpiry = true
	return nil
}
func (m *mockNotify) OnRoadsideBreakdownWithLoad(_ context.Context, _ EquipmentBreakdownPayload) error {
	m.roadsideBreak = true
	return nil
}

// wireHandler builds a Handler with the given mocks (no HTTP client or Auth0 config needed).
func wireHandler(hos *mockHOS, wb *mockWB, notify *mockNotify) *Handler {
	return &Handler{
		hosService:          hos,
		whiteboardService:   wb,
		notificationService: notify,
	}
}

func marshalPayload(t testing.TB, p any) json.RawMessage {
	t.Helper()
	raw, err := json.Marshal(p)
	require.NoError(t, err)
	return raw
}

// --- per-event routing ---

func TestRoute_StopLogged(t *testing.T) {
	hos := &mockHOS{}
	h := wireHandler(hos, &mockWB{}, &mockNotify{})
	evt := Event{Type: EventStopLogged, OccurredAt: time.Now(),
		Payload: marshalPayload(t, StopLoggedPayload{DriverID: uuid.New(), LoggedAt: time.Now()})}
	require.NoError(t, route(context.Background(), h, evt))
	assert.True(t, hos.called)
}

func TestRoute_MandatedStop(t *testing.T) {
	wb := &mockWB{}
	h := wireHandler(&mockHOS{}, wb, &mockNotify{})
	evt := Event{Type: EventMandatedStop, OccurredAt: time.Now(),
		Payload: marshalPayload(t, MandatedStopPayload{DriverID: uuid.New()})}
	require.NoError(t, route(context.Background(), h, evt))
	assert.True(t, wb.mandated)
}

func TestRoute_HOSLimitApproaching(t *testing.T) {
	notify := &mockNotify{}
	h := wireHandler(&mockHOS{}, &mockWB{}, notify)
	evt := Event{Type: EventHOSLimitApproaching, OccurredAt: time.Now(),
		Payload: marshalPayload(t, HOSAlertPayload{DriverID: uuid.New()})}
	require.NoError(t, route(context.Background(), h, evt))
	assert.True(t, notify.hosApproaching)
}

func TestRoute_HOSWeeklyLimitReached(t *testing.T) {
	notify := &mockNotify{}
	h := wireHandler(&mockHOS{}, &mockWB{}, notify)
	evt := Event{Type: EventHOSWeeklyLimitReached, OccurredAt: time.Now(),
		Payload: marshalPayload(t, HOSAlertPayload{DriverID: uuid.New()})}
	require.NoError(t, route(context.Background(), h, evt))
	assert.True(t, notify.hosWeekly)
}

func TestRoute_AssignmentDeparted(t *testing.T) {
	wb := &mockWB{}
	h := wireHandler(&mockHOS{}, wb, &mockNotify{})
	evt := Event{Type: EventAssignmentDeparted, OccurredAt: time.Now(),
		Payload: marshalPayload(t, AssignmentPayload{AssignmentID: uuid.New()})}
	require.NoError(t, route(context.Background(), h, evt))
	assert.True(t, wb.departed)
}

// EventAssignmentFulfilled must call BOTH whiteboard AND notification (workflow completed).
func TestRoute_AssignmentFulfilled_CallsBothServices(t *testing.T) {
	wb := &mockWB{}
	notify := &mockNotify{}
	h := wireHandler(&mockHOS{}, wb, notify)
	evt := Event{Type: EventAssignmentFulfilled, OccurredAt: time.Now(),
		Payload: marshalPayload(t, AssignmentPayload{AssignmentID: uuid.New(), DriverID: uuid.New(), PlanBOLID: uuid.New()})}
	require.NoError(t, route(context.Background(), h, evt))
	assert.True(t, wb.fulfilled, "whiteboard must be updated on fulfillment")
	assert.True(t, notify.bolCompleted, "notification service must fire BOL-completed alert")
}

func TestRoute_DeadheadConfirmed(t *testing.T) {
	wb := &mockWB{}
	h := wireHandler(&mockHOS{}, wb, &mockNotify{})
	evt := Event{Type: EventDeadheadConfirmed, OccurredAt: time.Now(),
		Payload: marshalPayload(t, AssignmentPayload{AssignmentID: uuid.New()})}
	require.NoError(t, route(context.Background(), h, evt))
	assert.True(t, wb.deadhead)
}

func TestRoute_DeadheadWindowExpiring(t *testing.T) {
	notify := &mockNotify{}
	h := wireHandler(&mockHOS{}, &mockWB{}, notify)
	evt := Event{Type: EventDeadheadWindowExpiring, OccurredAt: time.Now(),
		Payload: marshalPayload(t, DeadheadExpiryPayload{PairingID: uuid.New(), ActiveBOLID: uuid.New(), ExpiresAt: time.Now()})}
	require.NoError(t, route(context.Background(), h, evt))
	assert.True(t, notify.deadheadExpiry)
}

func TestRoute_BOLWorkflowCompleted(t *testing.T) {
	notify := &mockNotify{}
	h := wireHandler(&mockHOS{}, &mockWB{}, notify)
	evt := Event{Type: EventBOLWorkflowCompleted, OccurredAt: time.Now(),
		Payload: marshalPayload(t, BOLCompletedPayload{AssignmentID: uuid.New()})}
	require.NoError(t, route(context.Background(), h, evt))
	assert.True(t, notify.bolCompleted)
}

// Roadside breakdown with load attached must update whiteboard AND alert dispatcher.
func TestRoute_EquipmentBreakdown_RoadsideWithLoad(t *testing.T) {
	wb := &mockWB{}
	notify := &mockNotify{}
	h := wireHandler(&mockHOS{}, wb, notify)
	evt := Event{Type: EventEquipmentBreakdown, OccurredAt: time.Now(),
		Payload: marshalPayload(t, EquipmentBreakdownPayload{
			EquipmentID: uuid.New(), BreakdownType: "roadside", LoadAttached: true,
		})}
	require.NoError(t, route(context.Background(), h, evt))
	assert.True(t, wb.breakdown, "whiteboard must be updated for any breakdown")
	assert.True(t, notify.roadsideBreak, "dispatcher must be alerted for roadside+load")
}

// Depot breakdown without load must only update the whiteboard — no dispatcher alert.
func TestRoute_EquipmentBreakdown_DepotNoLoad(t *testing.T) {
	wb := &mockWB{}
	notify := &mockNotify{}
	h := wireHandler(&mockHOS{}, wb, notify)
	evt := Event{Type: EventEquipmentBreakdown, OccurredAt: time.Now(),
		Payload: marshalPayload(t, EquipmentBreakdownPayload{
			EquipmentID: uuid.New(), BreakdownType: "depot", LoadAttached: false,
		})}
	require.NoError(t, route(context.Background(), h, evt))
	assert.True(t, wb.breakdown)
	assert.False(t, notify.roadsideBreak, "depot breakdown must not trigger immediate alert")
}

func TestRoute_EquipmentResolved(t *testing.T) {
	wb := &mockWB{}
	h := wireHandler(&mockHOS{}, wb, &mockNotify{})
	evt := Event{Type: EventEquipmentResolved, OccurredAt: time.Now(),
		Payload: marshalPayload(t, EquipmentResolvedPayload{EquipmentID: uuid.New(), ResolvedAt: time.Now()})}
	require.NoError(t, route(context.Background(), h, evt))
	assert.True(t, wb.resolved)
}

// --- error paths ---

func TestRoute_UnknownEventType(t *testing.T) {
	h := wireHandler(&mockHOS{}, &mockWB{}, &mockNotify{})
	err := route(context.Background(), h, Event{Type: "unknown.event", Payload: json.RawMessage(`{}`)})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown event type")
}

func TestRoute_BadJSONPayload(t *testing.T) {
	h := wireHandler(&mockHOS{}, &mockWB{}, &mockNotify{})
	err := route(context.Background(), h, Event{Type: EventStopLogged, Payload: json.RawMessage(`not-json`)})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing")
}
