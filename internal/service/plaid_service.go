package service

// PlaidService defines the interface for Plaid operations
type PlaidService interface {
	CreateLinkToken(userID string) (string, error)
	ExchangePublicToken(publicToken string) (string, error)
}
