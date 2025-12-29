package models

type CreateLinkTokenRequest struct {
	UserID string `json:"user_id"`
}

type CreateLinkTokenResponse struct {
	LinkToken string `json:"link_token"`
}

type ExchangePublicTokenRequest struct {
	UserID      string `json:"user_id"`
	PublicToken string `json:"public_token"`
}

type Account struct {
	ItemID           string   // item_id
	PlaidAccountID   string   // plaid_account_id
	Name             string   // name
	OfficialName     *string  // official_name (nullable)
	Mask             *string  // mask (nullable)
	Type             string   // type
	Subtype          *string  // subtype (nullable)
	CurrentBalance   *float32 // current_balance (nullable)
	AvailableBalance *float32 // available_balance (nullable)
	Currency         *string  // currency (nullable)
}
