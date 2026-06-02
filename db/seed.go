package seeds

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

	log.Println("Running database migrations...")
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	log.Println("Database migration completed successfully!")

	seedUsers(db)
}

func seedUsers(db *gorm.DB) {
	users := []models.User{
		{
			ID:          uuid.New(),
			FirebaseUID: "fb_mock_admin_999",
			Email:       "superadmin@company.com",
			DisplayName: "System Admin",
			Scope:       "superadmin",
		},
	}
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoNothing: true,
	}).Create(&users).Error

	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println(("Database successfully seeded"))
}
