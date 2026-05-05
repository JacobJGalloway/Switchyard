package models

import (
	"time"

	"github.com/google/uuid"
)

type Driver struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	HomeWarehouseID string    `json:"home_warehouse_id"`
	Auth0UserID     string    `json:"auth0_user_id"`
	LicenseState    string    `json:"license_state"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}
