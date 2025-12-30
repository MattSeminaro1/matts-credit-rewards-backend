package main

import (
	"log"
	"matts-credit-rewards-app/backend/internal/api"
	"matts-credit-rewards-app/backend/internal/config"
	"matts-credit-rewards-app/backend/internal/db"
	"matts-credit-rewards-app/backend/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load Plaid config
	plaidCfg := config.LoadPlaidConfig("C:/Matts-Credit-Rewards-App/backend/.env")

	// Initialize Plaid client
	plaidClient := plaidCfg.InitializePlaidClient()
	log.Printf("Successfully loaded Plaid client for %s environment.\n", plaidCfg.Env)

	// Load Postgres config
	postgresCfg := config.LoadPostgresConfig("C:/Matts-Credit-Rewards-App/backend/.env")
	if err := db.Init(postgresCfg.DSN()); err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	// Initialize TokenServiceImpl & Handler
	tokenService := service.NewTokenServiceImpl(plaidClient)
	var plaidSvc service.PlaidService = tokenService

	plaidHandler := &api.PlaidHandler{PlaidService: plaidSvc}

	// Create Gin router
	r := gin.Default()

	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// Routes
	r.POST("/signup", api.SignupHandler)
	r.POST("/login", api.LoginHandler)
	r.POST("/create_link_token", plaidHandler.CreateLinkTokenHandler)
	r.POST("/exchange_public_token", plaidHandler.ExchangePublicTokenHandler)
	r.GET("/accounts", plaidHandler.GetAccountsHandler)
	r.GET("/transactions", plaidHandler.GetTransactionsHandler)

	// Start server
	log.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
