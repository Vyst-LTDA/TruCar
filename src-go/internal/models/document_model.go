package models

import (
	"time"
	"gorm.io/gorm"
)

type DocumentType string

const (
	DocumentTypeCNH    DocumentType = "CNH"
	DocumentTypeCRLV   DocumentType = "CRLV"
	DocumentTypeANTT   DocumentType = "ANTT"
	DocumentTypeASO    DocumentType = "ASO"
	DocumentTypeSeguro DocumentType = "Seguro"
	DocumentTypeOutro  DocumentType = "Outro"
)

type Document struct {
	gorm.Model
	DocumentType   DocumentType `gorm:"not null"`
	ExpiryDate     time.Time    `gorm:"type:date;not null;index"`
	FileURL        string       `gorm:"size:512;not null"`
	Notes          *string      `gorm:"type:text"`
	OrganizationID uint         `gorm:"not null"`
	VehicleID      *uint
	DriverID       *uint
	Organization   Organization
	Vehicle        *Vehicle
	Driver         *User
}
