package service

import (
	"errors"
	"matts-credit-rewards-app/backend/internal/auth"
	"matts-credit-rewards-app/backend/internal/models"
	"matts-credit-rewards-app/backend/internal/repository"

	"github.com/google/uuid"
)

func Signup(email, password, name string) error {
	if _, err := repository.GetUserByEmail(email); err == nil {
		return errors.New("user already exists")
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: hash,
		Name:         name,
	}

	return repository.CreateUser(user)
}

func Login(email, password string) (*models.User, error) {
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		// Could be sql.ErrNoRows — user not found
		return nil, errors.New("invalid email or password")
	}

	if !auth.CheckPassword(user.PasswordHash, password) {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
