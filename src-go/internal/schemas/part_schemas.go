package schemas

type PartCreate struct {
	Name            string   `json:"name" binding:"required"`
	Category        string   `json:"category" binding:"required"`
	MinimumStock    int      `json:"minimum_stock"`
	InitialQuantity int      `json:"initial_quantity"`
	PartNumber      *string  `json:"part_number"`
	Brand           *string  `json:"brand"`
	Location        *string  `json:"location"`
	Notes           *string  `json:"notes"`
	Value           *float64 `json:"value"`
	SerialNumber    *string  `json:"serial_number"`
	LifespanKM      *int     `json:"lifespan_km"`
}

type PartUpdate struct {
	Name         string   `json:"name"`
	Category     string   `json:"category"`
	MinimumStock int      `json:"minimum_stock"`
	PartNumber   *string  `json:"part_number"`
	Brand        *string  `json:"brand"`
	Location     *string  `json:"location"`
	Notes        *string  `json:"notes"`
	Value        *float64 `json:"value"`
	SerialNumber *string  `json:"serial_number"`
	LifespanKM   *int     `json:"lifespan_km"`
	Condition    *string  `json:"condition"`
}

type AddItemsPayload struct {
	Quantity int    `json:"quantity" binding:"required"`
	Notes    string `json:"notes"`
}

type SetItemStatusPayload struct {
	NewStatus        string `json:"new_status" binding:"required"`
	RelatedVehicleID *uint  `json:"related_vehicle_id"`
	Notes            *string `json:"notes"`
}
