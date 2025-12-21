package service

import (
	"context"
	"errors"

	"matts-credit-rewards-app/backend/internal/repository"

	"github.com/plaid/plaid-go/plaid"
)

var ErrUserNotFound = errors.New("user not found")

// LinkServiceImpl is the concrete implementation of PlaidService
type LinkServiceImpl struct {
	PlaidClient   *plaid.APIClient
	AppName       string
	Language      string
	Products      []plaid.Products
	Customization string
}

// NewLinkServiceImpl creates a new LinkServiceImpl
func NewLinkServiceImpl(client *plaid.APIClient) *LinkServiceImpl {
	return &LinkServiceImpl{
		PlaidClient:   client,
		AppName:       "Matts Credit Rewards",
		Language:      "en",
		Products:      []plaid.Products{plaid.PRODUCTS_AUTH},
		Customization: "default",
	}
}

// CreateLinkToken implements PlaidService
func (s *LinkServiceImpl) CreateLinkToken(userID string) (string, error) {
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
func (s *LinkServiceImpl) ExchangePublicToken(publicToken string) (string, error) {
	req := plaid.NewItemPublicTokenExchangeRequest(publicToken)
	resp, _, err := s.PlaidClient.PlaidApi.ItemPublicTokenExchange(context.Background()).
		ItemPublicTokenExchangeRequest(*req).
		Execute()
	if err != nil {
		return "", err
	}
	return resp.GetAccessToken(), nil
}
