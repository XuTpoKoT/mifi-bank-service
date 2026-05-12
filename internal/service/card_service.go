package service

import (
	"errors"

	"bank-service/internal/models"
	"bank-service/internal/repository"
	"bank-service/internal/security"
)

type CardService struct {
	cardRepo    *repository.CardRepository
	accountRepo *repository.AccountRepository
}

func NewCardService(
	cardRepo *repository.CardRepository,
	accountRepo *repository.AccountRepository,
) *CardService {
	return &CardService{
		cardRepo:    cardRepo,
		accountRepo: accountRepo,
	}
}

func (s *CardService) Create(
	userID int64,
	accountID int64,
) (*models.Card, string, string, error) {

	account, err := s.accountRepo.FindByID(accountID)
	if err != nil {
		return nil, "", "", err
	}

	if account.UserID != userID {
		return nil, "", "", errors.New("access denied")
	}

	pan, err := security.GeneratePAN()
	if err != nil {
		return nil, "", "", err
	}
	panHMAC := security.ComputeHMAC(pan)

	expiry := security.GenerateExpiry()

	cvv, err := security.GenerateCVV()
	if err != nil {
		return nil, "", "", err
	}

	cvvHash, err := security.Hash(cvv)
	if err != nil {
		return nil, "", "", err
	}

	encPAN, err := security.Encrypt(pan)
	if err != nil {
		return nil, "", "", err
	}

	encExpiry, err := security.Encrypt(expiry)
	if err != nil {
		return nil, "", "", err
	}

	card := &models.Card{
		AccountID:       accountID,
		EncryptedPAN:    encPAN,
		EncryptedExpiry: encExpiry,
		PanHMAC:         panHMAC,
		CVVHash:         cvvHash,
	}

	err = s.cardRepo.Create(card)
	if err != nil {
		return nil, "", "", err
	}

	// PAN и CVV показываем только один раз
	return card, pan, cvv, nil
}

func (s *CardService) GetAll(
	userID int64,
) ([]models.CardView, error) {

	cards, err := s.cardRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	result := make([]models.CardView, 0, len(cards))

	for _, card := range cards {

		pan, err := security.Decrypt(
			card.EncryptedPAN,
		)
		if err != nil {
			return nil, err
		}

		expiry, err := security.Decrypt(
			card.EncryptedExpiry,
		)
		if err != nil {
			return nil, err
		}

		// проверка целостности
		if security.ComputeHMAC(pan) != card.PanHMAC {
			return nil, errors.New(
				"card integrity violation",
			)
		}

		result = append(
			result,
			models.CardView{
				ID:        card.ID,
				AccountID: card.AccountID,
				PAN:       pan,
				Expiry:    expiry,
			},
		)
	}

	return result, nil
}
