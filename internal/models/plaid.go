package models

import "time"

type PlaidItem struct {
	ID              string    `db:"id"`
	UserID          string    `db:"user_id"`
	PlaidItemID     string    `db:"plaid_item_id"`
	AccessToken     string    `db:"access_token"`
	InstitutionName string    `db:"institution_name"`
	InstitutionID   string    `db:"institution_id"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

type Account struct {
	ID             string    `db:"id"`
	ItemID         string    `db:"item_id"`
	PlaidAccountID string    `db:"plaid_account_id"`
	Name           string    `db:"name"`
	OfficialName   string    `db:"official_name"`
	Mask           string    `db:"mask"`
	Type           string    `db:"type"`
	Subtype        string    `db:"subtype"`
	CurrentBalance float64   `db:"current_balance"`
	AvailableBal   float64   `db:"available_balance"`
	Currency       string    `db:"currency"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
