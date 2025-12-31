package repository

import (
	"matts-credit-rewards-app/backend/internal/db"
	"matts-credit-rewards-app/backend/internal/models"
)

// Inster Plaid Items
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

// Inster Plaid Accounts
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

// Get accounts by user ID and optional account type
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

// Get transactions by user ID and optional account ID
func GetTransactionsByUserAndAccount(
	userID string,
	accountID *string,
) ([]models.Transaction, error) {

	query := `
		SELECT
			t.id,
			t.account_id,
			t.plaid_transaction_id,
			t.name,
			t.amount,
			t.iso_currency_code,
			t.category,
			t.date,
			t.pending,
			t.created_at,
			t.updated_at
		FROM transactions t
		JOIN accounts a ON a.id = t.account_id
		JOIN plaid_items pi ON pi.id = a.item_id
		WHERE pi.user_id = $1
	`

	args := []interface{}{userID}

	if accountID != nil {
		query += " AND a.id = $2"
		args = append(args, *accountID)
	}

	query += " ORDER BY t.date DESC"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction

	for rows.Next() {
		var t models.Transaction

		err := rows.Scan(
			&t.ID,
			&t.AccountID,
			&t.PlaidTransactionID,
			&t.Name,
			&t.Amount,
			&t.IsoCurrencyCode,
			&t.Category,
			&t.Date,
			&t.Pending,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// Inster Plaid Transactions
func GetAccountIDByPlaidAccountID(plaidAccountID string) (string, error) {
	var accountID string

	err := db.DB.QueryRow(`
		SELECT id
		FROM rewards.accounts
		WHERE plaid_account_id = $1
	`, plaidAccountID).Scan(&accountID)

	return accountID, err
}

func UpsertTransaction(t models.Transaction) error {
	_, err := db.DB.Exec(`
		INSERT INTO rewards.transactions (
			account_id,
			plaid_transaction_id,
			name,
			amount,
			iso_currency_code,
			category,
			date,
			pending
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (plaid_transaction_id)
		DO UPDATE SET
			name = EXCLUDED.name,
			amount = EXCLUDED.amount,
			iso_currency_code = EXCLUDED.iso_currency_code,
			category = EXCLUDED.category,
			date = EXCLUDED.date,
			pending = EXCLUDED.pending,
			updated_at = NOW()
	`,
		t.AccountID,
		t.PlaidTransactionID,
		t.Name,
		t.Amount,
		t.IsoCurrencyCode,
		t.Category,
		t.Date,
		t.Pending,
	)

	return err
}

func DeleteTransactionByPlaidID(plaidTransactionID *string) error {
	_, err := db.DB.Exec(`
		DELETE FROM rewards.transactions
		WHERE plaid_transaction_id = $1
	`, plaidTransactionID)

	return err
}

func GetTransactionsCursor(itemID string) (*string, error) {
	var cursor *string

	err := db.DB.QueryRow(`
		SELECT transactions_cursor
		FROM rewards.plaid_items
		WHERE id = $1
	`, itemID).Scan(&cursor)

	return cursor, err
}

func UpdateTransactionsCursor(itemID string, cursor string) error {
	_, err := db.DB.Exec(`
		UPDATE rewards.plaid_items
		SET transactions_cursor = $1,
		    updated_at = NOW()
		WHERE id = $2
	`, cursor, itemID)

	return err
}
