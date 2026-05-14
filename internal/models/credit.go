package models

type Credit struct {
	ID             int64 `json:"id"`
	UserID         int64
	AccountID      int64   `json:"account_id"`
	Principal      float64 `json:"principal"`
	AnnualRate     float64 `json:"annual_rate"`
	TermMonths     int     `json:"term_months"`
	MonthlyPayment float64 `json:"monthly_payment"`
	RemainingDebt  float64 `json:"remaining_debt"`
	Status         string  `json:"status"`
}
