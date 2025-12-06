package main

import (
	"log"
	"matts-credit-rewards-app/backend/internal/api"
	"matts-credit-rewards-app/backend/internal/config"
	"matts-credit-rewards-app/backend/internal/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config from .env
	cfg := config.LoadConfig("C:/Matts-Credit-Rewards-App/backend/.env")

	// Initialize Postgres connection
	if err := db.Init(cfg.DSN()); err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	// Create Gin router
	r := gin.Default()

	// Enable CORS for your frontend dev server
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // adjust port if frontend is different
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// Routes
	r.POST("/signup", api.SignupHandler)
	r.POST("/login", api.LoginHandler)

	// Start server
	log.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
