package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

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

	query := buildSingleBarcodeQuery(normalizedBarcode)

	result, err := h.executeQuery(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	select {
	case h.batchQueue <- batchJob{
		requestID:   requestID,
		barcodes:    payload.Barcodes,
		callbackURL: payload.CallbackURL,
	}:
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(BatchResponse{
			Success:   true,
			RequestID: requestID,
			Message:   "Batch processing started",
		})
	default:
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "system_overloaded",
		})
	}
}