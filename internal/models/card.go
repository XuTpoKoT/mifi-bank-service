package models

type Card struct {
	ID              int64 `json:"id"`
	AccountID       int64 `json:"account_id"`
	EncryptedPAN    string
	EncryptedExpiry string
	PanHMAC         string
	CVVHash         string
}
