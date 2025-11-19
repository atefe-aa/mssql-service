package repository

import (
	"context"
	"mssql-api/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetMainBarcodeRecords fetches rows from View_Barcode_MainTube
func (r *Repository) GetMainBarcodeRecords(ctx context.Context, barcode string) ([]models.ParsicMainBarcode, error) {
	var records []models.ParsicMainBarcode
	if err := r.db.WithContext(ctx).
		Where("Str_BarcodeAdmitNum = ?", barcode).
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// GetChildBarcodeRecords fetches rows from View_Barcode_Devided
func (r *Repository) GetChildBarcodeRecords(ctx context.Context, barcode string) ([]models.ParsicChildBarcode, error) {
	var records []models.ParsicChildBarcode
	if err := r.db.WithContext(ctx).
		Where("Str_AdmitBarcodeNumber = ?", barcode).
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}
