package models

import (
	"gorm.io/gorm"
)

type ImplementStatus string

const (
	ImplementStatusAvailable   ImplementStatus = "available"
	ImplementStatusInUse       ImplementStatus = "in_use"
	ImplementStatusMaintenance ImplementStatus = "maintenance"
)

type Implement struct {
	gorm.Model
	Name           string          `gorm:"size:100;not null"`
	Brand          string          `gorm:"size:50;not null"`
	VehicleModel   string          `gorm:"size:50;not null"`
	Type           string          `gorm:"size:50"`
	Status         ImplementStatus `gorm:"size:20;not null;default:available"`
	Year           int             `gorm:"not null"`
	Identifier     string          `gorm:"size:50"`
	OrganizationID uint            `gorm:"not null"`
	Organization   Organization
}
