package config

import (
	"os"
	"time"
)

// Config holds the application configuration
type Config struct {
	Port               string
	DatabaseURL        string
	RabbitMQURL        string
	RedisURL           string
	JWTSecret          string
	JWTExpiry          time.Duration
	RefreshTokenExpiry time.Duration
	GoogleClientID     string
	GoogleClientSecret string
	GitHubClientID     string
	GitHubClientSecret string
}

// Load loads configuration from environment variables
func Load() *Config {
	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "15m"))
	refreshTokenExpiry, _ := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRY", "168h"))

	return &Config{
		Port:               getEnv("PORT", "8001"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://todoist:todoist_dev@localhost:5432/auth_db?sslmode=disable"),
		RabbitMQURL:        getEnv("RABBITMQ_URL", "amqp://todoist:todoist_dev@localhost:5672/"),
		RedisURL:           getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:          getEnv("JWT_SECRET", "dev_secret_key_change_in_production_please"),
		JWTExpiry:          jwtExpiry,
		RefreshTokenExpiry: refreshTokenExpiry,
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
