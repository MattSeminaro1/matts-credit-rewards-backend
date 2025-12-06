package repository

import (
	"matts-credit-rewards-app/backend/internal/db"
	"matts-credit-rewards-app/backend/internal/models"
)

func GetUserByEmail(email string) (*models.User, error) {
	u := &models.User{}
	row := db.DB.QueryRow("SELECT id, email, password_hash, name, created_at, updated_at FROM rewards.users WHERE email = $1", email)
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func CreateUser(u *models.User) error {
	_, err := db.DB.Exec("INSERT INTO rewards.users (id, email, password_hash, name) VALUES ($1, $2, $3, $4)",
		u.ID, u.Email, u.PasswordHash, u.Name)
	if err != nil {
		return err
	}
	return nil
}
