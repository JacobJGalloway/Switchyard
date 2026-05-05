package models

import (
	"time"

	"github.com/google/uuid"
)

type PairingStatus string

const (
	PairingStatusProposed  PairingStatus = "proposed"
	PairingStatusConfirmed PairingStatus = "confirmed"
	PairingStatusCancelled PairingStatus = "cancelled"
)

// PlanBOLPairing links an active BOL to its dead-head return BOL.
// EarliestValidAt enforces the 4-hour pre-arrangement hard constraint at the
// service layer — a pairing where paired_at < earliest_valid_at is rejected.
type PlanBOLPairing struct {
	ID              uuid.UUID     `json:"id"`
	ActiveBOLID     uuid.UUID     `json:"active_bol_id"`
	DeadheadBOLID   uuid.UUID     `json:"deadhead_bol_id"`
	PairedAt        time.Time     `json:"paired_at"`
	EarliestValidAt time.Time     `json:"earliest_valid_at"`
	OriginWarehouse string        `json:"origin_warehouse"`
	Status          PairingStatus `json:"status"`
}
