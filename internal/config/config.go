package config

func LoadConfig(envFile string) *PostgresConfig {
	// For now, we only have MySQL, but in future you could select based on env var
	return LoadPostgresConfig(envFile)
}
