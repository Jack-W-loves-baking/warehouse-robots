package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server Configuration
	Server ServerConfig

	// Robot SDK Configuration
	Robot RobotConfig

	// Logging Configuration
	Log LogConfig

	// CORS Configuration
	CORS CORSConfig

	// Environment
	Environment string
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port      string
	Host      string
	AdminPort string
}

// RobotConfig holds facades SDK-related configuration
type RobotConfig struct {
	EnableMock bool
}

// LogConfig holds logging-related configuration
type LogConfig struct {
	Level string
}

// CORSConfig holds CORS-related configuration
type CORSConfig struct {
	AllowedOrigins string
	AllowedMethods string
	AllowedHeaders string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	config := &Config{
		Server: ServerConfig{
			Port:      getEnv("PORT", "8080"),
			AdminPort: getEnv("ADMIN_PORT", "8081"),
			Host:      getEnv("HOST", "localhost"),
		},
		Robot: RobotConfig{
			EnableMock: getEnv("ENABLE_MOCK_ROBOT_SDK", "false") == "true",
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
			AllowedMethods: getEnv("CORS_ALLOWED_METHODS", "GET,POST,DELETE,OPTIONS"),
			AllowedHeaders: getEnv("CORS_ALLOWED_HEADERS", "Content-Type"),
		},
		Environment: getEnv("ENV", "development"),
	}

	return config
}

// Helper functions to get environment variables with default values
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
