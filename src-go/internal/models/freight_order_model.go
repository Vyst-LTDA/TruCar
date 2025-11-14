package models

import (
	"time"
	"gorm.io/gorm"
)

type FreightStatus string

const (
	FreightStatusOpen      FreightStatus = "Aberta"
	FreightStatusClaimed   FreightStatus = "Atribuída"
	FreightStatusInTransit FreightStatus = "Em Trânsito"
	FreightStatusDelivered FreightStatus = "Entregue"
	FreightStatusCanceled  FreightStatus = "Cancelado"
)

type StopPointType string

const (
	StopPointTypePickup  StopPointType = "Coleta"
	StopPointTypeDelivery StopPointType = "Entrega"
)

type StopPointStatus string

const (
	StopPointStatusPending   StopPointStatus = "Pendente"
	StopPointStatusCompleted StopPointStatus = "Concluído"
)

type FreightOrder struct {
	gorm.Model
	Description        *string        `gorm:"size:500"`
	Status             FreightStatus  `gorm:"not null;default:'Aberta'"`
	ScheduledStartTime *time.Time
	ScheduledEndTime   *time.Time
	ClientID           uint `gorm:"not null"`
	VehicleID          *uint
	DriverID           *uint
	OrganizationID     uint `gorm:"not null"`
	Client             Client
	Vehicle            *Vehicle
	Driver             *User
	Organization       Organization
	StopPoints         []StopPoint `gorm:"foreignKey:FreightOrderID"`
	Journeys           []Journey   `gorm:"foreignKey:FreightOrderID"`
}

type StopPoint struct {
	gorm.Model
	FreightOrderID    uint            `gorm:"not null"`
	SequenceOrder     int             `gorm:"not null"`
	Type              StopPointType   `gorm:"not null"`
	Status            StopPointStatus `gorm:"not null;default:'Pendente'"`
	Address           string          `gorm:"size:500;not null"`
	CargoDescription  *string         `gorm:"size:500"`
	ScheduledTime     time.Time       `gorm:"not null"`
	ActualArrivalTime *time.Time
}
