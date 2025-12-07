package service

import (
	"context"
	"matts-credit-rewards-app/backend/internal/models"
	"matts-credit-rewards-app/backend/internal/repository"

	plaid "github.com/plaid/plaid-go/plaid"
)

// ItemService handles Plaid item operations
type ItemService struct {
	Client     *plaid.APIClient
	Repository *repository.PlaidRepository
}

// Constructor
func NewItemService(client *plaid.APIClient, repo *repository.PlaidRepository) *ItemService {
	return &ItemService{
		Client:     client,
		Repository: repo,
	}
}

// LinkTokenMetadata is your simplified struct for passing Plaid metadata
type LinkTokenMetadata struct {
	Institution struct {
		Name          string
		InstitutionId string
	}
	Accounts []struct {
		Id           string
		Name         string
		OfficialName string
		Mask         string
		Type         string
		Subtype      string
		Balances     struct {
			Current         float64
			Available       float64
			IsoCurrencyCode string
		}
	}
}

// ExchangePublicToken exchanges a public token for an access token and stores the item/accounts
func (s *ItemService) ExchangePublicToken(userID, publicToken string, metadata *LinkTokenMetadata) error {
	ctx := context.Background()

	// Use the PlaidApi sub-client
	plaidResp, _, err := s.Client.PlaidApi.ItemPublicTokenExchange(ctx).
		ItemPublicTokenExchangeRequest(plaid.ItemPublicTokenExchangeRequest{
			PublicToken: publicToken,
		}).
		Execute()
	if err != nil {
		return err
	}

	accessToken := plaidResp.GetAccessToken()
	itemID := plaidResp.GetItemId()

	// Store PlaidItem
	item := &models.PlaidItem{
		UserID:          userID,
		PlaidItemID:     itemID,
		AccessToken:     accessToken,
		InstitutionName: metadata.Institution.Name,
		InstitutionID:   metadata.Institution.InstitutionId,
	}

	if err := s.Repository.CreateItem(item); err != nil {
		return err
	}

	// Store accounts
	for _, acc := range metadata.Accounts {
		a := &models.Account{
			ItemID:         item.ID,
			PlaidAccountID: acc.Id,
			Name:           acc.Name,
			OfficialName:   acc.OfficialName,
			Mask:           acc.Mask,
			Type:           acc.Type,
			Subtype:        acc.Subtype,
			CurrentBalance: acc.Balances.Current,
			AvailableBal:   acc.Balances.Available,
			Currency:       acc.Balances.IsoCurrencyCode,
		}
		if err := s.Repository.CreateAccount(a); err != nil {
			return err
		}
	}

	return nil
}
