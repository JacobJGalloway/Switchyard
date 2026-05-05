package models

import (
	"time"

	"github.com/google/uuid"
)

type LineItem struct {
	SKUID        string `json:"sku_id"`
	QtyDelivered int    `json:"qty_delivered"`
	UnitRef      string `json:"unit_ref"`
}

// DeliveryConfirmation is created when a driver logs a completed store stop.
// InvoiceID is set once the corresponding InternalInvoice has been generated.
type DeliveryConfirmation struct {
	ID            uuid.UUID  `json:"id"`
	PlanBOLStopID uuid.UUID  `json:"plan_bol_stop_id"`
	DriverID      uuid.UUID  `json:"driver_id"`
	ConfirmedAt   time.Time  `json:"confirmed_at"`
	InvoiceID     *uuid.UUID `json:"invoice_id,omitempty"`
}

// InternalInvoice is generated automatically on delivery confirmation.
// Output destination is controlled by INVOICE_OUTPUT_PATH env var.
type InternalInvoice struct {
	ID          uuid.UUID  `json:"id"`
	StoreID     string     `json:"store_id"`
	PlanBOLID   uuid.UUID  `json:"plan_bol_id"`
	LineItems   []LineItem `json:"line_items"`
	OutputPath  string     `json:"output_path"`
	GeneratedAt time.Time  `json:"generated_at"`
}
