package interfaces

import (
	"context"
	"ec-recommend/internal/dto"

	"github.com/google/uuid"
)

// RecommendationServiceInterface defines the interface for recommendation operations
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

// RecommendationRepositoryInterface defines the interface for recommendation data access
type RecommendationRepositoryInterface interface {
	// Customer-related methods
	GetCustomerByID(ctx context.Context, customerID uuid.UUID) (*dto.CustomerProfile, error)
	GetCustomerPurchaseHistory(ctx context.Context, customerID uuid.UUID, limit int) ([]dto.PurchaseItem, error)
	GetCustomerActivities(ctx context.Context, customerID uuid.UUID, limit int) ([]dto.ActivityItem, error)

	// Product-related methods
	GetProductsByIDs(ctx context.Context, productIDs []uuid.UUID) ([]dto.ProductRecommendation, error)
	GetProductsByCategory(ctx context.Context, categoryID int, limit int) ([]dto.ProductRecommendation, error)
	GetTrendingProducts(ctx context.Context, categoryID *int, limit int) ([]dto.ProductRecommendation, error)
	GetSimilarProductsByTags(ctx context.Context, tags []string, excludeProductID uuid.UUID, limit int) ([]dto.ProductRecommendation, error)
	GetProductsInPriceRange(ctx context.Context, minPrice, maxPrice float64, limit int) ([]dto.ProductRecommendation, error)

	// Analytics methods
	LogRecommendation(ctx context.Context, customerID uuid.UUID, recommendationType, contextType string, productIDs []uuid.UUID, sessionID uuid.UUID) error
	LogRecommendationInteraction(ctx context.Context, analytics *dto.RecommendationAnalytics) error

	// Collaborative filtering methods
	GetCustomersWithSimilarPurchases(ctx context.Context, customerID uuid.UUID, limit int) ([]uuid.UUID, error)
	GetPopularProductsAmongSimilarCustomers(ctx context.Context, similarCustomerIDs []uuid.UUID, excludeOwned []uuid.UUID, limit int) ([]dto.ProductRecommendation, error)
}
