package repository

import (
	"bank-service/internal/models"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(
	u *models.User,
) error {
	return r.db.QueryRow(`
        INSERT INTO users
        (username,email,password_hash)
        VALUES ($1,$2,$3)
        RETURNING id
    `,
		u.Username,
		u.Email,
		u.PasswordHash,
	).Scan(&u.ID)
}

func (r *UserRepository) FindByEmail(
	email string,
) (*models.User, error) {
	u := &models.User{}

	err := r.db.QueryRow(`
        SELECT id, username, email, password_hash
        FROM users
        WHERE email=$1
    `, email).
		Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.PasswordHash,
		)

	return u, err
}
