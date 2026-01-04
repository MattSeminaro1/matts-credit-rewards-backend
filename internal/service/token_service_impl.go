package service

import (
	"context"
	"errors"
	"log"

	"matts-credit-rewards-app/backend/internal/models"
	"matts-credit-rewards-app/backend/internal/repository"

	"github.com/plaid/plaid-go/plaid"
)

var ErrUserNotFound = errors.New("user not found")

// TokenServiceImpl is the concrete implementation of PlaidService
type TokenServiceImpl struct {
	PlaidClient   *plaid.APIClient
	AppName       string
	Language      string
	Products      []plaid.Products
	Customization string
}

// NewTokenServiceImpl creates a new TokenServiceImpl
func NewTokenServiceImpl(client *plaid.APIClient) *TokenServiceImpl {
	return &TokenServiceImpl{
		PlaidClient:   client,
		AppName:       "Matts Credit Rewards",
		Language:      "en",
		Products:      []plaid.Products{plaid.PRODUCTS_AUTH, plaid.PRODUCTS_TRANSACTIONS},
		Customization: "default",
	}
}

// CreateLinkToken implements PlaidService
func (s *TokenServiceImpl) CreateLinkToken(userID string) (string, error) {
	exists, err := repository.UserExists(userID)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", ErrUserNotFound
	}

	user := plaid.LinkTokenCreateRequestUser{ClientUserId: userID}
	req := plaid.NewLinkTokenCreateRequest(
		s.AppName,
		s.Language,
		[]plaid.CountryCode{plaid.COUNTRYCODE_US},
		user,
	)
	req.SetProducts(s.Products)
	req.SetLinkCustomizationName(s.Customization)

	resp, _, err := s.PlaidClient.PlaidApi.LinkTokenCreate(context.Background()).
		LinkTokenCreateRequest(*req).
		Execute()
	if err != nil {
		return "", err
	}

	return resp.GetLinkToken(), nil
}

// ExchangePublicToken implements PlaidService
func (s *TokenServiceImpl) ExchangePublicToken(userID string, publicToken string) error {
	// 1. Exchange public token
	req := plaid.NewItemPublicTokenExchangeRequest(publicToken)
	resp, _, err := s.PlaidClient.PlaidApi.ItemPublicTokenExchange(context.Background()).
		ItemPublicTokenExchangeRequest(*req).
		Execute()
	if err != nil {
		return err
	}

	accessToken := resp.GetAccessToken()
	itemID := resp.GetItemId()

	log.Printf("Exchanged public token. AccessToken: %s, ItemID: %s", accessToken, itemID)

	// 2. Get item info
	log.Printf("Fetching item info for AccessToken: %s", accessToken)
	itemResp, _, err := s.PlaidClient.PlaidApi.
		ItemGet(context.Background()).
		ItemGetRequest(*plaid.NewItemGetRequest(accessToken)).
		Execute()
	if err != nil {
		return err
	}

	log.Printf("Fetched item info: %+v", itemResp)

	item := itemResp.GetItem()

	institutionID := itemResp.Item.GetInstitutionId()

	institutionName := ""
	if v, ok := item.AdditionalProperties["institution_name"]; ok {
		if s, ok := v.(string); ok {
			institutionName = s
		}
	}

	log.Printf("Institution ID: %s, Name: %s", institutionID, institutionName)
	// 3. Persist plaid_items
	plaidItemID, err := repository.UpsertPlaidItem(
		userID,
		itemID,
		accessToken,
		institutionID,
		institutionName,
	)
	if err != nil {
		return err
	}
	log.Printf("Upserted Plaid Item with ID: %s", plaidItemID)

	// 4. Fetch accounts
	accResp, _, err := s.PlaidClient.PlaidApi.
		AccountsGet(context.Background()).
		AccountsGetRequest(*plaid.NewAccountsGetRequest(accessToken)).
		Execute()
	if err != nil {
		return err
	}
	log.Printf("Fetched accounts: %+v", accResp.GetAccounts())

	accounts := make([]models.Account, 0, len(accResp.GetAccounts()))

	for _, a := range accResp.GetAccounts() {
		acc := models.Account{
			ItemID:         plaidItemID,
			PlaidAccountID: a.AccountId,
			Name:           a.Name,
			Type:           string(a.Type),
		}

		if a.OfficialName.IsSet() {
			acc.OfficialName = a.OfficialName.Get()
		}
		if a.Mask.IsSet() {
			acc.Mask = a.Mask.Get()
		}
		if a.Subtype.IsSet() && a.Subtype.Get() != nil {
			subtypeValue := string(*a.Subtype.Get())
			acc.Subtype = &subtypeValue
		}
		if a.Balances.Current.IsSet() && a.Balances.Current.Get() != nil {
			acc.CurrentBalance = a.Balances.Current.Get()
		}
		if a.Balances.Available.IsSet() && a.Balances.Available.Get() != nil {
			acc.AvailableBalance = a.Balances.Available.Get()
		}
		if a.Balances.IsoCurrencyCode.IsSet() {
			acc.Currency = a.Balances.IsoCurrencyCode.Get()
		}

		accounts = append(accounts, acc)
	}

	// persist accounts
	if err := repository.UpsertAccounts(accounts); err != nil {
		return err
	}

	log.Printf("Upserted %d accounts for Plaid Item ID: %s", len(accounts), plaidItemID)
	log.Printf("Transaction Sync starting for Plaid Item ID: %s", plaidItemID)
	// 🔹 Sync transactions without changing struct
	txSyncService := NewTransactionSyncService(s.PlaidClient)
	if err := txSyncService.SyncItem(
		context.Background(),
		plaidItemID,
		accessToken,
	); err != nil {
		return err
	}

	// return accounts to frontend
	return nil
}
