package schemas

import "time"

type DocumentCreate struct {
	DocumentType string     `form:"document_type" binding:"required"`
	ExpiryDate   time.Time  `form:"expiry_date" binding:"required"`
	Notes        *string    `form:"notes"`
	VehicleID    *uint      `form:"vehicle_id"`
	DriverID     *uint      `form:"driver_id"`
}
