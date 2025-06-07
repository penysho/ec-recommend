package handler

import (
	"context"
)

// AIServiceInterface defines the interface for AI operations
type AIServiceInterface interface {
	GenerateResponse(ctx context.Context, prompt string) (*AIResponse, error)
}

// AIResponse represents the response from the AI model
type AIResponse struct {
	Content string `json:"content"`
	Usage   Usage  `json:"usage,omitempty"`
}

// Usage represents token usage information
type Usage struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
}
