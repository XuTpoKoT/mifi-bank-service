package repository

import (
	"bank-service/internal/models"
	"database/sql"
	"time"
)

type ScheduleRepository struct {
	db *sql.DB
}

func NewScheduleRepository(
	db *sql.DB,
) *ScheduleRepository {
	return &ScheduleRepository{
		db: db,
	}
}

func (r *ScheduleRepository) CreateSchedule(
	tx *sql.Tx,
	creditID int64,
	amount float64,
	months int,
) error {

	for i := 1; i <= months; i++ {
		dueDate := time.Now().
			AddDate(0, i, 0)

		_, err := tx.Exec(`
			INSERT INTO payment_schedules (
				credit_id,
				due_date,
				amount,
				status,
				penalty
			)
			VALUES ($1,$2,$3,'PENDING',0)
		`,
			creditID,
			dueDate,
			amount,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ScheduleRepository) FindByCreditID(
	creditID int64,
) ([]models.PaymentSchedule, error) {

	rows, err := r.db.Query(`
		SELECT
			id,
			credit_id,
			due_date,
			amount,
			status,
			penalty,
			paid_at
		FROM payment_schedules
		WHERE credit_id = $1
		ORDER BY due_date
	`, creditID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.PaymentSchedule

	for rows.Next() {
		var p models.PaymentSchedule
		var paidAt sql.NullTime

		err := rows.Scan(
			&p.ID,
			&p.CreditID,
			&p.DueDate,
			&p.Amount,
			&p.Status,
			&p.Penalty,
			&paidAt,
		)
		if err != nil {
			return nil, err
		}

		if paidAt.Valid {
			p.PaidAt = &paidAt.Time
		}

		result = append(result, p)
	}

	return result, rows.Err()
}

func (r *ScheduleRepository) FindDue() (
	[]models.PaymentSchedule,
	error,
) {

	rows, err := r.db.Query(`
		SELECT
			id,
			credit_id,
			due_date,
			amount,
			status,
			penalty
		FROM payment_schedules
		WHERE due_date <= CURRENT_DATE
		  AND status = 'PENDING'
		ORDER BY due_date
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.PaymentSchedule

	for rows.Next() {
		var p models.PaymentSchedule

		err := rows.Scan(
			&p.ID,
			&p.CreditID,
			&p.DueDate,
			&p.Amount,
			&p.Status,
			&p.Penalty,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, p)
	}

	return result, rows.Err()
}

func (r *ScheduleRepository) MarkPaid(
	id int64,
) error {

	_, err := r.db.Exec(`
		UPDATE payment_schedules
		SET
			status = 'PAID',
			paid_at = now()
		WHERE id = $1
	`, id)

	return err
}

func (r *ScheduleRepository) MarkOverdue(
	id int64,
	penalty float64,
) error {

	_, err := r.db.Exec(`
		UPDATE payment_schedules
		SET
			status = 'OVERDUE',
			penalty = $2
		WHERE id = $1
	`, id, penalty)

	return err
}
