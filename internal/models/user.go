package models

import "time"

type User struct {
	ID             int64
	Username       string
	PasswordHash   string
	MFAEnabled     bool
	TOTPSecret     string
	FailedAttempts int
	LockedUntil    *time.Time
	CreatedAt      time.Time
	LastLogin      *time.Time
}