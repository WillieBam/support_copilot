package models

import (
	"time"

	"github.com/google/uuid"
)

type TeamIncident struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	IncidentID uuid.UUID `gorm:"type:uuid;not null"`
	TeamID     uuid.UUID `gorm:"type:uuid;not null"`
	AssignedBy uuid.UUID `gorm:"type:uuid;not null"`
	Title      string    `gorm:"type:varchar(255);not null"`
	Status     string    `gorm:"type:varchar(20)"`
	Details    string    `gorm:"type:text"`
	AssignedAt time.Time `gorm:"type:timestamp(3);default:CURRENT_TIMESTAMP"`
}

