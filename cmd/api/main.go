package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"mssql-api/internal/config"
	"mssql-api/internal/database"
	"mssql-api/internal/handlers"
	"mssql-api/internal/repository"
	"mssql-api/internal/templates"
)

func main() {
	setupLogging()

	cfg, err := config.Load()
	if err != nil {
		log.Printf("Config error: %v", err)
	}

	mux := http.NewServeMux()

	if err := setupSettingsRoutes(mux); err != nil {
		log.Fatalf("Failed to setup settings routes: %v", err)
	}

	if cfg != nil {
		if err := setupAPIRoutes(mux, cfg); err != nil {
			log.Printf("API setup error: %v", err)
		}
	}

	port := "8080"
	if cfg != nil {
		port = cfg.ServerPort
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server running on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func setupLogging() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open log file:", err)
		return
	}
	log.SetOutput(logFile)
}

func setupSettingsRoutes(mux *http.ServeMux) error {
	tpl, err := templates.New()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	settingsHandler := handlers.NewSettingsHandler(tpl)

	mux.HandleFunc("/settings", settingsHandler.SettingsPage)
	mux.HandleFunc("/settings/save", settingsHandler.SaveSettings)
	mux.HandleFunc("/settings/test-connection", settingsHandler.TestConnection)

	return nil
}

func setupAPIRoutes(mux *http.ServeMux, cfg *config.Config) error {
	db, err := database.NewGormDatabase(cfg)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Note: In production, handle cleanup properly
	// defer sqlDB.Close() won't work here as it returns immediately

	repo := repository.NewRepository(db)
	handler := handlers.NewHandler(repo)

	mux.HandleFunc("/health", handler.HealthCheck)
	mux.HandleFunc("/api/barcode", handler.GetBarcodeRecords)

	log.Printf("API routes registered, DB pool stats: %+v", sqlDB.Stats())
	return nil
}