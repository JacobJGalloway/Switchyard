package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Event is the envelope for all workflow events entering the Go backend.
type Event struct {
	Type       string          `json:"type"`
	OccurredAt time.Time       `json:"occurred_at"`
	Payload    json.RawMessage `json:"payload"`
}

// Event type constants follow <domain>.<action> naming.
const (
	EventStopLogged            = "driver.stop_logged"
	EventMandatedStop          = "driver.mandated_stop"
	EventHOSLimitApproaching   = "driver.hos_limit_approaching"
	EventHOSWeeklyLimitReached = "driver.hos_weekly_limit_reached"

	EventAssignmentDeparted     = "assignment.departed"
	EventAssignmentFulfilled    = "assignment.fulfilled"
	EventDeadheadConfirmed      = "assignment.deadhead_confirmed"
	EventDeadheadWindowExpiring = "assignment.deadhead_window_expiring"

	EventBOLWorkflowCompleted = "bol.workflow_completed"

	EventEquipmentBreakdown     = "equipment.breakdown_reported"
	EventEquipmentResolved      = "equipment.resolved"
	EventMaintenanceScheduled   = "equipment.maintenance_scheduled"
)

// route dispatches an event to the appropriate service handler(s).
// Events that touch multiple services call them in dependency order.
// If a later call fails after an earlier one succeeds, the error is returned
// but the earlier state change is not rolled back — acceptable for v1.1.
func route(ctx context.Context, h *Handler, evt Event) error {
	switch evt.Type {

	case EventStopLogged:
		var p StopLoggedPayload
		if err := json.Unmarshal(evt.Payload, &p); err != nil {
			return fmt.Errorf("parsing %s payload: %w", evt.Type, err)
		}
		return h.hosService.OnStopLogged(ctx, p)

	case EventMandatedStop:
		var p MandatedStopPayload
		if err := json.Unmarshal(evt.Payload, &p); err != nil {
			return fmt.Errorf("parsing %s payload: %w", evt.Type, err)
		}
		return h.whiteboardService.OnMandatedStop(ctx, p)

	case EventHOSLimitApproaching:
		var p HOSAlertPayload
		if err := json.Unmarshal(evt.Payload, &p); err != nil {
			return fmt.Errorf("parsing %s payload: %w", evt.Type, err)
		}
		return h.notificationService.OnHOSLimitApproaching(ctx, p)

	case EventHOSWeeklyLimitReached:
		var p HOSAlertPayload
		if err := json.Unmarshal(evt.Payload, &p); err != nil {
			return fmt.Errorf("parsing %s payload: %w", evt.Type, err)
		}
		return h.notificationService.OnHOSWeeklyLimitReached(ctx, p)

	case EventAssignmentDeparted:
		var p AssignmentPayload
		if err := json.Unmarshal(evt.Payload, &p); err != nil {
			return fmt.Errorf("parsing %s payload: %w", evt.Type, err)
		}
		return h.whiteboardService.OnAssignmentDeparted(ctx, p)

	case EventAssignmentFulfilled:
		var p AssignmentPayload
		if err := json.Unmarshal(evt.Payload, &p); err != nil {
			return fmt.Errorf("parsing %s payload: %w", evt.Type, err)
		}
		if err := h.whiteboardService.OnAssignmentFulfilled(ctx, p); err != nil {
			return err
		}
		return h.notificationService.OnBOLWorkflowCompleted(ctx, BOLCompletedPayload{
			AssignmentID: p.AssignmentID,
			DriverID:     p.DriverID,
			PlanBOLID:    p.PlanBOLID,
		})

	case EventDeadheadConfirmed:
		var p AssignmentPayload
		if err := json.Unmarshal(evt.Payload, &p); err != nil {
			return fmt.Errorf("parsing %s payload: %w", evt.Type, err)
		}
		return h.whiteboardService.OnDeadheadConfirmed(ctx, p)

	case EventDeadheadWindowExpiring:
		var p DeadheadExpiryPayload
		if err := json.Unmarshal(evt.Payload, &p); err != nil {
			return fmt.Errorf("parsing %s payload: %w", evt.Type, err)
		}
		return h.notificationService.OnDeadheadWindowExpiring(ctx, p)

	case EventBOLWorkflowCompleted:
		var p BOLCompletedPayload
		if err := json.Unmarshal(evt.Payload, &p); err != nil {
			return fmt.Errorf("parsing %s payload: %w", evt.Type, err)
		}
		return h.notificationService.OnBOLWorkflowCompleted(ctx, p)

	case EventEquipmentBreakdown:
		var p EquipmentBreakdownPayload
		if err := json.Unmarshal(evt.Payload, &p); err != nil {
			return fmt.Errorf("parsing %s payload: %w", evt.Type, err)
		}
		if err := h.whiteboardService.OnEquipmentBreakdown(ctx, p); err != nil {
			return err
		}
		// Roadside breakdowns with load attached require immediate dispatcher alert.
		if p.BreakdownType == "roadside" && p.LoadAttached {
			return h.notificationService.OnRoadsideBreakdownWithLoad(ctx, p)
		}
		return nil

	case EventEquipmentResolved:
		var p EquipmentResolvedPayload
		if err := json.Unmarshal(evt.Payload, &p); err != nil {
			return fmt.Errorf("parsing %s payload: %w", evt.Type, err)
		}
		return h.whiteboardService.OnEquipmentResolved(ctx, p)

	default:
		return fmt.Errorf("unknown event type: %s", evt.Type)
	}
}
