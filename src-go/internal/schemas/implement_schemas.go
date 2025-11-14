package schemas

import "go-api/internal/models"

type ImplementCreate struct {
	Name       string `json:"name" binding:"required"`
	Brand      string `json:"brand" binding:"required"`
	VehicleModel      string `json:"model" binding:"required"`
	Year       int    `json:"year" binding:"required"`
	Identifier string `json:"identifier"`
	Type       string `json:"type"`
}

type ImplementUpdate struct {
	Name       string `json:"name"`
	Brand      string `json:"brand"`
	VehicleModel      string `json:"model"`
	Year       int    `json:"year"`
	Identifier string `json:"identifier"`
	Type       string `json:"type"`
	Status     string `json:"status"`
}

type ImplementPublic struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Brand      string `json:"brand"`
	VehicleModel      string `json:"model"`
	Year       int    `json:"year"`
	Identifier string `json:"identifier"`
	Type       string `json:"type"`
	Status     string `json:"status"`
	OrganizationID uint `json:"organization_id"`
}

func ToImplementPublic(implement models.Implement) ImplementPublic {
	return ImplementPublic{
		ID:            implement.ID,
		Name:          implement.Name,
		Brand:         implement.Brand,
		VehicleModel:         implement.VehicleModel,
		Year:          implement.Year,
		Identifier:    implement.Identifier,
		Type:          implement.Type,
		Status:        string(implement.Status),
		OrganizationID: implement.OrganizationID,
	}
}
