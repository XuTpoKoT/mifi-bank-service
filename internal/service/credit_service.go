package service

import (
	"bank-service/internal/models"
	"bank-service/internal/repository"
	"database/sql"
)

type CreditService struct {
	db           *sql.DB
	creditRepo   *repository.CreditRepository
	accountRepo  *repository.AccountRepository
	scheduleRepo *repository.ScheduleRepository
}

func NewCreditService(
	db *sql.DB,
	creditRepo *repository.CreditRepository,
	accountRepo *repository.AccountRepository,
	scheduleRepo *repository.ScheduleRepository,
) *CreditService {
	return &CreditService{
		db:           db,
		creditRepo:   creditRepo,
		accountRepo:  accountRepo,
		scheduleRepo: scheduleRepo,
	}
}

func (s *CreditService) Create(
	userID int64,
	accountID int64,
	principal float64,
	months int,
) error {

	account, err := s.accountRepo.FindByID(accountID)
	if err != nil {
		return err
	}

	if account.UserID != userID {
		return ErrAccessDenied
	}

	rate := 20.0 // позже заменим на ЦБ + margin

	monthly := CalculateAnnuity(
		principal,
		rate,
		months,
	)

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	credit := &models.Credit{
		UserID:         userID,
		AccountID:      accountID,
		Principal:      principal,
		AnnualRate:     rate,
		TermMonths:     months,
		MonthlyPayment: monthly,
		RemainingDebt:  principal,
	}

	err = s.creditRepo.Create(tx, credit)
	if err != nil {
		return err
	}

	err = s.scheduleRepo.CreateSchedule(
		tx,
		credit.ID,
		monthly,
		months,
	)
	if err != nil {
		return err
	}

	err = s.accountRepo.TopUpTx(
		tx,
		accountID,
		principal,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *CreditService) GetSchedule(
	userID int64,
	creditID int64,
) ([]models.PaymentSchedule, error) {

	credit, err := s.creditRepo.
		FindByID(creditID)
	if err != nil {
		return nil, err
	}

	if credit.UserID != userID {
		return nil, ErrAccessDenied
	}

	return s.scheduleRepo.
		FindByCreditID(creditID)
}

func (s *CreditService) ProcessDuePayments() {
	payments, err := s.scheduleRepo.FindDue()
	if err != nil {
		return
	}

	for _, p := range payments {

		credit, err := s.creditRepo.
			FindByID(p.CreditID)
		if err != nil {
			continue
		}

		account, err := s.accountRepo.
			FindByID(credit.AccountID)
		if err != nil {
			continue
		}

		if account.Balance >= p.Amount {

			err = s.accountRepo.Withdraw(
				account.ID,
				p.Amount,
			)
			if err != nil {
				continue
			}

			_ = s.scheduleRepo.
				MarkPaid(p.ID)

		} else {

			penalty :=
				p.Amount * 0.10

			_ = s.scheduleRepo.
				MarkOverdue(
					p.ID,
					penalty,
				)
		}
	}
}
