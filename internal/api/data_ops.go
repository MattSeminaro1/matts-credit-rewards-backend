package api

import (
	"log"
	"net/http"

	"matts-credit-rewards-app/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

// GET /accounts
func GetAccountsHandler(c *gin.Context) {
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
func GetTransactionsHandler(c *gin.Context) {
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
