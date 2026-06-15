package auth

import (
	"database/sql"
	"time"

	"cli-auth/internal/models"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) CreateUser(
	username string,
	passwordHash string,
) error {

	query := `
	INSERT INTO users (
		username,
		password_hash,
		created_at
	)
	VALUES (?, ?, ?)
	`

	_, err := r.DB.Exec(
		query,
		username,
		passwordHash,
		time.Now(),
	)

	return err
}

func (r *Repository) GetUserByUsername(
	username string,
) (*models.User, error) {

	query := `
	SELECT
		id,
		username,
		password_hash,
		mfa_enabled,
		totp_secret,
		failed_attempts,
		locked_until,
		created_at,
		last_login
	FROM users
	WHERE username = ?
	`

	user := &models.User{}

	err := r.DB.QueryRow(
		query,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.MFAEnabled,
		&user.TOTPSecret,
		&user.FailedAttempts,
		&user.LockedUntil,
		&user.CreatedAt,
		&user.LastLogin,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
