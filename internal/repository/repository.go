// internal/repository/repository.go
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"mssql-api/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// GetMainBarcodeRecords queries View_Barcode_MainTube for barcodes containing 'Xz'
func (r *Repository) GetMainBarcodeRecords(ctx context.Context, barcode string) ([]models.ParsicMainBarcode, error) {
	query := `SELECT * FROM View_Barcode_MainTube WHERE Str_BarcodeAdmitNum = @p1`
	
	rows, err := r.db.QueryContext(ctx, query, barcode)
	if err != nil {
		return nil, fmt.Errorf("failed to query main barcode: %w", err)
	}
	defer rows.Close()

	var records []models.ParsicMainBarcode
	for rows.Next() {
		var record models.ParsicMainBarcode
		// Scan all columns - adjust based on your actual view columns
		if err := rows.Scan(&record.StrBarcodeAdmitNum /* add other fields here */); err != nil {
			return nil, fmt.Errorf("failed to scan main barcode: %w", err)
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return records, nil
}

// GetChildBarcodeRecords queries View_Barcode_Devided for regular barcodes
func (r *Repository) GetChildBarcodeRecords(ctx context.Context, barcode string) ([]models.ParsicChildBarcode, error) {
	query := `SELECT * FROM View_Barcode_Devided WHERE Str_AdmitBarcodeNumber = @p1`
	
	rows, err := r.db.QueryContext(ctx, query, barcode)
	if err != nil {
		return nil, fmt.Errorf("failed to query child barcode: %w", err)
	}
	defer rows.Close()

	var records []models.ParsicChildBarcode
	for rows.Next() {
		var record models.ParsicChildBarcode
		// Scan all columns - adjust based on your actual view columns
		if err := rows.Scan(&record.StrAdmitBarcodeNumber /* add other fields here */); err != nil {
			return nil, fmt.Errorf("failed to scan child barcode: %w", err)
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return records, nil
}

// Add more repository methods as needed