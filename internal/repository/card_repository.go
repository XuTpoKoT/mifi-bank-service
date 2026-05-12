package repository

import (
	"bank-service/internal/models"
	"database/sql"
)

type CardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{db: db}
}

func (r *CardRepository) Create(card *models.Card) error {
	return r.db.QueryRow(`
		INSERT INTO cards (
			account_id,
			encrypted_pan,
			encrypted_expiry,
			pan_hmac,
			cvv_hash
		)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id
	`,
		card.AccountID,
		card.EncryptedPAN,
		card.EncryptedExpiry,
		card.PanHMAC,
		card.CVVHash,
	).Scan(&card.ID)
}

func (r *CardRepository) FindByUserID(
	userID int64,
) ([]models.Card, error) {

	rows, err := r.db.Query(`
		SELECT
			c.id,
			c.account_id,
			c.encrypted_pan,
			c.encrypted_expiry,
			c.pan_hmac
		FROM cards c
		JOIN accounts a
			ON a.id = c.account_id
		WHERE a.user_id = $1
		ORDER BY c.id
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []models.Card

	for rows.Next() {
		var c models.Card

		err := rows.Scan(
			&c.ID,
			&c.AccountID,
			&c.EncryptedPAN,
			&c.EncryptedExpiry,
			&c.PanHMAC,
		)
		if err != nil {
			return nil, err
		}

		cards = append(cards, c)
	}

	return cards, rows.Err()
}
