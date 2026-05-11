package service

import (
	"bank-service/internal/models"
	"bank-service/internal/repository"
	"bank-service/internal/security"
	"errors"
)

type AuthService struct {
	users *repository.UserRepository
}

func NewAuthService(
	r *repository.UserRepository,
) *AuthService {
	return &AuthService{r}
}

func (s *AuthService) Register(
	username, email, password string,
) error {
	hash, err := security.Hash(password)
	if err != nil {
		return err
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
	}

	return s.users.Create(user)
}

func (s *AuthService) Login(
	email, password string,
) (string, error) {

	user, err := s.users.FindByEmail(email)
	if err != nil {
		return "", err
	}

	if !security.Verify(
		user.PasswordHash,
		password,
	) {
		return "", errors.New("invalid credentials")
	}

	return security.GenerateJWT(user.ID)
}
