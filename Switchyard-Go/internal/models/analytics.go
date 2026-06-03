package models

import "github.com/google/uuid"

type BOLStatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type BOLOperatingCost struct {
	PlanBOLID   uuid.UUID `json:"plan_bol_id"`
	DriverID    uuid.UUID `json:"driver_id"`
	WarehouseID string    `json:"warehouse_id"`
	MilesDriven float64   `json:"miles_driven"`
	BaseRate    float64   `json:"base_rate"`
	TowCost     float64   `json:"tow_cost"`
	TotalCost   float64   `json:"total_cost"`
}

type DriverOperatingCost struct {
	DriverID   uuid.UUID `json:"driver_id"`
	TotalMiles float64   `json:"total_miles"`
	TotalCost  float64   `json:"total_cost"`
	BOLCount   int       `json:"bol_count"`
}

type WarehouseOperatingCost struct {
	WarehouseID string  `json:"warehouse_id"`
	TotalMiles  float64 `json:"total_miles"`
	TotalCost   float64 `json:"total_cost"`
	BOLCount    int     `json:"bol_count"`
}
