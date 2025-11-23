package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mssql-api/internal/config"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) GetBarcodeRecords(w http.ResponseWriter, r *http.Request) {
	barcode := r.URL.Query().Get("barcode")
	if barcode == "" {
		http.Error(w, "barcode parameter is required", http.StatusBadRequest)
		return
	}
	normalizedBarcode := strings.TrimSpace(barcode)

	var query string
	if strings.Contains(strings.ToLower(normalizedBarcode), "xz") {
		query = fmt.Sprintf(`SELECT * FROM View_Barcode_MainTube WHERE Str_BarcodeAdmitNum = '%s'`, decodeBarcode(normalizedBarcode))
	} else {
		query = fmt.Sprintf(`SELECT * FROM View_Barcode_Devided WHERE Str_AdmitBarcodeNumber = '%s'`, decodeBarcode(normalizedBarcode))
	}

	client := http.Client{Timeout: 10 * time.Second}
	payload := map[string]string{"query": query}
	payloadBytes, _ := json.Marshal(payload)

	resp, err := client.Post(h.cfg.GatewayUrl+"/api/query", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "failed to reach gateway: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("gateway returned status %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		http.Error(w, "failed to decode gateway response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func decodeBarcode(barcode string) string {
	parts := strings.SplitN(barcode, "z", 2)

	if len(parts) != 2 {
		return barcode
	}

	prefix := parts[0]
	suffix := parts[1]

	hexMap := map[byte]string{
		'A': "10", 'a': "10",
		'B': "11", 'b': "11",
		'C': "12", 'c': "12",
		'D': "13", 'd': "13",
		'E': "14", 'e': "14",
		'F': "15", 'f': "15",
		'X': "99", 'x': "99",
	}

	var decodedPrefix strings.Builder
	for i := 0; i < len(prefix); i++ {
		char := prefix[i]
		if replacement, exists := hexMap[char]; exists {
			decodedPrefix.WriteString(replacement)
		} else {
			decodedPrefix.WriteByte(char)
		}
	}

	decodedPrefixStr := decodedPrefix.String()
	var formattedPrefix string

	if len(decodedPrefixStr) > 1 {
		firstPart := decodedPrefixStr[0:1]
		secondPart := decodedPrefixStr[1:]
		formattedPrefix = firstPart + "." + secondPart
	} else {
		formattedPrefix = decodedPrefixStr
	}

	return formattedPrefix + "-" + suffix
}