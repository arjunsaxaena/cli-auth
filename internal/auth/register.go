package auth

import (
	"errors"
)

func Register(
	repo *Repository,
	username string,
	password string,
) error {

	if username == "" {
		return errors.New("username cannot be empty")
	}

	if password == "" {
		return errors.New("password cannot be empty")
	}

	existing, _ := repo.GetUserByUsername(username)

	if existing != nil {
		return errors.New("username already exists")
	}

	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	return repo.CreateUser(
		username,
		hash,
	)
}
