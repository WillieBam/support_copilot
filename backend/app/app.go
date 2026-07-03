package app

import (
	"fmt"
	"log"

	"github.com/WillieBam/support_copilot/backend/app/config"
	"github.com/WillieBam/support_copilot/backend/internal/interfaces"
	firebaseRepo "github.com/WillieBam/support_copilot/backend/internal/repository/firebase"
	postgresRepo "github.com/WillieBam/support_copilot/backend/internal/repository/postgres"
	"github.com/WillieBam/support_copilot/backend/internal/service"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	Client      *appClient
	Repository  *appRepository
	Service     interfaces.IAppService
	AuthService interfaces.IAuthService
}

func NewApp() *App {
	cfg := config.Get()

	// Open DB connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)
	gormDB, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	appClient := newAppClient()
	appRepository := newAppRepository(appClient)

	// Create internal repositories using the DB connection and configuration
	userRepo := postgresRepo.NewUserRepository(gormDB)
	firebaseRepository, err := firebaseRepo.NewFirebaseRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase Repository: %v", err)
	}

	// Initialize the Authentication Service
	authService := service.New(service.AuthServiceParam{
		UserRepo:     userRepo,
		FirebaseRepo: firebaseRepository,
	})

	appService := service.NewAppService()

	return &App{
		Client:      appClient,
		Repository:  appRepository,
		Service:     appService,
		AuthService: authService,
	}
}
