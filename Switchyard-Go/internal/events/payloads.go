package events

import (
	"time"

	"github.com/google/uuid"
)

// Each payload type carries the minimum fields the target service needs.
// The Event envelope wraps these as json.RawMessage so the router can
// decode only what it needs for the matched event type.

type StopLoggedPayload struct {
	DriverID      uuid.UUID `json:"driver_id"`
	PlanBOLStopID uuid.UUID `json:"plan_bol_stop_id"`
	PlanBOLID     uuid.UUID `json:"plan_bol_id"`
	LocationID    string    `json:"location_id"`
	StopType      string    `json:"stop_type"`
	LoggedAt      time.Time `json:"logged_at"`
}

type MandatedStopPayload struct {
	DriverID     uuid.UUID `json:"driver_id"`
	AssignmentID uuid.UUID `json:"assignment_id"`
	StopAt       time.Time `json:"stop_at"`
	ELDStopRef   *string   `json:"eld_stop_ref,omitempty"`
}

type HOSAlertPayload struct {
	DriverID        uuid.UUID `json:"driver_id"`
	DailyHoursUsed  float64   `json:"daily_hours_used"`
	WeeklyHoursUsed float64   `json:"weekly_hours_used"`
	StateCode       string    `json:"state_code"`
	CycleLabel      string    `json:"cycle_label"`
}

type AssignmentPayload struct {
	AssignmentID uuid.UUID `json:"assignment_id"`
	DriverID     uuid.UUID `json:"driver_id"`
	PlanBOLID    uuid.UUID `json:"plan_bol_id"`
	EquipmentID  uuid.UUID `json:"equipment_id"`
}

type BOLCompletedPayload struct {
	AssignmentID uuid.UUID `json:"assignment_id"`
	DriverID     uuid.UUID `json:"driver_id"`
	PlanBOLID    uuid.UUID `json:"plan_bol_id"`
}

type DeadheadExpiryPayload struct {
	PairingID   uuid.UUID `json:"pairing_id"`
	ActiveBOLID uuid.UUID `json:"active_bol_id"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type EquipmentBreakdownPayload struct {
	EquipmentID   uuid.UUID  `json:"equipment_id"`
	BreakdownID   uuid.UUID  `json:"breakdown_id"`
	BreakdownType string     `json:"breakdown_type"`
	LocationDesc  *string    `json:"location_desc,omitempty"`
	DriverID      *uuid.UUID `json:"driver_id,omitempty"`
	LoadAttached  bool       `json:"load_attached"`
	ReportedAt    time.Time  `json:"reported_at"`
}

type EquipmentResolvedPayload struct {
	EquipmentID   uuid.UUID  `json:"equipment_id"`
	BreakdownID   *uuid.UUID `json:"breakdown_id,omitempty"`
	MaintenanceID *uuid.UUID `json:"maintenance_id,omitempty"`
	ResolvedAt    time.Time  `json:"resolved_at"`
}
