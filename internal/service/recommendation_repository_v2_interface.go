package service

import (
	"context"
	"ec-recommend/internal/dto"

	"github.com/google/uuid"
)

// RecommendationRepositoryV2Interface defines the interface for repository operations
// This interface is defined in the service package as it is consumed by services
type RecommendationRepositoryV2Interface interface {
	// Enhanced customer and product methods
	GetCustomerByID(ctx context.Context, customerID uuid.UUID) (*dto.CustomerProfile, error)
	GetCustomerPurchaseHistory(ctx context.Context, customerID uuid.UUID, limit int) ([]dto.PurchaseItem, error)
	GetCustomerActivities(ctx context.Context, customerID uuid.UUID, limit int) ([]dto.ActivityItem, error)
	GetProductsByIDs(ctx context.Context, productIDs []uuid.UUID) ([]dto.ProductRecommendationV2, error)
	GetProductsByCategory(ctx context.Context, categoryID int, limit int) ([]dto.ProductRecommendationV2, error)

	// Vector and semantic search methods
	GetProductsWithVectorSearch(ctx context.Context, vector []float64, limit int, filters map[string]interface{}) ([]dto.ProductRecommendationV2, error)
	GetProductsWithSemanticSearch(ctx context.Context, query string, limit int, filters map[string]interface{}) ([]dto.ProductRecommendationV2, error)
	GetProductsWithHybridSearch(ctx context.Context, query string, vector []float64, limit int, filters map[string]interface{}) ([]dto.ProductRecommendationV2, error)

	// Enhanced analytics and trending
	GetTrendingProductsV2(ctx context.Context, categoryID *int, timeRange string, limit int) ([]dto.TrendingProductV2, error)
	GetProductPerformanceMetrics(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID]*dto.ProductPerformanceMetrics, error)
	GetMarketAnalysis(ctx context.Context, categoryID *int, timeRange string) (*dto.MarketAnalysis, error)

	// Logging and analytics
	LogRecommendation(ctx context.Context, customerID uuid.UUID, recommendationType, contextType string, productIDs []uuid.UUID, sessionID uuid.UUID) error
	LogRecommendationInteraction(ctx context.Context, analytics *dto.RecommendationAnalytics) error
	LogSemanticSearch(ctx context.Context, customerID *uuid.UUID, query string, results []uuid.UUID, processingTimeMs int64) error

	// Cache management
	GetCachedRecommendations(ctx context.Context, key string) ([]dto.ProductRecommendationV2, error)
	SetCachedRecommendations(ctx context.Context, key string, recommendations []dto.ProductRecommendationV2, ttl int64) error
	InvalidateCache(ctx context.Context, pattern string) error
}
