package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	FirebaseUID   string     `gorm:"type:varchar(255);not null;unique"`
	Email         string     `gorm:"type:varchar(255);not null;unique"`
	DisplayName   string     `gorm:"type:varchar(100)"`
	CreatedAt     time.Time  `gorm:"type:timestamp(0);default:CURRENT_TIMESTAMP"`
	DeactivatedAt *time.Time `gorm:"type:timestamp(0)"`
	Scope         string     `gorm:"type:varchar(50);not null"`
}
