package config

import (
	// "log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	HasuraEndpoint    string
	HasuraAdminSecret string
	JwtSecret        string
	JwtAlgorithm     string
}

func NewConfig() *Config {
	// Try loading .env file but don't fail if it doesn't exist
	_ = godotenv.Load() // Changed from Fatalf to silent load
	
	return &Config{
		HasuraEndpoint:    getEnv("HASURA_ENDPOINT", ""),
		HasuraAdminSecret: getEnv("HASURA_ADMIN_SECRET", ""),
		JwtSecret:        getEnv("JWT_SECRET", ""),
		JwtAlgorithm:     getEnv("JWT_ALGORITHM", "HS256"),
	}
}

// Helper function to get env with fallback
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}