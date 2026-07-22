package models

import (
	"github.com/google/uuid"
)

type TeamMember struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TeamID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_team"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_team"`
	Role   string    `gorm:"type:varchar(20);not null;default:'member'"`

	Team Team `gorm:"foreignKey:TeamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
