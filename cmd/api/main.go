package main

import (
	"fmt"
	"log"
	"net/http"

	"mssql-api/internal/config"
	"mssql-api/internal/database"
	"mssql-api/internal/handlers"
	"mssql-api/internal/repository"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize GORM database
	db, err := database.NewGormDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Optional: ping database
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close() // Close underlying sql.DB when program exits

	log.Println("Database connection established successfully")

	// Initialize repository with *gorm.DB
	repo := repository.NewRepository(db)

	// Initialize handlers
	handler := handlers.NewHandler(repo)

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthCheck)
	mux.HandleFunc("/api/barcode", handler.GetBarcodeRecords)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)
	log.Printf("Database: %s (auth mode: %s)", cfg.DBName, cfg.DBAuthMode)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
