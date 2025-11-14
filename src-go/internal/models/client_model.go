package models

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	Name            string `gorm:"size:255;not null"`
	ContactPerson   string `gorm:"size:255"`
	Email           string `gorm:"size:255"`
	Phone           string `gorm:"size:50"`
	Address         string `gorm:"size:255"`
	City            string `gorm:"size:100"`
	State           string `gorm:"size:100"`
	ZipCode         string `gorm:"size:20"`
	Country         string `gorm:"size:100"`
	OrganizationID  uint   `gorm:"not null"`
	Organization    Organization
}
