package config

import (
	"fmt"
	"os"
)

// Config represents the application configuration
type Config struct {
	// Server configuration
	Port string `json:"port"`

	// Database configuration
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`

	// AWS configuration
	AWSRegion string `json:"aws_region"`

	// Bedrock configuration
	BedrockModelID string `json:"bedrock_model_id"`

	// Bedrock Knowledge Base configuration (for V2)
	KnowledgeBaseID  string `json:"knowledge_base_id"`
	EmbeddingModelID string `json:"embedding_model_id"`

	// OpenSearch configuration (for V2)
	OpenSearchEndpoint  string `json:"opensearch_endpoint"`
	OpenSearchUsername  string `json:"opensearch_username"`
	OpenSearchPassword  string `json:"opensearch_password"`
	OpenSearchIndexName string `json:"opensearch_index_name"`

	// Logging configuration
	LogLevel string `json:"log_level"`
}

// Load loads configuration from environment variables with default values
func Load() (*Config, error) {
	config := &Config{
		Port:           getEnvWithDefault("PORT", "8080"),
		DBHost:         getEnvWithDefault("DB_HOST", "localhost"),
		DBPort:         getEnvWithDefault("DB_PORT", "5436"),
		DBUser:         getEnvWithDefault("DB_USER", "postgres"),
		DBPassword:     getEnvWithDefault("DB_PASSWORD", "postgres"),
		DBName:         getEnvWithDefault("DB_NAME", "postgres"),
		AWSRegion:      getEnvWithDefault("AWS_REGION", "ap-northeast-1"),
		BedrockModelID: getEnvWithDefault("BEDROCK_MODEL_ID", "amazon.nova-lite-v1:0"),

		// Bedrock Knowledge Base configuration
		KnowledgeBaseID:  getEnvWithDefault("KNOWLEDGE_BASE_ID", ""),
		EmbeddingModelID: getEnvWithDefault("EMBEDDING_MODEL_ID", "amazon.titan-embed-text-v1"),

		// OpenSearch configuration
		OpenSearchEndpoint:  getEnvWithDefault("OPENSEARCH_ENDPOINT", ""),
		OpenSearchUsername:  getEnvWithDefault("OPENSEARCH_USERNAME", ""),
		OpenSearchPassword:  getEnvWithDefault("OPENSEARCH_PASSWORD", ""),
		OpenSearchIndexName: getEnvWithDefault("OPENSEARCH_INDEX_NAME", "product-vectors"),

		LogLevel: getEnvWithDefault("LOG_LEVEL", "info"),
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

	if c.DBHost == "" {
		return fmt.Errorf("database host cannot be empty")
	}

	if c.DBPort == "" {
		return fmt.Errorf("database port cannot be empty")
	}

	if c.DBUser == "" {
		return fmt.Errorf("database user cannot be empty")
	}

	if c.DBName == "" {
		return fmt.Errorf("database name cannot be empty")
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
