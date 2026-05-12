package models

type CardView struct {
	ID        int64  `json:"id"`
	AccountID int64  `json:"account_id"`
	PAN       string `json:"pan"`
	Expiry    string `json:"expiry"`
}
