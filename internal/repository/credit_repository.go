package repository

import (
	"bank-service/internal/models"
	"database/sql"
)

type CreditRepository struct {
	db *sql.DB
}

func NewCreditRepository(db *sql.DB) *CreditRepository {
	return &CreditRepository{db: db}
}

func (r *CreditRepository) Create(
	tx *sql.Tx,
	c *models.Credit,
) error {
	return tx.QueryRow(`
		INSERT INTO credits (
			user_id,
			account_id,
			principal,
			annual_rate,
			term_months,
			monthly_payment,
			remaining_debt
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id
	`,
		c.UserID,
		c.AccountID,
		c.Principal,
		c.AnnualRate,
		c.TermMonths,
		c.MonthlyPayment,
		c.RemainingDebt,
	).Scan(&c.ID)
}

func (r *CreditRepository) FindByID(
	id int64,
) (*models.Credit, error) {

	var c models.Credit

	err := r.db.QueryRow(`
		SELECT
			id,
			user_id,
			account_id,
			principal,
			annual_rate,
			term_months,
			monthly_payment,
			remaining_debt,
			status
		FROM credits
		WHERE id = $1
	`, id).Scan(
		&c.ID,
		&c.UserID,
		&c.AccountID,
		&c.Principal,
		&c.AnnualRate,
		&c.TermMonths,
		&c.MonthlyPayment,
		&c.RemainingDebt,
		&c.Status,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}
