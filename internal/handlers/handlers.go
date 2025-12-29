package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mssql-edge/internal/config"
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

type BatchRequest struct {
	Barcodes    []string `json:"barcodes"`
	CallbackURL string   `json:"callback_url"`
}

type BatchResponse struct {
	Success   bool   `json:"success"`
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
}
type CallbackPayload struct {
	Success   bool                     `json:"success"`
	RequestID string                   `json:"request_id"`
	Records      []map[string]interface{} `json:"records,omitempty"`
	Count     int                      `json:"count,omitempty"`
	Error     string                   `json:"error,omitempty"`
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
	if strings.Contains(strings.ToLower(normalizedBarcode), "xz") || strings.Contains(strings.ToLower(normalizedBarcode), ".99-") {
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

func (h *Handler) GetBatchBarcodeRecords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "method_not_allowed",
		})
		return
	}

	var payload BatchRequest

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "invalid_json",
			"details": err.Error(),
		})
		return
	}

    if payload.CallbackURL == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "success": false,
            "error":   "callback_url_required",
        })
        return
    }

    requestID := fmt.Sprintf("batch_%d", time.Now().UnixNano())
    w.WriteHeader(http.StatusAccepted)
    json.NewEncoder(w).Encode(BatchResponse{
        Success:   true,
        RequestID: requestID,
        Message:   "Batch processing started",
    })
    go h.processBatchAsync(requestID, payload.Barcodes, payload.CallbackURL)
}

func (h *Handler) processBatchAsync(requestID string, barcodes []string, callbackURL string) {
	log.Printf("[%s] Starting background processing for %d barcodes", requestID, len(barcodes))

	// Separate barcodes by type
	var xzBarcodes, otherBarcodes []string
	for _, bc := range barcodes {
		bc = strings.TrimSpace(bc)
		if bc == "" {
			continue
		}
		if strings.Contains(strings.ToLower(bc), "xz") || strings.Contains(strings.ToLower(bc), ".99-") {
			xzBarcodes = append(xzBarcodes, bc)
		} else {
			otherBarcodes = append(otherBarcodes, bc)
		}
	}

	client := http.Client{Timeout: 0}
	allRows := []map[string]interface{}{}
	totalCount := 0

	// Process barcodes
	if err := h.processBarcodes(&client, xzBarcodes, "View_Barcode_MainTube", &allRows, &totalCount); err != nil {
		log.Printf("[%s] Error processing xz barcodes: %v", requestID, err)
		h.sendCallback(requestID, callbackURL, CallbackPayload{
			Success:   false,
			RequestID: requestID,
			Error:     fmt.Sprintf("Failed to process xz barcodes: %v", err),
		})
		return
	}

	if err := h.processBarcodes(&client, otherBarcodes, "View_Barcode_Devided", &allRows, &totalCount); err != nil {
		log.Printf("[%s] Error processing other barcodes: %v", requestID, err)
		h.sendCallback(requestID, callbackURL, CallbackPayload{
			Success:   false,
			RequestID: requestID,
			Error:     fmt.Sprintf("Failed to process other barcodes: %v", err),
		})
		return
	}

	// Send success callback
	log.Printf("[%s] Processing complete. Sending callback with %d rows", requestID, totalCount)
	h.sendCallback(requestID, callbackURL, CallbackPayload{
		Success:   true,
		RequestID: requestID,
		Records:      allRows,
		Count:     totalCount,
	})
}

func (h *Handler) processBarcodes(client *http.Client, barcodes []string, table string, allRows *[]map[string]interface{}, totalCount *int) error {

	if len(barcodes) == 0 {
		return nil
	}

	// Build WHERE clause
	conditions := []string{}
	for _, bc := range barcodes {
		conditions = append(conditions, fmt.Sprintf("'%s'", decodeBarcode(bc)))
	}

	var query string
	if table == "View_Barcode_MainTube" {
		query = fmt.Sprintf(`
			SELECT PatientName, Str_SpecialName, Str_BarcodeAdmitNum
			FROM View_Barcode_MainTube
			WHERE Str_BarcodeAdmitNum IN (%s)
		`, strings.Join(conditions, ","))
	} else {
		query = fmt.Sprintf(`
			SELECT PatientName, Str_SpecialName, Str_AdmitBarcodeNumber
			FROM View_Barcode_Devided
			WHERE Str_AdmitBarcodeNumber IN (%s)
		`, strings.Join(conditions, ","))
	}

	payloadBytes, _ := json.Marshal(map[string]string{"query": query})

	resp, err := client.Post(h.cfg.GatewayUrl+"/api/query", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to query gateway: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gateway returned status %d", resp.StatusCode)
	}

	var result struct {
		Success bool                     `json:"success"`
		Rows    []map[string]interface{} `json:"rows"`
		Count   int                      `json:"count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode gateway response: %w", err)
	}

	*allRows = append(*allRows, result.Rows...)
	*totalCount += result.Count
	return nil
}
func (h *Handler) sendCallback(requestID, callbackURL string, payload CallbackPayload) {
	client := http.Client{Timeout: 30 * time.Second}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[%s] Failed to marshal callback payload: %v", requestID, err)
		return
	}

	// Retry logic for callback
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := client.Post(callbackURL, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			log.Printf("[%s] Callback attempt %d failed: %v", requestID, attempt, err)
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt*2) * time.Second)
				continue
			}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.Printf("[%s] Callback successful (status: %d)", requestID, resp.StatusCode)
			return
		}

		log.Printf("[%s] Callback attempt %d returned status %d", requestID, attempt, resp.StatusCode)
		if attempt < maxRetries {
			time.Sleep(time.Duration(attempt*2) * time.Second)
		}
	}

	log.Printf("[%s] All callback attempts failed", requestID)
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
