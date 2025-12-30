package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

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
