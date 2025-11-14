package models

import (
	"time"
	"gorm.io/gorm"
)

type PartCategory string

const (
	PartCategoryPeca       PartCategory = "Peça"
	PartCategoryFluido     PartCategory = "Fluído"
	PartCategoryConsumivel PartCategory = "Consumível"
	PartCategoryPneu       PartCategory = "Pneu"
	PartCategoryOutro      PartCategory = "Outro"
)

type Part struct {
	gorm.Model
	Name           string       `gorm:"size:255;not null;index"`
	Category       PartCategory `gorm:"not null;default:'Peça'"`
	Value          *float64
	InvoiceURL     *string `gorm:"size:512"`
	SerialNumber   *string `gorm:"size:100;index;unique"`
	MinimumStock   int     `gorm:"default:0"`
	PartNumber     *string `gorm:"size:100;index"`
	Brand          *string `gorm:"size:100"`
	Location       *string `gorm:"size:100"`
	Notes          *string `gorm:"type:text"`
	PhotoURL       *string `gorm:"size:512"`
	LifespanKM     *int
	OrganizationID uint `gorm:"not null"`
	Organization   Organization
	Items          []InventoryItem
}

type InventoryItemStatus string

const (
	InventoryItemStatusDisponivel InventoryItemStatus = "Disponível"
	InventoryItemStatusEmUso      InventoryItemStatus = "Em Uso"
	InventoryItemStatusFimDeVida  InventoryItemStatus = "Fim de Vida"
)

type InventoryItem struct {
	gorm.Model
	ItemIdentifier        int                 `gorm:"not null;index"`
	Status                InventoryItemStatus `gorm:"not null;default:'Disponível';index"`
	PartID                uint                `gorm:"not null"`
	OrganizationID        uint                `gorm:"not null"`
	InstalledOnVehicleID  *uint
	InstalledAt           *time.Time
	Part                  Part
	Organization          Organization
	InstalledOnVehicle    *Vehicle
	Transactions          []InventoryTransaction `gorm:"foreignKey:ItemID"`
}
