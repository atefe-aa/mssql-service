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

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	showConfigGUI()

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

func showConfigGUI() {
	a := app.New()
	w := a.NewWindow("Configure Database")

	server := widget.NewEntry()
	server.SetPlaceHolder("DB_SERVER")

	dbName := widget.NewEntry()
	dbName.SetPlaceHolder("DB_NAME")

	user := widget.NewEntry()
	user.SetPlaceHolder("DB_USER")

	pass := widget.NewPasswordEntry()
	pass.SetPlaceHolder("DB_PASSWORD")

	port := widget.NewEntry()
	port.SetPlaceHolder("SERVER_PORT")
	port.SetText("8080") // default

	save := widget.NewButton("Save & Start Server", func() {
		f, err := os.Create(".env")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		fmt.Fprintf(f, "DB_SERVER=%s\nDB_NAME=%s\nDB_USER=%s\nDB_PASSWORD=%s\nSERVER_PORT=%s\n",
			server.Text, dbName.Text, user.Text, pass.Text, port.Text)

		w.Close() // close GUI and continue to start server
	})

	w.SetContent(container.NewVBox(
		server,
		dbName,
		user,
		pass,
		port,
		save,
	))

	w.ShowAndRun()
}
