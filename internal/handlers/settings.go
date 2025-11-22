package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"time"

	"mssql-api/internal/config"
	"mssql-api/internal/database"
	"mssql-api/internal/templates"
)

type SettingsHandler struct {
	templates *templates.Templates
}

func NewSettingsHandler(tpl *templates.Templates) *SettingsHandler {
	return &SettingsHandler{templates: tpl}
}

type SettingsPageData struct {
	CurrentUser       string
	Username          string
	Userdomain        string
	Config            *config.Config
	HasExistingConfig bool
	ConfigError       string
	ConnectionStatus  string
	ConnectionError   string
	ConnectionSuccess bool
}

func (h *SettingsHandler) SettingsPage(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := user.Current()
	username := os.Getenv("USERNAME")
	userdomain := os.Getenv("USERDOMAIN")

	cfg, err := config.Load()
	var currentConfig *config.Config
	if err == nil {
		currentConfig = cfg
	} else {
		currentConfig = &config.Config{
			DBServer:   "",
			DBPort:     "1433",
			DBName:     "",
			DBAuthMode: "windows",
			DBUser:     "",
			DBPassword: "",
			ServerPort: "8080",
		}
	}

	data := SettingsPageData{
		CurrentUser:       currentUser.Username,
		Username:          username,
		Userdomain:        userdomain,
		Config:            currentConfig,
		HasExistingConfig: err == nil && currentConfig.DBName != "",
	}

	if err != nil {
		data.ConfigError = err.Error()
	}

	h.templates.Render(w, "settings.html", data)
}

type TestConnectionResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func (h *SettingsHandler) TestConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ParseMultipartForm is needed for FormData sent via fetch
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		// Fallback to ParseForm for URL-encoded data
		if err := r.ParseForm(); err != nil {
			sendJSON(w, TestConnectionResponse{Success: false, Error: "Failed to parse form: " + err.Error()})
			return
		}
	}

	// Build config from form values (not from saved config)
	testConfig := &config.Config{
		DBServer:   r.FormValue("DB_SERVER"),
		DBPort:     r.FormValue("DB_PORT"),
		DBName:     r.FormValue("DB_NAME"),
		DBAuthMode: r.FormValue("DB_AUTH_MODE"),
		DBUser:     r.FormValue("DB_USER"),
		DBPassword: r.FormValue("DB_PASSWORD"),
		ServerPort: r.FormValue("SERVER_PORT"),
	}

	// Validate required fields
	if testConfig.DBServer == "" {
		sendJSON(w, TestConnectionResponse{Success: false, Error: "Server is required"})
		return
	}
	if testConfig.DBName == "" {
		sendJSON(w, TestConnectionResponse{Success: false, Error: "Database name is required"})
		return
	}

	// Test the connection
	db, err := database.NewGormDatabase(testConfig)
	if err != nil {
		sendJSON(w, TestConnectionResponse{Success: false, Error: err.Error()})
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		sendJSON(w, TestConnectionResponse{Success: false, Error: err.Error()})
		return
	}
	defer sqlDB.Close()

	// Test with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result int
	if err := db.WithContext(ctx).Raw("SELECT 1").Scan(&result).Error; err != nil {
		sendJSON(w, TestConnectionResponse{Success: false, Error: err.Error()})
		return
	}

	if result != 1 {
		sendJSON(w, TestConnectionResponse{Success: false, Error: "Unexpected query result"})
		return
	}

	sendJSON(w, TestConnectionResponse{Success: true})
}

func (h *SettingsHandler) SaveSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	f, err := os.Create(".env")
	if err != nil {
		http.Error(w, "Failed to save .env", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	content := fmt.Sprintf(`DB_SERVER=%s
DB_PORT=%s
DB_NAME=%s
DB_AUTH_MODE=%s
DB_USER=%s
DB_PASSWORD=%s
SERVER_PORT=%s
`,
		r.FormValue("DB_SERVER"),
		r.FormValue("DB_PORT"),
		r.FormValue("DB_NAME"),
		r.FormValue("DB_AUTH_MODE"),
		r.FormValue("DB_USER"),
		r.FormValue("DB_PASSWORD"),
		r.FormValue("SERVER_PORT"),
	)

	if _, err := f.WriteString(content); err != nil {
		http.Error(w, "Failed to write config", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Saved. Please restart the server manually.")
}

func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}