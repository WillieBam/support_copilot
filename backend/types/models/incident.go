package models

import (
	"time"

	"github.com/google/uuid"
)

type Incident struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	AlertID   uuid.UUID `gorm:"type:uuid;not null"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt time.Time `gorm:"type:timestamp(3);default:CURRENT_TIMESTAMP"`
	Details   string    `gorm:"type:text"`
	Status    string    `gorm:"type:varchar(20)"`
}