package cmd

import (
	"fmt"
	"log"
	"os"

	seeds "github.com/WillieBam/support_copilot/db"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var rootCmd = &cobra.Command{
	Use:   "support-copilot",
	Short: "Support copilot",
	Long:  "Support Copilot",
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database auto-migration and seeding",
	Long:  "Run database auto-migration and seeding",
	Run: func(cmd *cobra.Command, args []string) {
		dbHost := os.Getenv("DB_HOST")
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbPort := os.Getenv("DB_PORT")

		if dbHost == "" {
			dbHost = "localhost"
		}
		if dbPort == "" {
			dbPort = "5432"
		}

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			dbHost, dbUser, dbPass, dbName, dbPort)

		gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		seeds.InitDatabase(gormDB)

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	rootCmd.Flags().BoolP("toggle", "t", false, "yeehaa")
}
