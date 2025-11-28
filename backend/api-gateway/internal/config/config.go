package config

import "os"

// Config holds the API Gateway configuration
type Config struct {
	Port                   string
	AuthServiceURL         string
	TaskServiceURL         string
	ProjectServiceURL      string
	NotificationServiceURL string
	WebSocketGatewayURL    string
	JWTSecret              string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:                   getEnv("PORT", "8000"),
		AuthServiceURL:         getEnv("AUTH_SERVICE_URL", "http://localhost:8001"),
		TaskServiceURL:         getEnv("TASK_SERVICE_URL", "http://localhost:8002"),
		ProjectServiceURL:      getEnv("PROJECT_SERVICE_URL", "http://localhost:8003"),
		NotificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8004"),
		WebSocketGatewayURL:    getEnv("WEBSOCKET_GATEWAY_URL", "ws://localhost:8005"),
		JWTSecret:              getEnv("JWT_SECRET", "dev_secret_key_change_in_production_please"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
