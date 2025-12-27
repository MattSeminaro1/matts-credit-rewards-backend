package service

import "github.com/plaid/plaid-go/plaid"

// PlaidService defines the interface for Plaid operations
type PlaidService interface {
	CreateLinkToken(userID string) (string, error)
	ExchangePublicToken(userID string, publicToken string) ([]plaid.AccountBase, error)
}
