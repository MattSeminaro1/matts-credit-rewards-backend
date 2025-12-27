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
