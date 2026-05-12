package service

import (
	"bank-service/internal/logger"
	"bank-service/internal/models"
	"bank-service/internal/repository"
	"errors"
)

type AccountService struct {
	repo *repository.AccountRepository
}

func NewAccountService(
	repo *repository.AccountRepository,
) *AccountService {
	return &AccountService{repo}
}

func (s *AccountService) Create(
	userID int64,
) (*models.Account, error) {

	account := &models.Account{
		UserID:   userID,
		Balance:  0,
		Currency: "RUB",
	}
	logger.Log.Infof(
		"creating account for user=%d",
		userID,
	)
	err := s.repo.Create(account)

	return account, err
}

func (s *AccountService) TopUp(
	userID int64,
	accountID int64,
	amount float64,
) error {

	account, err := s.repo.FindByID(accountID)
	if err != nil {
		return err
	}

	if account.UserID != userID {
		return errors.New("access denied")
	}

	return s.repo.TopUp(accountID, amount)
}

func (s *AccountService) Transfer(
	userID int64,
	fromID int64,
	toID int64,
	amount float64,
) error {

	account, err := s.repo.FindByID(fromID)
	if err != nil {
		return err
	}

	if account.UserID != userID {
		return errors.New("access denied")
	}

	logger.Log.Infof(
		"transfer: user=%d from=%d to=%d amount=%.2f",
		userID,
		fromID,
		toID,
		amount,
	)
	return s.repo.Transfer(
		fromID,
		toID,
		amount,
	)
}

func (s *AccountService) GetAll(
	userID int64,
) ([]models.Account, error) {
	return s.repo.FindByUserID(userID)
}
