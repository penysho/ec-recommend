package handler

import (
	"context"
	"ec-recommend/internal/dto"

	"github.com/google/uuid"
)

// RecommendationServiceV2Interface defines the interface for advanced recommendation service operations
// This interface is defined in the handler package as it is consumed by handlers
type RecommendationServiceV2Interface interface {
	// GetRecommendationsV2 generates advanced product recommendations using RAG and vector search
	GetRecommendationsV2(ctx context.Context, req *dto.RecommendationRequestV2) (*dto.RecommendationResponseV2, error)

	// SemanticSearch performs semantic search using natural language queries
	SemanticSearch(ctx context.Context, req *dto.SemanticSearchRequest) (*dto.SemanticSearchResponse, error)

	// GetVectorSimilarProducts finds products similar to a given product using vector similarity
	GetVectorSimilarProducts(ctx context.Context, req *dto.VectorSimilarityRequest) (*dto.VectorSimilarityResponse, error)

	// GetKnowledgeBasedRecommendations generates recommendations based on comprehensive knowledge base analysis
	GetKnowledgeBasedRecommendations(ctx context.Context, req *dto.KnowledgeBasedRecommendationRequest) (*dto.KnowledgeBasedRecommendationResponse, error)

	// GetRecommendationExplanation provides detailed explanation for a specific recommendation
	GetRecommendationExplanation(ctx context.Context, recommendationID, customerID uuid.UUID) (*dto.RecommendationExplanationResponse, error)

	// GetTrendingProductsV2 returns trending products with AI-powered insights
	GetTrendingProductsV2(ctx context.Context, req *dto.TrendingProductsRequestV2) (*dto.TrendingProductsResponseV2, error)

	// GetCustomerProfile retrieves customer profile data for personalized recommendations
	GetCustomerProfile(ctx context.Context, customerID uuid.UUID) (*dto.CustomerProfile, error)

	// LogRecommendationInteraction logs customer interactions with recommendations for analytics
	LogRecommendationInteraction(ctx context.Context, analytics *dto.RecommendationAnalytics) error
}
