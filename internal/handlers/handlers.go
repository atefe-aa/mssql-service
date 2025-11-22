package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"mssql-api/internal/repository"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
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

	var response interface{}
	var err error

	if strings.Contains(normalizedBarcode, "Xz") || strings.Contains(normalizedBarcode, "XZ") || strings.Contains(normalizedBarcode, "xz") {
		records, queryErr := h.repo.GetMainBarcodeRecords(r.Context(), normalizedBarcode)
		if queryErr != nil {
			http.Error(w, queryErr.Error(), http.StatusInternalServerError)
			return
		}
		response = records
	} else {
		records, queryErr := h.repo.GetChildBarcodeRecords(r.Context(), normalizedBarcode)
		if queryErr != nil {
			http.Error(w, queryErr.Error(), http.StatusInternalServerError)
			return
		}
		response = records
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}