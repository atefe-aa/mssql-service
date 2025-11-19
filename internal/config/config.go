package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database configuration
	DBServer         string
	DBPort           string
	DBName           string
	DBAuthMode       string // "windows" or "sql"
	DBUser           string // Only for SQL auth
	DBPassword       string // Only for SQL auth
	
	// Server configuration
	ServerPort       string
}

func Load() (*Config, error) {
	godotenv.Load() 
	
	cfg := &Config{
		DBServer:   getEnv("DB_SERVER", "localhost"),
		DBPort:     getEnv("DB_PORT", "1433"),
		DBName:     getEnv("DB_NAME", ""),
		DBAuthMode: getEnv("DB_AUTH_MODE", "windows"), // default to windows auth
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	if c.DBAuthMode != "windows" && c.DBAuthMode != "sql" {
		return fmt.Errorf("DB_AUTH_MODE must be either 'windows' or 'sql'")
	}

	if c.DBAuthMode == "sql" {
		if c.DBUser == "" || c.DBPassword == "" {
			return fmt.Errorf("DB_USER and DB_PASSWORD are required when using SQL authentication")
		}
	}

	if _, err := strconv.Atoi(c.ServerPort); err != nil {
		return fmt.Errorf("SERVER_PORT must be a valid port number")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}