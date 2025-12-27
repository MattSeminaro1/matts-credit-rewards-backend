package api

import (
	"log"
	"net/http"

	"matts-credit-rewards-app/backend/internal/models"
	"matts-credit-rewards-app/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type PlaidHandler struct {
	PlaidService service.PlaidService
}

// POST /api/create_link_token
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

// POST /api/exchange_public_token
func (h *PlaidHandler) ExchangePublicTokenHandler(c *gin.Context) {
	var req models.ExchangePublicTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Exchange public token and get accounts
	accounts, err := h.PlaidService.ExchangePublicToken(req.UserID, req.PublicToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange public token"})
		return
	}

	// Return accounts to frontend
	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}
