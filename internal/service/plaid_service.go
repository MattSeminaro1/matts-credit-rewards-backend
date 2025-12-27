package service

// PlaidService defines the interface for Plaid operations
type PlaidService interface {
	CreateLinkToken(userID string) (string, error)
	ExchangePublicToken(userID string, publicToken string) error
}
