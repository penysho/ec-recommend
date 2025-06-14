package service

import (
	"context"
	"ec-recommend/internal/dto"
	"ec-recommend/internal/types"
)

// ChatServiceInterface defines the interface for chat service operations
// This interface is defined in the service package as it is consumed by services
type ChatServiceInterface interface {
	// Chat performs a chat interaction using the AI model
	Chat(ctx context.Context, req *dto.ChatRequest) (*dto.ChatResponse, error)

	// GenerateResponse generates a response for the given prompt
	GenerateResponse(ctx context.Context, prompt string) (*types.AIResponse, error)
}
