package handler

import (
	"log"
	"mssql-edge/internal/config"
)

type Handler struct {
	cfg        *config.Config
	batchQueue chan batchJob
	querySem   chan struct{}
}

func NewHandler(cfg *config.Config) *Handler {
	h := &Handler{
		cfg:        cfg,
		batchQueue: make(chan batchJob, 50), // queue size
		querySem:   make(chan struct{}, 5),  // max 5 concurrent gateway queries
	}

	// start workers
	for i := 0; i < 3; i++ {
		go h.batchWorker(i)
	}

	return h
}

func (h *Handler) batchWorker(id int) {
	for job := range h.batchQueue {
		log.Printf("[worker-%d] Processing %s", id, job.requestID)
		h.processBatchAsync(job.requestID, job.barcodes, job.callbackURL)
	}
}