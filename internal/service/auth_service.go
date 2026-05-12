package service

import (
	"bank-service/internal/models"
	"bank-service/internal/repository"
	"bank-service/internal/security"
	"errors"
	"regexp"
	"strings"
)

type AuthService struct {
	users *repository.UserRepository
}

func NewAuthService(
	r *repository.UserRepository,
) *AuthService {
	return &AuthService{r}
}

var emailRegex = regexp.MustCompile(
	`^[^\s@]+@[^\s@]+\.[^\s@]+$`,
)

func (s *AuthService) Register(
	username, email, password string,
) error {
	if len(username) < 3 {
		return errors.New(
			"username too short",
		)
	}

	if !emailRegex.MatchString(email) {
		return errors.New(
			"invalid email",
		)
	}

	if len(password) < 6 {
		return errors.New(
			"password too short",
		)
	}

	hash, err := security.Hash(password)
	if err != nil {
		return err
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
	}

	err = s.users.Create(user)
	if err != nil {
		if strings.Contains(
			err.Error(),
			"users_email_key",
		) {
			return errors.New(
				"email already exists",
			)
		}
		return err
	}
	return nil
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
