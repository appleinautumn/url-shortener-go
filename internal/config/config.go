package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string
	DBFile  string
}

func LoadConfig() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", "error", err)
		// Continue even if .env is not found, environment variables may be set elsewhere
	}

	dbFile := os.Getenv("DB_FILE")
	if dbFile == "" {
		return nil, &ConfigError{"DB_FILE must be set in the environment"}
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	return &Config{
		AppEnv:  appEnv,
		AppPort: appPort,
		DBFile:  dbFile,
	}, nil
}

type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return e.Message
}
