package service

import (
	"context"
	"errors"

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
		Products:      []plaid.Products{plaid.PRODUCTS_AUTH},
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
func (s *TokenServiceImpl) ExchangePublicToken(publicToken string) (string, error) {
	req := plaid.NewItemPublicTokenExchangeRequest(publicToken)
	resp, _, err := s.PlaidClient.PlaidApi.ItemPublicTokenExchange(context.Background()).
		ItemPublicTokenExchangeRequest(*req).
		Execute()
	if err != nil {
		return "", err
	}
	return resp.GetAccessToken(), nil
}
