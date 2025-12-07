package models

// CategorySpending представляет траты по категории
type CategorySpending struct {
	CategoryName string
	Amount       float64
	Percentage   float64 // процент от зарплаты
}

// CategoryRecommendation представляет рекомендацию по одной категории
type CategoryRecommendation struct {
	CategoryName     string  `json:"category_name"`
	ActualAmount     float64 `json:"actual_amount"`
	ActualPercentage float64 `json:"actual_percentage"`
	RecommendedMin   float64 `json:"recommended_min"`
	RecommendedMax   float64 `json:"recommended_max"`
	Status           string  `json:"status"` // excellent, normal, warning, critical
	Message          string  `json:"message"`
	Deviation        float64 `json:"deviation"`
}

// RecommendationStatus константы для статусов
const (
	StatusExcellent = "excellent"
	StatusNormal    = "normal"
	StatusWarning   = "warning"
	StatusCritical  = "critical"
)

// OverallStatus константы для общего статуса
const (
	OverallStatusExcellent         = "excellent"
	OverallStatusGood              = "good"
	OverallStatusAttentionRequired = "attention_required"
	OverallStatusCritical          = "critical"
)

// AnalyticsResult представляет полный результат аналитики
type AnalyticsResult struct {
	UserUID         string                   `json:"user_uid"`
	Salary          float64                  `json:"salary"`
	SalaryBracket   string                   `json:"salary_bracket"`
	TotalCategories int                      `json:"total_categories"`
	ExcellentCount  int                      `json:"excellent_count"`
	NormalCount     int                      `json:"normal_count"`
	WarningCount    int                      `json:"warning_count"`
	CriticalCount   int                      `json:"critical_count"`
	OverallStatus   string                   `json:"overall_status"`
	OverallMessage  string                   `json:"overall_message"`
	Recommendations []CategoryRecommendation `json:"recommendations"`
	CalculatedAt    int64                    `json:"calculated_at"` // Unix timestamp
}
