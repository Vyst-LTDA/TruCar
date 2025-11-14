package models

import (
	"time"
	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionTypeEntrada       TransactionType = "Entrada"
	TransactionTypeSaidaUso      TransactionType = "Instalação (Uso)"
	TransactionTypeFimDeVida     TransactionType = "Fim de Vida"
	TransactionTypeAjusteInicial TransactionType = "Ajuste Inicial"
	TransactionTypeInstalacao    TransactionType = "Instalação"
	TransactionTypeDescarte      TransactionType = "Descarte"
)

type InventoryTransaction struct {
	gorm.Model
	ItemID             uint            `gorm:"not null"`
	PartID             *uint
	UserID             *uint
	TransactionType    TransactionType `gorm:"not null"`
	Notes              *string         `gorm:"type:text"`
	RelatedVehicleID   *uint
	RelatedUserID      *uint
	Timestamp          time.Time       `gorm:"not null;default:CURRENT_TIMESTAMP"`
	Item               InventoryItem
	Part               *Part
	User               *User
	RelatedVehicle     *Vehicle
	RelatedUser        *User
}
