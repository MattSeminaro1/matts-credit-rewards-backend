package api

import (
	"log"
	"net/http"

	"matts-credit-rewards-app/backend/internal/models"
	"matts-credit-rewards-app/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type PlaidTokenHandler struct {
	PlaidService service.PlaidTokenService
}

// POST /create_link_token
func (h *PlaidTokenHandler) CreateLinkTokenHandler(c *gin.Context) {
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
func (h *PlaidTokenHandler) ExchangePublicTokenHandler(c *gin.Context) {
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
