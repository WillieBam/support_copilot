package db

import (
	"log"

	"github.com/WillieBam/support_copilot/backend/types/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func InitDatabase(db *gorm.DB) {
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Fatalf("Failed to create UUID extension: %v", err)
	}

	err := db.AutoMigrate(&models.User{})
	if err == nil {
		err = db.AutoMigrate(&models.Alert{})
	}
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	log.Println("Database migration completed successfully!")

	seedUsers(db)
}

func seedUsers(db *gorm.DB) {
	defaultUsers := []models.User{
		{
			ID:          uuid.New(),
			FirebaseUID: "fb_superadmin_111",
			Email:       "superadmin@company.com",
			DisplayName: "System Boss",
			Scope:       "superadmin",
		},
	}

	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoNothing: true,
	}).Create(&defaultUsers).Error

	if err != nil {
		log.Printf("Warning: Seeding failed: %v", err)
	} else {
		log.Println("Database seeding done!")
	}
}
