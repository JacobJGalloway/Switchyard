package models

import (
	"time"

	"github.com/google/uuid"
)

// DriverBOLAssignment links a driver, a planned BOL, and a piece of equipment
// for a single run. Answers: who is taking what, in which vehicle, to where.
type DriverBOLAssignment struct {
	ID                  uuid.UUID  `json:"id"`
	DriverID            uuid.UUID  `json:"driver_id"`
	PlanBOLID           uuid.UUID  `json:"plan_bol_id"`
	EquipmentID         uuid.UUID  `json:"equipment_id"`
	BaseRatePerMile     float64    `json:"base_rate_per_mile"`
	AssignedAt          time.Time  `json:"assigned_at"`
	DepartedAt          *time.Time `json:"departed_at,omitempty"`
	FulfilledAt         *time.Time `json:"fulfilled_at,omitempty"`
	DeadheadConfirmedAt *time.Time `json:"deadhead_confirmed_at,omitempty"`
}
