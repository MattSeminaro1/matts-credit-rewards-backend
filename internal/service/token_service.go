package service

// PlaidTokenService defines the interface for Plaid operations
type PlaidTokenService interface {
	CreateLinkToken(userID string) (string, error)
	ExchangePublicToken(userID string, publicToken string) error
}
