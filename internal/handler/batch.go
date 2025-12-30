package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func (h *Handler) processBatchAsync(requestID string, barcodes []string, callbackURL string) {
	log.Printf("[%s] Starting background processing for %d barcodes", requestID, len(barcodes))

	// Separate barcodes by type
	xzBarcodes, otherBarcodes := categorizeBarcodes(barcodes)

	client := http.Client{Timeout: 60 * time.Second}
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
		Records:   allRows,
		Count:     totalCount,
	})
}

func categorizeBarcodes(barcodes []string) (xzBarcodes, otherBarcodes []string) {
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
	return xzBarcodes, otherBarcodes
}