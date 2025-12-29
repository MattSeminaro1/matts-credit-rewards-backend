package repository

import (
	"matts-credit-rewards-app/backend/internal/db"
	"matts-credit-rewards-app/backend/internal/models"
)

func UpsertPlaidItem(
	userID, plaidItemID, accessToken, institutionID, institutionName string,
) (string, error) {
	var id string
	err := db.DB.QueryRow(`
		INSERT INTO rewards.plaid_items (user_id, plaid_item_id, access_token, institution_id, institution_name)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (plaid_item_id)
		DO UPDATE SET
			access_token = EXCLUDED.access_token,
			institution_id = EXCLUDED.institution_id,
			institution_name = EXCLUDED.institution_name,
			updated_at = NOW()
		RETURNING id
	`, userID, plaidItemID, accessToken, institutionID, institutionName).Scan(&id)
	return id, err
}

func UpsertAccounts(accounts []models.Account) error {
	if len(accounts) == 0 {
		return nil
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, a := range accounts {
		_, err := tx.Exec(`
			INSERT INTO rewards.accounts (
				item_id,
				plaid_account_id,
				name,
				official_name,
				mask,
				type,
				subtype,
				current_balance,
				available_balance,
				currency
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
			ON CONFLICT (plaid_account_id)
			DO UPDATE SET
				current_balance = EXCLUDED.current_balance,
				available_balance = EXCLUDED.available_balance,
				updated_at = NOW()
		`,
			a.ItemID,
			a.PlaidAccountID,
			a.Name,
			a.OfficialName,
			a.Mask,
			a.Type,
			a.Subtype,
			a.CurrentBalance,
			a.AvailableBalance,
			a.Currency,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// repository/accounts.go
func GetAccountsByUserAndType(userID string, accountType *string) ([]models.Account, error) {
	query := `
		SELECT 
			item_id, plaid_account_id, name, official_name, mask, type, subtype, current_balance, available_balance, currency
		FROM rewards.accounts
		WHERE item_id IN (
			SELECT id FROM rewards.plaid_items WHERE user_id = $1
		)
	`

	args := []interface{}{userID}

	if accountType != nil {
		query += " AND type = $2"
		args = append(args, *accountType)
	}

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var a models.Account
		err := rows.Scan(
			&a.ItemID,
			&a.PlaidAccountID,
			&a.Name,
			&a.OfficialName,
			&a.Mask,
			&a.Type,
			&a.Subtype,
			&a.CurrentBalance,
			&a.AvailableBalance,
			&a.Currency,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}

	return accounts, nil
}
