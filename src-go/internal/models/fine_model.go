package models

import (
	"time"
	"gorm.io/gorm"
)

type FineStatus string

const (
	FineStatusPending  FineStatus = "Pendente"
	FineStatusPaid     FineStatus = "Paga"
	FineStatusAppealed FineStatus = "Em Recurso"
	FineStatusCanceled FineStatus = "Cancelada"
)

type Fine struct {
	gorm.Model
	Description    string     `gorm:"size:255;not null"`
	InfractionCode string     `gorm:"size:50"`
	Date           time.Time  `gorm:"type:date;not null"`
	Value          float64    `gorm:"not null"`
	Status         FineStatus `gorm:"size:20;not null;default:'Pendente'"`
	VehicleID      uint       `gorm:"not null"`
	DriverID       *uint
	OrganizationID uint       `gorm:"not null"`
	Vehicle        Vehicle
	Driver         *User
	Organization   Organization
	// O custo associado será tratado na camada de serviço,
	// pois a relação em GORM é mais complexa de configurar do que em SQLAlchemy.
}
