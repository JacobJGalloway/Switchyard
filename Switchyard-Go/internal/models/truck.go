package models

import (
	"time"

	"github.com/google/uuid"
)

// TruckInventorySnapshot records the quantity of a single SKU loaded and
// remaining at a specific stop in a planned run. One snapshot per SKU per stop.
type TruckInventorySnapshot struct {
	ID                uuid.UUID `json:"id"`
	PlanBOLID         uuid.UUID `json:"plan_bol_id"`
	PlanBOLStopID     uuid.UUID `json:"plan_bol_stop_id"`
	SKUID             string    `json:"sku_id"`
	QuantityLoaded    int       `json:"quantity_loaded"`
	QuantityRemaining int       `json:"quantity_remaining"`
	SnapshotAt        time.Time `json:"snapshot_at"`
}
