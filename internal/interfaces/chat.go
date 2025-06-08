package interfaces

import (
	"context"
	"ec-recommend/internal/types"
)

// ChatServiceInterface defines the interface for chat operations
type ChatServiceInterface interface {
	GenerateResponse(ctx context.Context, prompt string) (*types.AIResponse, error)
}
