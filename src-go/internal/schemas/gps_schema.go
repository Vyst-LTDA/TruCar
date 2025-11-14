package schemas

// LocationCreate é o schema para o payload do ping de GPS.
type LocationCreate struct {
	VehicleID uint    `json:"vehicle_id" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}
