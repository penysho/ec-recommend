package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original env vars
	originalPort := os.Getenv("PORT")
	originalRegion := os.Getenv("AWS_REGION")
	originalModelID := os.Getenv("BEDROCK_MODEL_ID")
	originalLogLevel := os.Getenv("LOG_LEVEL")

	// Clean up after test
	defer func() {
		os.Setenv("PORT", originalPort)
		os.Setenv("AWS_REGION", originalRegion)
		os.Setenv("BEDROCK_MODEL_ID", originalModelID)
		os.Setenv("LOG_LEVEL", originalLogLevel)
	}()

	t.Run("load with default values", func(t *testing.T) {
		// Clear environment variables
		os.Unsetenv("PORT")
		os.Unsetenv("AWS_REGION")
		os.Unsetenv("BEDROCK_MODEL_ID")
		os.Unsetenv("LOG_LEVEL")

		config, err := Load()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if config.Port != "8080" {
			t.Errorf("Expected port 8080, got %s", config.Port)
		}

		if config.AWSRegion != "ap-northeast-1" {
			t.Errorf("Expected region ap-northeast-1, got %s", config.AWSRegion)
		}

		if config.BedrockModelID != "amazon.nova-lite-v1:0" {
			t.Errorf("Expected default model ID, got %s", config.BedrockModelID)
		}

		if config.LogLevel != "info" {
			t.Errorf("Expected log level info, got %s", config.LogLevel)
		}
	})

	t.Run("load with custom values", func(t *testing.T) {
		// Set custom environment variables
		os.Setenv("PORT", "3000")
		os.Setenv("AWS_REGION", "ap-northeast-1")
		os.Setenv("BEDROCK_MODEL_ID", "anthropic.claude-3-sonnet-20240229-v1:0")
		os.Setenv("LOG_LEVEL", "debug")

		config, err := Load()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if config.Port != "3000" {
			t.Errorf("Expected port 3000, got %s", config.Port)
		}

		if config.AWSRegion != "ap-northeast-1" {
			t.Errorf("Expected region ap-northeast-1, got %s", config.AWSRegion)
		}

		if config.BedrockModelID != "anthropic.claude-3-sonnet-20240229-v1:0" {
			t.Errorf("Expected custom model ID, got %s", config.BedrockModelID)
		}

		if config.LogLevel != "debug" {
			t.Errorf("Expected log level debug, got %s", config.LogLevel)
		}
	})
}

func TestValidate(t *testing.T) {
	t.Run("valid configuration", func(t *testing.T) {
		config := &Config{
			Port:           "8080",
			DBHost:         "localhost",
			DBPort:         "5436",
			DBUser:         "postgres",
			DBPassword:     "postgres",
			DBName:         "postgres",
			AWSRegion:      "us-east-1",
			BedrockModelID: "amazon.nova-lite-v1:0",
			LogLevel:       "info",
		}

		err := config.validate()
		if err != nil {
			t.Errorf("Expected no error for valid config, got %v", err)
		}
	})

	t.Run("empty port", func(t *testing.T) {
		config := &Config{
			Port:           "",
			DBHost:         "localhost",
			DBPort:         "5436",
			DBUser:         "postgres",
			DBPassword:     "postgres",
			DBName:         "postgres",
			AWSRegion:      "us-east-1",
			BedrockModelID: "anthropic.claude-3-haiku-20240307-v1:0",
			LogLevel:       "info",
		}

		err := config.validate()
		if err == nil {
			t.Error("Expected error for empty port")
		}
	})

	t.Run("empty AWS region", func(t *testing.T) {
		config := &Config{
			Port:           "8080",
			DBHost:         "localhost",
			DBPort:         "5436",
			DBUser:         "postgres",
			DBPassword:     "postgres",
			DBName:         "postgres",
			AWSRegion:      "",
			BedrockModelID: "anthropic.claude-3-haiku-20240307-v1:0",
			LogLevel:       "info",
		}

		err := config.validate()
		if err == nil {
			t.Error("Expected error for empty AWS region")
		}
	})

	t.Run("empty model ID", func(t *testing.T) {
		config := &Config{
			Port:           "8080",
			DBHost:         "localhost",
			DBPort:         "5436",
			DBUser:         "postgres",
			DBPassword:     "postgres",
			DBName:         "postgres",
			AWSRegion:      "us-east-1",
			BedrockModelID: "",
			LogLevel:       "info",
		}

		err := config.validate()
		if err == nil {
			t.Error("Expected error for empty model ID")
		}
	})

	t.Run("empty database host", func(t *testing.T) {
		config := &Config{
			Port:           "8080",
			DBHost:         "",
			DBPort:         "5436",
			DBUser:         "postgres",
			DBPassword:     "postgres",
			DBName:         "postgres",
			AWSRegion:      "us-east-1",
			BedrockModelID: "amazon.nova-lite-v1:0",
			LogLevel:       "info",
		}

		err := config.validate()
		if err == nil {
			t.Error("Expected error for empty database host")
		}
	})

	t.Run("empty database port", func(t *testing.T) {
		config := &Config{
			Port:           "8080",
			DBHost:         "localhost",
			DBPort:         "",
			DBUser:         "postgres",
			DBPassword:     "postgres",
			DBName:         "postgres",
			AWSRegion:      "us-east-1",
			BedrockModelID: "amazon.nova-lite-v1:0",
			LogLevel:       "info",
		}

		err := config.validate()
		if err == nil {
			t.Error("Expected error for empty database port")
		}
	})

	t.Run("empty database user", func(t *testing.T) {
		config := &Config{
			Port:           "8080",
			DBHost:         "localhost",
			DBPort:         "5436",
			DBUser:         "",
			DBPassword:     "postgres",
			DBName:         "postgres",
			AWSRegion:      "us-east-1",
			BedrockModelID: "amazon.nova-lite-v1:0",
			LogLevel:       "info",
		}

		err := config.validate()
		if err == nil {
			t.Error("Expected error for empty database user")
		}
	})

	t.Run("empty database name", func(t *testing.T) {
		config := &Config{
			Port:           "8080",
			DBHost:         "localhost",
			DBPort:         "5436",
			DBUser:         "postgres",
			DBPassword:     "postgres",
			DBName:         "",
			AWSRegion:      "us-east-1",
			BedrockModelID: "amazon.nova-lite-v1:0",
			LogLevel:       "info",
		}

		err := config.validate()
		if err == nil {
			t.Error("Expected error for empty database name")
		}
	})
}
