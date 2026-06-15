package auth

import (
	"errors"
	"fmt"
	"time"
)

const (
	MaxFailedAttempts = 5
	LockDuration      = 15 * time.Minute
)

func Login(
	repo *Repository,
	username string,
	password string,
) error {

	user, err := repo.GetUserByUsername(username)
	if err != nil {
		return errors.New("invalid username or password")
	}

	if user.LockedUntil != nil &&
		time.Now().Before(*user.LockedUntil) {

		return fmt.Errorf(
			"account locked until %s",
			user.LockedUntil.Format(time.RFC1123),
		)
	}

	err = VerifyPassword(
		user.PasswordHash,
		password,
	)

	if err != nil {

		newAttempts := user.FailedAttempts + 1

		_ = repo.IncrementFailedAttempts(user.ID)

		if newAttempts >= MaxFailedAttempts {

			lockUntil := time.Now().
				Add(LockDuration)

			_ = repo.LockAccount(
				user.ID,
				lockUntil,
			)

			return fmt.Errorf(
				"account locked for %d minutes",
				int(LockDuration.Minutes()),
			)
		}

		return fmt.Errorf(
			"invalid username or password (%d/%d attempts)",
			newAttempts,
			MaxFailedAttempts,
		)
	}

	_ = repo.ResetFailedAttempts(user.ID)

	_ = repo.UpdateLastLogin(user.ID)

	return nil
}
