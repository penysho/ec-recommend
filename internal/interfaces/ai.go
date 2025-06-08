package interfaces

import (
	"context"
	"ec-recommend/internal/types"
)

// AIServiceInterface defines the interface for AI operations
type AIServiceInterface interface {
	GenerateResponse(ctx context.Context, prompt string) (*types.AIResponse, error)
}
