package service

import (
	"context"
	"log"

	"matts-credit-rewards-app/backend/internal/models"
	"matts-credit-rewards-app/backend/internal/repository"
	"matts-credit-rewards-app/backend/internal/util"

	"github.com/plaid/plaid-go/plaid"
)

// TransactionSyncService handles syncing transactions from Plaid
type TransactionSyncService struct {
	PlaidClient *plaid.APIClient
}

// NewTransactionSyncService creates a new TransactionSyncService
func NewTransactionSyncService(client *plaid.APIClient) *TransactionSyncService {
	return &TransactionSyncService{PlaidClient: client}
}

// SyncItem syncs transactions for a given Plaid item
func (s *TransactionSyncService) SyncItem(
	ctx context.Context,
	itemID string,
	accessToken string,
) error {
	log.Printf("Starting sync for item %s", itemID)
	cursor, err := repository.GetTransactionsCursor(itemID)
	if err != nil {
		return err
	}
	log.Printf("Using cursor: %v", cursor)
	for {
		log.Printf("starting for loop")
		req := plaid.TransactionsSyncRequest{
			AccessToken: accessToken,
		}
		if cursor != nil {
			req.Cursor = cursor
		}
		log.Printf("Created request: %+v", req)
		resp, _, err := s.PlaidClient.PlaidApi.
			TransactionsSync(ctx).
			TransactionsSyncRequest(req).
			Execute()
		if err != nil {
			return err
		}

		// Added + Modified
		log.Printf("starting inner for loop")
		for _, tx := range append(resp.GetAdded(), resp.GetModified()...) {
			log.Printf("Get Account ID")
			accountID, err := repository.
				GetAccountIDByPlaidAccountID(tx.AccountId)
			if err != nil {
				return err
			}

			date, err := util.ParseDate(tx.Date)
			if err != nil {
				return err
			}

			model := models.Transaction{
				AccountID:          accountID,
				PlaidTransactionID: tx.TransactionId,
				Name:               tx.Name,
				Amount:             tx.Amount,
				Date:               date,
				Pending:            tx.Pending,
				IsoCurrencyCode:    tx.IsoCurrencyCode.Get(),
				Category:           util.CategoryToString(tx.Category),
			}
			log.Print("inserting transaction")
			log.Printf("Upserting transaction: %+v", model)
			if err := repository.UpsertTransaction(model); err != nil {
				log.Printf("Upsert failed: %v", err)
				return err
			}
		}

		// Removed
		log.Printf("removing transactions")
		for _, r := range resp.GetRemoved() {
			if err := repository.
				DeleteTransactionByPlaidID(r.TransactionId); err != nil {
				return err
			}
		}
		log.Printf("updating transactions cursor")
		next := resp.GetNextCursor()
		if err := repository.
			UpdateTransactionsCursor(itemID, next); err != nil {
			return err
		}
		cursor = &next

		if !resp.GetHasMore() {
			break
		}
	}

	return nil
}
