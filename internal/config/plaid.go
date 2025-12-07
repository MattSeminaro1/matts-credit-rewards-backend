package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/plaid/plaid-go/plaid"
)

// PlaidConfig holds the Plaid client and environment settings
type PlaidConfig struct {
	ClientID    string // Renamed to clarify
	Secret      string // Renamed to clarify
	Env         string
	PlaidClient *plaid.APIClient // Added to store the initialized client
}

// LoadPlaidConfig loads configuration from the .env file into the struct
func LoadPlaidConfig(envFile string) *PlaidConfig {
	// Load the env file
	if err := godotenv.Load(envFile); err != nil {
		log.Println("Could not load .env file, using system environment variables")
	}

	return &PlaidConfig{
		ClientID: os.Getenv("PLAID_CLIENT_ID"),
		Secret:   os.Getenv("PLAID_SECRET"),
		Env:      os.Getenv("PLAID_ENV"),
	}
}

// InitializePlaidClient uses the loaded config values to create a functional Plaid API Client
func (c *PlaidConfig) InitializePlaidClient() *plaid.APIClient {
	if c.ClientID == "" || c.Secret == "" || c.Env == "" {
		log.Fatal("Plaid environment variables (PLAID_CLIENT_ID, PLAID_SECRET, PLAID_ENV) must be set.")
	}

	var plaidEnv plaid.Environment
	switch c.Env {
	case "sandbox":
		plaidEnv = plaid.Sandbox
	case "development":
		plaidEnv = plaid.Development
	case "production":
		plaidEnv = plaid.Production
	default:
		log.Fatalf("Invalid PLAID_ENV '%s'. Use 'sandbox', 'development', or 'production'.", c.Env)
	}

	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", c.ClientID)
	configuration.AddDefaultHeader("PLAID-SECRET", c.Secret)
	configuration.UseEnvironment(plaidEnv)

	apiClient := plaid.NewAPIClient(configuration)
	c.PlaidClient = apiClient // Optionally store it back in the struct

	return apiClient
}
