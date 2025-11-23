package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GatewayUrl         string
	ServerPort       string
}

func Load() (*Config, error) {
	godotenv.Load() 
	
	cfg := &Config{
		GatewayUrl:   getEnv("GATEWAY_URL", "http://localhost:8080"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.GatewayUrl == "" {
		return fmt.Errorf("GATEWAY_URL is required")
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