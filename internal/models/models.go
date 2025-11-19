package models

type ParsicMainBarcode struct {
	StrBarcodeAdmitNum          string  `gorm:"column:Str_BarcodeAdmitNum"`
	PrkTestInternalNum          int     `gorm:"column:Prk_TestInternalNum"`
	FrkTubeNumber               int     `gorm:"column:Frk_TubeNumber"`
	FrkDepartment               int     `gorm:"column:Frk_Department"`
	PatientName                 string  `gorm:"column:PatientName"`
	BitVirtualMotherTube        bool    `gorm:"column:Bit_VirtualMotherTube"`
	FrkTubeID                   int     `gorm:"column:Frk_TubeID"`
	IntTestWorkInProgress       string  `gorm:"column:Int_TestWorkInProgress"`
	PrkAVTRelation              int     `gorm:"column:Prk_AvTRelation"`
	StrSampleName1              string  `gorm:"column:Str_SampleName1"`
	StrTestNameGroup            string  `gorm:"column:Str_TestNameGroup"`
	StrDepartmentName1          string  `gorm:"column:Str_DepartmentName1"`
	StrTubeName                 string  `gorm:"column:Str_TubeName"`
	StrStationName              string  `gorm:"column:Str_StationName"`
	StrSpecialName              string  `gorm:"column:Str_SpecialName"`
	FrkParentTestNum            int     `gorm:"column:Frk_ParentTestNum"`
	StrParentNameTests           string  `gorm:"column:Str_ParentNameTests"`
	IntParentNumberTests        int     `gorm:"column:Int_ParentNumberTests"`
	IntPatientWorkInProgressState int   `gorm:"column:Int_PatientWorkInProgressState"`
	IntIndependantPrice         bool    `gorm:"column:Int_IndependantPrice"`
	FrkContractors              int     `gorm:"column:Frk_Contractors"`
	IntIsContractors            bool    `gorm:"column:Int_IsContractors"`
	StrOutSourceName            string  `gorm:"column:Str_OutSourceName"`
	StrBarcodeNum               string  `gorm:"column:Str_BarcodeNum"`
	PrkAdmitPatient             int     `gorm:"column:PRK_AdmitPatient"`
	IntSampleVGroupVSiteID      string  `gorm:"column:Int_SampleVGroupVSiteID"`
	FrkMultiSite                int     `gorm:"column:Frk_MultiSite"`
	FrkSampleID                 int     `gorm:"column:Frk_SampleID"`
	IntGroupID                  int     `gorm:"column:Int_GroupID"`
	StrBarcode                  string  `gorm:"column:Str_Barcode"`
	StrDepartmentName           string  `gorm:"column:Str_DepartmentName"`
	FrkParentTubeID             int     `gorm:"column:Frk_ParentTubeID"`
	FltHSampleDeadValue         float64 `gorm:"column:Flt_HSampleDeadValue"`
	IntSampleSize               float64 `gorm:"column:Int_SampleSize"`
	StrContainerName            string  `gorm:"column:Str_ContainerName"`
	IntCapacity                 float64 `gorm:"column:Int_Capacity"`
	StrSampleName               string  `gorm:"column:Str_SampleName"`
}

type ParsicChildBarcode struct {
	PatientName            string `gorm:"column:PatientName"`
	StrAdmitBarcodeNumber  string `gorm:"column:Str_AdmitBarcodeNumber"`
	FrkTestInternalNum     int    `gorm:"column:Frk_TestInternalNum"`
	StrTestNameGroup       string `gorm:"column:Str_TestNameGroup"`
	PrkTestGroupID         int    `gorm:"column:Prk_TestGroupID"`
	StrResultDate          string `gorm:"column:Str_ResultDate"`
	IntIsDateEffectivly    bool   `gorm:"column:Int_IsDateEffectivly"`
	FrkDepartment          int    `gorm:"column:Frk_Department"`
	StrTubeName            string `gorm:"column:Str_TubeName"`
	IntPatientWorkInProgressState int `gorm:"column:Int_PatientWorkInProgressState"`
	IntTestWorkInProgress  string `gorm:"column:Int_TestWorkInProgress"`
	FrkSampleID            int    `gorm:"column:Frk_SampleID"`
	PrkAVTRelation         int    `gorm:"column:Prk_AvTRelation"`
	FrkMultiSite           int    `gorm:"column:Frk_MultiSite"`
	FrkMultiRequestFor     int    `gorm:"column:Frk_MultiRequestFor"`
	StrSpecialName         string `gorm:"column:Str_SpecialName"`
	StrStationName         string `gorm:"column:Str_StationName"`
	StrSampleName          string `gorm:"column:Str_SampleName"`
	IntTestsOrderNumber    int    `gorm:"column:Int_TestsOrderNumber"`
	IntIndependantPrice    bool   `gorm:"column:Int_IndependantPrice"`
	PrkAdmitPatient        int    `gorm:"column:PRK_AdmitPatient"`
	PrkTubeID              int    `gorm:"column:Prk_TubeID"`
	FrkPatientInfo         int    `gorm:"column:Frk_PatientInfo"`
	FrkAdmitPatientNum     int    `gorm:"column:Frk_AdmitPatientNum"`
	StrBarcode             string `gorm:"column:Str_Barcode"`
	FrkContractors         int    `gorm:"column:Frk_Contractors"`
	IntIsContractors       bool   `gorm:"column:Int_IsContractors"`
	StrDepartmentName1     string `gorm:"column:Str_DepartmentName1"`
	StrOutSourceName       string `gorm:"column:Str_OutSourceName"`
	IntTestResult          string `gorm:"column:int_TestResult"`
	IntEmeregencyType      int    `gorm:"column:Int_EmeregencyType"`
	StrTimeOfEmeregency    string `gorm:"column:Str_TimeOfEmeregency"`
	StrDateOfEmeregency    string `gorm:"column:Str_DateOfEmeregency"`
	StrGeneralName         string `gorm:"column:Str_GeneralName"`
	StrDepartmentName      string `gorm:"column:Str_DepartmentName"`
}

func (ParsicMainBarcode) TableName() string {
	return "View_Barcode_MainTube"
}

func (ParsicChildBarcode) TableName() string {
	return "View_Barcode_Devided"
}
