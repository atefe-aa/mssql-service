package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"mssql-edge/internal/config"
	"mssql-edge/internal/templates"
)

type SettingsHandler struct {
	templates *templates.Templates
}

func NewSettingsHandler(tpl *templates.Templates) *SettingsHandler {
	return &SettingsHandler{templates: tpl}
}

type SettingsPageData struct {
	Config            *config.Config
	HasExistingConfig bool
	ConfigError       string
	ConnectionStatus  string
	ConnectionError   string
	ConnectionSuccess bool
	APIKey            string
}

func (h *SettingsHandler) SettingsPage(w http.ResponseWriter, r *http.Request) {

	cfg, err := config.Load()
	var currentConfig *config.Config
	if err == nil {
		currentConfig = cfg
	} else {
		currentConfig = &config.Config{
			GatewayUrl: "",
			ServerPort: "8080",
			APIKey:     "",
		}
	}

	data := SettingsPageData{
		Config:            currentConfig,
		HasExistingConfig: err == nil && currentConfig.GatewayUrl != "",
		APIKey:            currentConfig.APIKey,
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
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		if err := r.ParseForm(); err != nil {
			sendJSON(w, TestConnectionResponse{Success: false, Error: "Failed to parse form: " + err.Error()})
			return
		}
	}

	// Build config from form values (not from saved config)
	testConfig := &config.Config{
		GatewayUrl: r.FormValue("GATEWAY_URL"),
	}

	if testConfig.GatewayUrl == "" {
		sendJSON(w, TestConnectionResponse{Success: false, Error: "GATEWAY_URL cannot be empty"})
		return
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("%s/health", testConfig.GatewayUrl))
	if err != nil {
		sendJSON(w, TestConnectionResponse{Success: false, Error: "Failed to reach gateway: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		sendJSON(w, TestConnectionResponse{Success: false, Error: fmt.Sprintf("Unexpected status code: %d", resp.StatusCode)})
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		sendJSON(w, TestConnectionResponse{Success: false, Error: "Failed to parse response: " + err.Error()})
		return
	}

	if body.Status != "ok" {
		sendJSON(w, TestConnectionResponse{Success: false, Error: "Unexpected response body"})
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

	cfg, err := config.Load()
	apiKey := ""
	if err == nil {
		apiKey = cfg.APIKey
	}

	f, err := os.Create(".env")
	if err != nil {
		http.Error(w, "Failed to save .env", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	content := fmt.Sprintf(`GATEWAY_URL=%s
					SERVER_PORT=%s
					API_KEY=%s
					`,
		r.FormValue("GATEWAY_URL"),
		r.FormValue("SERVER_PORT"),
		apiKey,
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
