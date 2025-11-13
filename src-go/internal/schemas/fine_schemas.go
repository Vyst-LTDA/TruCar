package schemas

import "time"

type FineCreate struct {
	Description    string    `json:"description" binding:"required"`
	InfractionCode string    `json:"infraction_code"`
	Date           time.Time `json:"date" binding:"required"`
	Value          float64   `json:"value" binding:"required"`
	Status         string    `json:"status" binding:"required"`
	VehicleID      uint      `json:"vehicle_id" binding:"required"`
	DriverID       *uint     `json:"driver_id"`
}

type FineUpdate struct {
	Description    string    `json:"description"`
	InfractionCode string    `json:"infraction_code"`
	Date           time.Time `json:"date"`
	Value          float64   `json:"value"`
	Status         string    `json:"status"`
	VehicleID      uint      `json:"vehicle_id"`
	DriverID       *uint     `json:"driver_id"`
}
