package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GatewayUrl string
	ServerPort string
	APIKey     string
}

func Load() (*Config, error) {
	godotenv.Load()

	cfg := &Config{
		GatewayUrl: getEnv("GATEWAY_URL", "http://localhost:8080"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		APIKey:     getEnv("API_KEY", ""),
	}

	if cfg.APIKey == "" {
		apiKey, err := generateAPIKey()
		if err != nil {
			return nil, fmt.Errorf("failed to generate API key: %w", err)
		}
		cfg.APIKey = apiKey

		if err := saveAPIKeyToEnv(apiKey); err != nil {
			return nil, fmt.Errorf("failed to save API key: %w", err)
		}
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

func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
func saveAPIKeyToEnv(apiKey string) error {
	envMap, _ := godotenv.Read(".env")
	if envMap == nil {
		envMap = make(map[string]string)
	}
	
	envMap["API_KEY"] = apiKey
	
	return godotenv.Write(envMap, ".env")
}