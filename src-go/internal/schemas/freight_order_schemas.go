package schemas

import "time"

type StopPointCreate struct {
	SequenceOrder    int       `json:"sequence_order" binding:"required"`
	Type             string    `json:"type" binding:"required"`
	Address          string    `json:"address" binding:"required"`
	CargoDescription *string   `json:"cargo_description"`
	ScheduledTime    time.Time `json:"scheduled_time" binding:"required"`
}

type FreightOrderCreate struct {
	Description        *string           `json:"description"`
	ScheduledStartTime *time.Time        `json:"scheduled_start_time"`
	ScheduledEndTime   *time.Time        `json:"scheduled_end_time"`
	ClientID           uint              `json:"client_id" binding:"required"`
	StopPoints         []StopPointCreate `json:"stop_points" binding:"required"`
}

type FreightOrderClaim struct {
	VehicleID uint `json:"vehicle_id" binding:"required"`
}
