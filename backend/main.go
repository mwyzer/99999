package main

import (
	"log"

	"property-hub-backend/config"
	"property-hub-backend/database"
	"property-hub-backend/handlers"
	"property-hub-backend/routes"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set config globally for handlers
	handlers.AppConfig = cfg

	// Connect to database
	database.Connect(cfg)

	// Run auto-migration via GORM (handled in Connect)

	// Seed default data (development only)
	if cfg.AppEnv == "development" {
		database.SeedAllData()
	}

	// Setup routes
	router := routes.Setup(cfg)

	// Start server
	addr := ":" + cfg.AppPort
	log.Printf("🚀 %s server starting on %s", cfg.AppName, addr)
	log.Printf("📋 Environment: %s", cfg.AppEnv)
	log.Printf("🌐 API URL: %s/api/v1", cfg.AppURL)

	if err := router.Run(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
