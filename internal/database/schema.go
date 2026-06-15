package database

import "database/sql"

func CreateSchema(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		mfa_enabled BOOLEAN NOT NULL DEFAULT 0,
		totp_secret TEXT,
		failed_attempts INTEGER NOT NULL DEFAULT 0,
		locked_until DATETIME,
		created_at DATETIME NOT NULL,
		last_login DATETIME
	);
	`

	_, err := db.Exec(query)
	return err
}
