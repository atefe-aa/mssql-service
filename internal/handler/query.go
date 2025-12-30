package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func buildSingleBarcodeQuery(barcode string) string {
	decoded := decodeBarcode(barcode)

	if strings.Contains(strings.ToLower(barcode), "xz") || strings.Contains(strings.ToLower(barcode), ".99-") {
		return fmt.Sprintf(`SELECT * FROM View_Barcode_MainTube WHERE Str_BarcodeAdmitNum = '%s'`, decoded)
	}
	return fmt.Sprintf(`SELECT * FROM View_Barcode_Devided WHERE Str_AdmitBarcodeNumber = '%s'`, decoded)
}

func (h *Handler) executeQuery(query string) (interface{}, error) {
	client := http.Client{Timeout: 30 * time.Second}
	payload := map[string]string{"query": query}
	payloadBytes, _ := json.Marshal(payload)

	resp, err := client.Post(h.cfg.GatewayUrl+"/api/query", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to reach gateway: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gateway returned status %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode gateway response: %w", err)
	}

	return result, nil
}

func (h *Handler) processBarcodes(client *http.Client, barcodes []string, table string, allRows *[]map[string]interface{}, totalCount *int) error {
	if len(barcodes) == 0 {
		return nil
	}

	// Build WHERE clause
	conditions := make([]string, 0, len(barcodes))
	for _, bc := range barcodes {
		conditions = append(conditions, fmt.Sprintf("'%s'", decodeBarcode(bc)))
	}

	query := buildBatchQuery(table, strings.Join(conditions, ","))
	payloadBytes, _ := json.Marshal(map[string]string{"query": query})

	// ACQUIRE SLOT (BLOCKS IF FULL)
	h.querySem <- struct{}{}
	defer func() { <-h.querySem }()

	resp, err := client.Post(h.cfg.GatewayUrl+"/api/query", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to query gateway: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("gateway error %d: %s", resp.StatusCode, string(body))
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

func buildBatchQuery(table, conditions string) string {
	if table == "View_Barcode_MainTube" {
		return fmt.Sprintf(`
			SELECT PatientName, Str_SpecialName, Str_BarcodeAdmitNum
			FROM View_Barcode_MainTube
			WHERE Str_BarcodeAdmitNum IN (%s)
		`, conditions)
	}
	return fmt.Sprintf(`
		SELECT PatientName, Str_SpecialName, Str_AdmitBarcodeNumber
		FROM View_Barcode_Devided
		WHERE Str_AdmitBarcodeNumber IN (%s)
	`, conditions)
}