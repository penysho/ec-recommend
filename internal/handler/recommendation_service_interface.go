package handler

import (
	"context"
	"ec-recommend/internal/dto"

	"github.com/google/uuid"
)

// RecommendationServiceInterface defines the interface for basic recommendation service operations
// This interface is defined in the handler package as it is consumed by handlers
type RecommendationServiceInterface interface {
	// GetRecommendations generates product recommendations for a customer
	GetRecommendations(ctx context.Context, req *dto.RecommendationRequest) (*dto.RecommendationResponse, error)

	// GetCustomerProfile retrieves customer profile data for recommendations
	GetCustomerProfile(ctx context.Context, customerID uuid.UUID) (*dto.CustomerProfile, error)

	// LogRecommendationInteraction logs customer interactions with recommendations
	LogRecommendationInteraction(ctx context.Context, analytics *dto.RecommendationAnalytics) error

	// GetSimilarProducts finds products similar to a given product
	GetSimilarProducts(ctx context.Context, productID uuid.UUID, limit int) ([]dto.ProductRecommendation, error)

	// GetTrendingProducts returns currently trending products
	GetTrendingProducts(ctx context.Context, categoryID *int, limit int) ([]dto.ProductRecommendation, error)

	// GetPersonalizedRecommendations generates AI-powered personalized recommendations
	GetPersonalizedRecommendations(ctx context.Context, profile *dto.CustomerProfile, limit int) ([]dto.ProductRecommendation, error)
}
