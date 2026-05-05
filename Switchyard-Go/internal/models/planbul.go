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
	PlanBOLStatusDraft     PlanBOLStatus = "draft"
	PlanBOLStatusValidated PlanBOLStatus = "validated"
	PlanBOLStatusSubmitted PlanBOLStatus = "submitted"
	PlanBOLStatusFulfilled PlanBOLStatus = "fulfilled"
)

// PlanBOLRecord is the Go backend's planning-phase representation of a BOL.
// Once submitted it becomes read-only; the .NET Logistics API owns the committed BOL.
type PlanBOLRecord struct {
	ID              uuid.UUID     `json:"id"`
	DriverID        uuid.UUID     `json:"driver_id"`
	OriginatingWhID string        `json:"originating_wh_id"`
	Status          PlanBOLStatus `json:"status"`
	CreatedAt       time.Time     `json:"created_at"`
	SubmittedAt     *time.Time    `json:"submitted_at,omitempty"`
	FulfilledAt     *time.Time    `json:"fulfilled_at,omitempty"`
}

type PlanBOLStop struct {
	ID           uuid.UUID  `json:"id"`
	PlanBOLID    uuid.UUID  `json:"plan_bol_id"`
	Sequence     int        `json:"sequence"`
	LocationID   string     `json:"location_id"`
	StopType     StopType   `json:"stop_type"`
	IsProcessed  bool       `json:"is_processed"`
	ProcessedAt  *time.Time `json:"processed_at,omitempty"`
	DriverLogRef *string    `json:"driver_log_ref,omitempty"`
}
