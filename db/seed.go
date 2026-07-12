package db

import (
	"log"
	"time"

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
	seedAlerts(db)
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

func seedAlerts(db *gorm.DB) {
	// Generating deterministic UUIDs for test incident tracking
	mockIncidentID := uuid.MustParse("9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d")

	mockAlerts := []models.Alert{
		{
			ID:          uuid.New(),
			IncidentID:  mockIncidentID,
			ReceivedAt:  time.Now().Add(-15 * time.Minute),
			ServiceName: "payment-gateway-service",
			Severity:    "CRITICAL",
			// Raw backticks allow native unescaped double quotes for cross-team json mapping
			Metrics: `{"container.cpu.usage": 94.2, "runtime.go.mem_stats.total_alloc": 4825800, "error_rate": 0.06}`,
		},
		{
			ID:          uuid.New(),
			IncidentID:  mockIncidentID,
			ReceivedAt:  time.Now().Add(-5 * time.Minute),
			ServiceName: "authentication-service",
			Severity:    "WARNING",
			Metrics:     `{"trace.grpc.server.request.hits": 4500, "system.cpu.system": 78.1}`,
		},
	}

	// We conflict check using the primary key ID since alerts are event snapshots
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&mockAlerts).Error

	if err != nil {
		log.Printf("Warning: Alert seeding failed: %v", err)
	} else {
		log.Println("Alert database seeding done!")
	}
}
