package service

import (
	"context"
	"ec-recommend/internal/dto"

	"github.com/google/uuid"
)

// RecommendationRepositoryInterface defines the interface for basic repository operations
// This interface is defined in the service package as it is consumed by services
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
