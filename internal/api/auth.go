package api

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"matts-credit-rewards-app/backend/internal/models"
	"matts-credit-rewards-app/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// SignupHandler handles user signup
func SignupHandler(c *gin.Context) {
	// Read the raw body for debugging
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	log.Println("Raw signup body:", string(body))

	// Reset the request body so ShouldBindJSON can read it
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Parse JSON
	var req models.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("BindJSON error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	log.Println("Parsed signup request:", req)

	// Call the service layer
	err = service.Signup(req.Email, req.Password, req.Name)
	if err != nil {
		log.Println("Signup service error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// LoginHandler handles user login
func LoginHandler(c *gin.Context) {
	// Read the raw body for debugging
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}
	log.Println("Raw login body:", string(body))

	// Reset the request body so ShouldBindJSON can read it
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Parse JSON
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("BindJSON error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	log.Println("Parsed login request:", req)

	// Call the service layer
	user, err := service.Login(req.Email, req.Password)
	if err != nil {
		log.Println("Login service error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// TODO: Generate JWT token
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
		"user": gin.H{
			"uid":   user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
