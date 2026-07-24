package models

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID        uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TeamName  string       `gorm:"type:varchar(20);not null;uniqueIndex" json:"team_name"`
	CreatedAt time.Time    `gorm:"type:timestamp(0);default:CURRENT_TIMESTAMP" json:"created_at"`
	Members   []TeamMember `gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE" json:"members,omitempty"`
}
