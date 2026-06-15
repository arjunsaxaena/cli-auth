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

func (r *Repository) UpdateLastLogin(
	userID int64,
) error {

	query := `
	UPDATE users
	SET last_login = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	_, err := r.DB.Exec(query, userID)

	return err
}

func (r *Repository) ResetFailedAttempts(
	userID int64,
) error {

	query := `
	UPDATE users
	SET failed_attempts = 0,
	    locked_until = NULL
	WHERE id = ?
	`

	_, err := r.DB.Exec(query, userID)

	return err
}

func (r *Repository) IncrementFailedAttempts(
	userID int64,
) error {

	query := `
	UPDATE users
	SET failed_attempts = failed_attempts + 1
	WHERE id = ?
	`

	_, err := r.DB.Exec(
		query,
		userID,
	)

	return err
}

func (r *Repository) LockAccount(
	userID int64,
	until time.Time,
) error {

	query := `
	UPDATE users
	SET locked_until = ?
	WHERE id = ?
	`

	_, err := r.DB.Exec(
		query,
		until,
		userID,
	)

	return err
}
