package repository

import (
	"database/sql"
	"matts-credit-rewards-app/backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type PlaidRepository struct {
	DB *sql.DB
}

func NewPlaidRepository(db *sql.DB) *PlaidRepository {
	return &PlaidRepository{DB: db}
}

func (r *PlaidRepository) CreateItem(item *models.PlaidItem) error {
	if item.ID == "" {
		item.ID = uuid.New().String()
	}
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	_, err := r.DB.Exec(`
		INSERT INTO plaid_items (id, user_id, plaid_item_id, access_token, institution_name, institution_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`,
		item.ID, item.UserID, item.PlaidItemID, item.AccessToken, item.InstitutionName, item.InstitutionID, item.CreatedAt, item.UpdatedAt)
	return err
}

func (r *PlaidRepository) CreateAccount(acc *models.Account) error {
	if acc.ID == "" {
		acc.ID = uuid.New().String()
	}
	now := time.Now()
	acc.CreatedAt = now
	acc.UpdatedAt = now

	_, err := r.DB.Exec(`
		INSERT INTO accounts (id, item_id, plaid_account_id, name, official_name, mask, type, subtype, current_balance, available_balance, currency, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`,
		acc.ID, acc.ItemID, acc.PlaidAccountID, acc.Name, acc.OfficialName, acc.Mask,
		acc.Type, acc.Subtype, acc.CurrentBalance, acc.AvailableBal, acc.Currency, acc.CreatedAt, acc.UpdatedAt)
	return err
}
