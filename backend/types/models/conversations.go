package models

import "github.com/google/uuid"

type Conversation struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TeamID         uuid.UUID `gorm:"type:uuid;not null"`
	TeamIncidentID uuid.UUID `gorm:"type:uuid;not null"`
	UserID         uuid.UUID `gorm:"type:uuid;not null"`
	Title          string    `gorm:"type:varchar(255)"`
}
