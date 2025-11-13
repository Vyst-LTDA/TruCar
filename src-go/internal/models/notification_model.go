package models

import (
	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationTypeMaintenanceDueDate      NotificationType = "maintenance_due_date"
	NotificationTypeMaintenanceDueKm        NotificationType = "maintenance_due_km"
	NotificationTypeDocumentExpiring        NotificationType = "document_expiring"
	NotificationTypeLowStock                NotificationType = "low_stock"
	NotificationTypeTireStatusBad           NotificationType = "tire_status_bad"
	NotificationTypeAbnormalFuelConsumption NotificationType = "abnormal_fuel_consumption"
	NotificationTypeCostExceeded            NotificationType = "cost_exceeded"
	NotificationTypeNewFineRegistered       NotificationType = "new_fine_registered"
	NotificationTypeFinePaymentDue          NotificationType = "fine_payment_due"
	NotificationTypeFreightAssigned         NotificationType = "freight_assigned"
	NotificationTypeFreightUpdated          NotificationType = "freight_updated"
	NotificationTypeMaintenanceRequestNew   NotificationType = "maintenance_request_new"
	NotificationTypeMaintenanceStatusUpdate NotificationType = "maintenance_request_status_update"
	NotificationTypeMaintenanceNewComment   NotificationType = "maintenance_request_new_comment"
	NotificationTypeJourneyStarted          NotificationType = "journey_started"
	NotificationTypeJourneyEnded            NotificationType = "journey_ended"
	NotificationTypeAchievementUnlocked     NotificationType = "achievement_unlocked"
	NotificationTypeLeaderboardTop3         NotificationType = "leaderboard_top3"
)

type Notification struct {
	gorm.Model
	OrganizationID    uint             `gorm:"not null"`
	UserID            uint             `gorm:"not null"`
	Message           string           `gorm:"type:text;not null"`
	IsRead            bool             `gorm:"default:false;not null"`
	NotificationType  NotificationType `gorm:"not null"`
	RelatedEntityType string           `gorm:"size:50"`
	RelatedEntityID   *uint
	RelatedVehicleID  *uint
	User              User
	Vehicle           *Vehicle
	Organization      Organization
}
