package models

import "time"

type Warehouse struct {
	ID        string     `json:"id"`
	Region    *string    `json:"region,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
