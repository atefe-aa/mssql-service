package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"mssql-api/internal/config"
	"mssql-api/internal/database"
	"mssql-api/internal/handlers"
	"mssql-api/internal/repository"
)

func main() {
	// Log to file
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(logFile)
	} else {
		log.Println("Failed to open log file:", err)
	}

	// Load configuration (after .env exists)
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Config error: %v", err)
	}

	// Setup HTTP routes
	mux := http.NewServeMux()

	// CONFIG UI
	mux.HandleFunc("/settings", settingsPage)
	mux.HandleFunc("/settings/save", saveSettings)

	// API
	if cfg != nil {
		db, err := database.NewGormDatabase(cfg)
		if err == nil {
			sqlDB, _ := db.DB()
			defer sqlDB.Close()

			repo := repository.NewRepository(db)
			handler := handlers.NewHandler(repo)

			mux.HandleFunc("/health", handler.HealthCheck)
			mux.HandleFunc("/api/barcode", handler.GetBarcodeRecords)
		}
	}

	// Start server
	port := "8080"
	if cfg != nil {
		port = cfg.ServerPort
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server running on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Printf("Server failed: %v", err)
	}
}

func settingsPage(w http.ResponseWriter, r *http.Request) {
	tpl := `
<html>
<body>
<p>To use Windows Authentication:</p>
<ul>
<li>Ensure SQL Server allows Windows Authentication</li>
<li>Check that the Windows account has login permissions</li>
<li>Verify the account has access to the specific database</li>
</ul>
	<h2>Configure Database</h2>
	<form action="/settings/save" method="POST">
		DB Server: <input name="DB_SERVER"><br><br>
		DB Port: <input name="DB_PORT" value="1433"><br><br>
		DB Name: <input name="DB_NAME"><br><br>
		Authentication Mode:
		<select name="DB_AUTH_MODE">
			<option value="windows">Windows</option>
			<option value="sql">SQL</option>
		</select><br><br>
		DB User: <input name="DB_USER"><br><br>
		DB Password: <input type="password" name="DB_PASSWORD"><br><br>
		Server Port: <input name="SERVER_PORT" value="8080"><br><br>
		<button type="submit">Save & Restart</button>
	</form>
</body>
</html>
`
	t, _ := template.New("settings").Parse(tpl)
	t.Execute(w, nil)
}

func saveSettings(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	f, err := os.Create(".env")
	if err != nil {
		http.Error(w, "Failed to save .env", 500)
		return
	}
	defer f.Close()

	fmt.Fprintf(f, `DB_SERVER=%s
DB_PORT=%s
DB_NAME=%s
DB_AUTH_MODE=%s
DB_USER=%s
DB_PASSWORD=%s
SERVER_PORT=%s
`,
		r.Form.Get("DB_SERVER"),
		r.Form.Get("DB_PORT"),
		r.Form.Get("DB_NAME"),
		r.Form.Get("DB_AUTH_MODE"),
		r.Form.Get("DB_USER"),
		r.Form.Get("DB_PASSWORD"),
		r.Form.Get("SERVER_PORT"),
	)

	fmt.Fprintf(w, "Saved. Please restart the server manually.")
}
