package handler

type batchJob struct {
	requestID   string
	barcodes    []string
	callbackURL string
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
	Records   []map[string]interface{} `json:"records,omitempty"`
	Count     int                      `json:"count,omitempty"`
	Error     string                   `json:"error,omitempty"`
}