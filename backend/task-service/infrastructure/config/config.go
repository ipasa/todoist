package config

import "os"

type Config struct {
	Port                string
	DatabaseURL         string
	RabbitMQURL         string
	JWTSecret           string
	JWTExpiry           string
	RefreshTokenExpiry   string
}

func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "8002"),
		DatabaseURL:         getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/todoist_tasks?sslmode=disable"),
		RabbitMQURL:         getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		JWTSecret:           getEnv("JWT_SECRET", "dev_secret_key_change_in_production_please"),
		JWTExpiry:           getEnv("JWT_EXPIRY", "15m"),
		RefreshTokenExpiry:   getEnv("REFRESH_TOKEN_EXPIRY", "168h"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
