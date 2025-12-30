package api

import (
	"log"
	"net/http"

	"matts-credit-rewards-app/backend/internal/models"
	"matts-credit-rewards-app/backend/internal/repository"
	"matts-credit-rewards-app/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type PlaidHandler struct {
	PlaidService service.PlaidService
}

// POST /create_link_token
func (h *PlaidHandler) CreateLinkTokenHandler(c *gin.Context) {
	var req models.CreateLinkTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	log.Printf("Creating link token for user: %s", req.UserID)
	linkToken, err := h.PlaidService.CreateLinkToken(req.UserID)
	log.Printf("Link token creation result: %s, error: %v", linkToken, err)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create link token"})
		return
	}

	c.JSON(http.StatusOK, models.CreateLinkTokenResponse{LinkToken: linkToken})
}

// POST /exchange_public_token
func (h *PlaidHandler) ExchangePublicTokenHandler(c *gin.Context) {
	var req models.ExchangePublicTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Exchange public token and store new account
	err := h.PlaidService.ExchangePublicToken(req.UserID, req.PublicToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to link account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "account successfully added"})
}

// GET /accounts
func (h *PlaidHandler) GetAccountsHandler(c *gin.Context) {
	userID := c.Query("userId")
	log.Printf("userID: %s", userID)

	accountType := c.Query("type")
	log.Printf("accountType: %s", accountType)

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	var acctTypePtr *string
	if accountType != "" {
		acctTypePtr = &accountType
	}

	accounts, err := repository.GetAccountsByUserAndType(userID, acctTypePtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch accounts"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// GET /transactions
func (h *PlaidHandler) GetTransactionsHandler(c *gin.Context) {
	userID := c.Query("userId")
	log.Printf("userID: %s", userID)

	accountId := c.Query("accountId")
	log.Printf("accountId: %s", accountId)

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	var accountIdPtr *string
	if accountId != "" {
		accountIdPtr = &accountId
	}

	accounts, err := repository.GetAccountsByUserAndType(userID, accountIdPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch accounts"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}
