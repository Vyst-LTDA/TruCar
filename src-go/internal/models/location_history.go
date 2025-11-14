package models

import (
	"time"
)

// LocationHistory armazena o histórico de localizações de um veículo.
type LocationHistory struct {
	ID        uint      `gorm:"primaryKey"`
	VehicleID uint      `gorm:"not null;index"`
	Latitude  float64   `gorm:"not null"`
	Longitude float64   `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`

	Vehicle Vehicle `gorm:"foreignKey:VehicleID"`
}
