package models

import (
	"time"

	"github.com/google/uuid"
)

type StopType string

const (
	StopTypeWarehouse   StopType = "warehouse"
	StopTypeStore       StopType = "store"
	StopTypeReturnDepot StopType = "return_depot"
)

type PlanBOLStatus string

const (
	PlanBOLStatusDraft        PlanBOLStatus = "draft"
	PlanBOLStatusPlanProgress PlanBOLStatus = "plan-progress"
	PlanBOLStatusLoading      PlanBOLStatus = "loading"
	PlanBOLStatusValidated    PlanBOLStatus = "validated"
	PlanBOLStatusSubmitted    PlanBOLStatus = "submitted"
	PlanBOLStatusFulfilled    PlanBOLStatus = "fulfilled"
)

// PlanBOLRecord is the Go backend's planning-phase representation of a BOL.
// Once submitted it becomes read-only; the .NET Logistics API owns the committed BOL.
type PlanBOLRecord struct {
	ID              uuid.UUID     `json:"id"`
	DriverID        uuid.UUID     `json:"driver_id"`
	OriginatingWhID string        `json:"originating_wh_id"`
	Status          PlanBOLStatus `json:"status"`
	CreatedAt       time.Time     `json:"created_at"`
	SubmittedAt            *time.Time `json:"submitted_at,omitempty"`
	FulfilledAt            *time.Time `json:"fulfilled_at,omitempty"`
	// SubmittedTransactionID is the .NET Logistics API transaction ID assigned
	// when the plan is committed via POST /api/BillOfLading. Required for
	// ProcessStop calls on subsequent driver stop logs.
	SubmittedTransactionID *string    `json:"submitted_transaction_id,omitempty"`
}

type BOLStatusHistory struct {
	ID         uuid.UUID      `json:"id"`
	PlanBOLID  uuid.UUID      `json:"plan_bol_id"`
	FromStatus *PlanBOLStatus `json:"from_status"` // nil on initial creation
	ToStatus   PlanBOLStatus  `json:"to_status"`
	ChangedAt  time.Time      `json:"changed_at"`
}

type PlanBOLStop struct {
	ID            uuid.UUID      `json:"id"`
	PlanBOLID     uuid.UUID      `json:"plan_bol_id"`
	Sequence      int            `json:"sequence"`
	LocationID    string         `json:"location_id"`
	StopType      StopType       `json:"stop_type"`
	// DeliveryItems holds items to load (warehouse stops) or deliver (store stops).
	// Nil for return_depot stops.
	DeliveryItems map[string]int `json:"delivery_items,omitempty"`
	IsProcessed   bool           `json:"is_processed"`
	ProcessedAt   *time.Time     `json:"processed_at,omitempty"`
	DriverLogRef  *string        `json:"driver_log_ref,omitempty"`
}
