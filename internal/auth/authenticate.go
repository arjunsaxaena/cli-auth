package auth

import (
	"cli-auth/internal/config"
	"cli-auth/internal/models"
	"errors"
	"fmt"
	"time"
)

func AuthenticatePassword(
	repo *Repository,
	username string,
	password string,
) (*models.User, error) {

	user, err := repo.GetUserByUsername(
		username,
	)

	if err != nil {
		return nil,
			errors.New(
				"invalid username or password",
			)
	}

	if user.LockedUntil != nil &&
		time.Now().Before(*user.LockedUntil) {

		return nil,
			fmt.Errorf(
				"account locked until %s",
				user.LockedUntil.Format(
					time.RFC1123,
				),
			)
	}

	err = VerifyPassword(
		user.PasswordHash,
		password,
	)

	if err != nil {

		newAttempts :=
			user.FailedAttempts + 1

		_ = repo.IncrementFailedAttempts(
			user.ID,
		)

		if newAttempts >= config.App.MaxFailedAttempt {

			lockUntil :=
				time.Now().
					Add(config.App.LockDuration)

			_ = repo.LockAccount(
				user.ID,
				lockUntil,
			)

			return nil,
				fmt.Errorf(
					"account locked for %d minutes",
					int(
						config.App.LockDuration.Minutes(),
					),
				)
		}

		return nil,
			fmt.Errorf(
				"invalid username or password (%d/%d attempts)",
				newAttempts,
				config.App.MaxFailedAttempt,
			)
	}

	return user, nil
}

func CompleteLogin(
	repo *Repository,
	userID int64,
) error {

	_ = repo.ResetFailedAttempts(
		userID,
	)

	return repo.UpdateLastLogin(
		userID,
	)
}
