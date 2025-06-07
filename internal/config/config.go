package config

import (
	"fmt"
	"os"
)

// Config represents the application configuration
type Config struct {
	// Server configuration
	Port string `json:"port"`

	// AWS configuration
	AWSRegion string `json:"aws_region"`

	// Bedrock configuration
	BedrockModelID string `json:"bedrock_model_id"`

	// Logging configuration
	LogLevel string `json:"log_level"`
}

// Load loads configuration from environment variables with default values
func Load() (*Config, error) {
	config := &Config{
		Port:           getEnvWithDefault("PORT", "8080"),
		AWSRegion:      getEnvWithDefault("AWS_REGION", "ap-northeast-1"),
		BedrockModelID: getEnvWithDefault("BEDROCK_MODEL_ID", "amazon.nova-lite-v1:0"),
		LogLevel:       getEnvWithDefault("LOG_LEVEL", "info"),
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// validate performs basic validation of the configuration
func (c *Config) validate() error {
	if c.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}

	if c.AWSRegion == "" {
		return fmt.Errorf("AWS region cannot be empty")
	}

	if c.BedrockModelID == "" {
		return fmt.Errorf("bedrock model ID cannot be empty")
	}

	return nil
}

// getEnvWithDefault returns the value of an environment variable or a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
