package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
	SSLMode  string
}

func LoadPostgresConfig(envFile string) *PostgresConfig {
	// Load the env file
	if err := godotenv.Load(envFile); err != nil {
		log.Println("Could not load .env file, using system environment variables")
	}

	sslMode := os.Getenv("SSLMode")
	if sslMode == "" {
		sslMode = "disable" // default
	}

	return &PostgresConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
		SSLMode:  sslMode,
	}
}

// Returns the DSN for pq (Postgres driver)
func (c *PostgresConfig) DSN() string {
	log.Println("Postgres DSN:", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode))
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.SSLMode)
}
