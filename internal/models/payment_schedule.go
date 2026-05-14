package models

import "time"

type PaymentSchedule struct {
	ID       int64      `json:"id"`
	CreditID int64      `json:"credit_id"`
	DueDate  time.Time  `json:"due_date"`
	Amount   float64    `json:"amount"`
	Status   string     `json:"status"`
	Penalty  float64    `json:"penalty"`
	PaidAt   *time.Time `json:"paid_at,omitempty"`
}
