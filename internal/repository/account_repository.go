package repository

import (
	"bank-service/internal/models"
	"database/sql"
	"errors"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db}
}

func (r *AccountRepository) Create(
	account *models.Account,
) error {
	return r.db.QueryRow(`
		INSERT INTO accounts (user_id, balance, currency)
		VALUES ($1, $2, $3)
		RETURNING id
	`,
		account.UserID,
		account.Balance,
		account.Currency,
	).Scan(&account.ID)
}

func (r *AccountRepository) FindByID(
	id int64,
) (*models.Account, error) {

	account := &models.Account{}

	err := r.db.QueryRow(`
		SELECT id, user_id, balance, currency
		FROM accounts
		WHERE id=$1
	`, id).Scan(
		&account.ID,
		&account.UserID,
		&account.Balance,
		&account.Currency,
	)

	return account, err
}

func (r *AccountRepository) TopUp(
	accountID int64,
	amount float64,
) error {

	_, err := r.db.Exec(`
		UPDATE accounts
		SET balance = balance + $1
		WHERE id=$2
	`,
		amount,
		accountID,
	)

	return err
}

func (r *AccountRepository) Transfer(
	fromID int64,
	toID int64,
	amount float64,
) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var balance float64

	// блокируем sender
	err = tx.QueryRow(`
		SELECT balance
		FROM accounts
		WHERE id=$1
		FOR UPDATE
	`, fromID).Scan(&balance)

	if err != nil {
		return err
	}

	if balance < amount {
		return errors.New("insufficient funds")
	}

	// списываем
	_, err = tx.Exec(`
		UPDATE accounts
		SET balance = balance - $1
		WHERE id=$2
	`,
		amount,
		fromID,
	)
	if err != nil {
		return err
	}

	// начисляем
	_, err = tx.Exec(`
		UPDATE accounts
		SET balance = balance + $1
		WHERE id=$2
	`,
		amount,
		toID,
	)
	if err != nil {
		return err
	}

	// логируем транзакцию
	_, err = tx.Exec(`
		INSERT INTO transactions
		(from_account_id, to_account_id, amount, type)
		VALUES ($1,$2,$3,'TRANSFER')
	`,
		fromID,
		toID,
		amount,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *AccountRepository) FindByUserID(
	userID int64,
) ([]models.Account, error) {

	rows, err := r.db.Query(`
		SELECT id, user_id, balance, currency
		FROM accounts
		WHERE user_id = $1
		ORDER BY id
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account

	for rows.Next() {
		var a models.Account

		err := rows.Scan(
			&a.ID,
			&a.UserID,
			&a.Balance,
			&a.Currency,
		)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, a)
	}

	return accounts, rows.Err()
}

func (r *AccountRepository) TopUpTx(
	tx *sql.Tx,
	accountID int64,
	amount float64,
) error {

	_, err := tx.Exec(`
		UPDATE accounts
		SET balance = balance + $2
		WHERE id = $1
	`,
		accountID,
		amount,
	)

	return err
}

func (r *AccountRepository) Withdraw(
	accountID int64,
	amount float64,
) error {

	_, err := r.db.Exec(`
		UPDATE accounts
		SET balance = balance - $2
		WHERE id = $1
	`,
		accountID,
		amount,
	)

	return err
}
