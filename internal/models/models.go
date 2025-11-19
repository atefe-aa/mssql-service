package models

type ParsicMainBarcode struct {
	StrBarcodeAdmitNum string `json:"str_barcode_admit_num" db:"Str_BarcodeAdmitNum"`
	// Add other fields from your view here
	// Example fields (adjust based on your actual view structure):
	// PatientName        string    `json:"patient_name" db:"PatientName"`
	// AdmitDate          time.Time `json:"admit_date" db:"AdmitDate"`
	// TestType           string    `json:"test_type" db:"TestType"`
}

// ParsicChildBarcode represents data from View_Barcode_Devided
type ParsicChildBarcode struct {
	StrAdmitBarcodeNumber string `json:"str_admit_barcode_number" db:"Str_AdmitBarcodeNumber"`
	// Add other fields from your view here
	// Example fields (adjust based on your actual view structure):
	// PatientName           string    `json:"patient_name" db:"PatientName"`
	// AdmitDate             time.Time `json:"admit_date" db:"AdmitDate"`
	// TestType              string    `json:"test_type" db:"TestType"`
	// ChildBarcodeSeq       int       `json:"child_barcode_seq" db:"ChildBarcodeSeq"`
}