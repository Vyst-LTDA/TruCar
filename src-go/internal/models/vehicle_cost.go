package models

import (
	"time"

	"gorm.io/gorm"
)

// CostType define um tipo para padronizar os tipos de custo
type CostType string

const (
	CostTypeManutencao      CostType = "Manutenção"
	CostTypeCombustivel     CostType = "Combustível"
	CostTypePedagio         CostType = "Pedágio"
	CostTypeSeguro          CostType = "Seguro"
	CostTypePneu            CostType = "Pneu"
	CostTypePecasComponentes CostType = "Peças e Componentes"
	CostTypeMulta           CostType = "Multa"
	CostTypeOutros          CostType = "Outros"
)

// VehicleCost representa o modelo de dados para custos de veículos.
type VehicleCost struct {
	gorm.Model
	Description    string    `json:"description"`
	Amount         float64   `json:"amount"`
	Date           time.Time `json:"date"`
	CostType       CostType  `json:"cost_type" gorm:"type:varchar(50)"`
	VehicleID      uint      `json:"vehicle_id"`
	OrganizationID uint      `json:"organization_id"`
	FineID         *uint     `json:"fine_id" gorm:"unique;default:null"` // Ponteiro para permitir nulos

	// Associações
	Vehicle      Vehicle      `json:"vehicle,omitempty"`
	Organization Organization `json:"organization,omitempty"`
	Fine         *Fine        `json:"fine,omitempty"`
}
