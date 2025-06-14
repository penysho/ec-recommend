package dto

import (
	"time"

	"github.com/google/uuid"
)

// RecommendationRequestV2 represents an enhanced request for product recommendations using RAG and vector search
type RecommendationRequestV2 struct {
	CustomerID         uuid.UUID           `json:"customer_id" binding:"required"`
	RecommendationType string              `json:"recommendation_type,omitempty"`  // "hybrid", "semantic", "collaborative", "vector_search", "knowledge_based"
	ContextType        string              `json:"context_type,omitempty"`         // "homepage", "product_page", "cart", "checkout", "search_results"
	QueryText          string              `json:"query_text,omitempty"`           // Natural language query for semantic search
	ProductID          *uuid.UUID          `json:"product_id,omitempty"`           // For product-based recommendations
	CategoryID         *int                `json:"category_id,omitempty"`          // For category-based recommendations
	PriceRangeMin      *float64            `json:"price_range_min,omitempty"`      // Minimum price filter
	PriceRangeMax      *float64            `json:"price_range_max,omitempty"`      // Maximum price filter
	Limit              int                 `json:"limit,omitempty"`                // Number of recommendations to return (default: 10)
	ExcludeOwned       bool                `json:"exclude_owned,omitempty"`        // Exclude already purchased products
	EnableExplanation  bool                `json:"enable_explanation,omitempty"`   // Include AI-generated explanations
	VectorSearchConfig *VectorSearchConfig `json:"vector_search_config,omitempty"` // Advanced vector search configuration
}

// VectorSearchConfig represents configuration for vector search operations
type VectorSearchConfig struct {
	SimilarityThreshold float64                `json:"similarity_threshold,omitempty"` // Minimum similarity score for results
	SearchType          string                 `json:"search_type,omitempty"`          // "semantic", "hybrid", "keyword"
	RerankerEnabled     bool                   `json:"reranker_enabled,omitempty"`     // Enable result reranking
	MetadataFilters     map[string]interface{} `json:"metadata_filters,omitempty"`     // Additional metadata-based filters
}

// RecommendationResponseV2 represents the enhanced response containing product recommendations with RAG capabilities
type RecommendationResponseV2 struct {
	CustomerID         uuid.UUID                 `json:"customer_id"`
	Recommendations    []ProductRecommendationV2 `json:"recommendations"`
	RecommendationType string                    `json:"recommendation_type"`
	ContextType        string                    `json:"context_type"`
	GeneratedAt        time.Time                 `json:"generated_at"`
	Metadata           RecommendationMetadataV2  `json:"metadata"`
	SemanticInsights   *SemanticInsights         `json:"semantic_insights,omitempty"`
	QueryUnderstanding *QueryUnderstanding       `json:"query_understanding,omitempty"`
}

// ProductRecommendationV2 represents an enhanced product recommendation with AI-powered features
type ProductRecommendationV2 struct {
	ProductID        uuid.UUID          `json:"product_id"`
	Name             string             `json:"name"`
	Description      string             `json:"description,omitempty"`
	Price            float64            `json:"price"`
	OriginalPrice    *float64           `json:"original_price,omitempty"`
	Brand            string             `json:"brand,omitempty"`
	CategoryID       int                `json:"category_id"`
	CategoryName     string             `json:"category_name"`
	RatingAverage    float64            `json:"rating_average"`
	RatingCount      int                `json:"rating_count"`
	PopularityScore  int                `json:"popularity_score"`
	ConfidenceScore  float64            `json:"confidence_score"` // AI-generated confidence
	SimilarityScore  float64            `json:"similarity_score"` // Vector similarity score
	Reason           string             `json:"reason"`           // AI-generated explanation
	Tags             []string           `json:"tags,omitempty"`
	ImageURL         string             `json:"image_url,omitempty"`
	VectorMetadata   *VectorMetadata    `json:"vector_metadata,omitempty"`
	AIInsights       *ProductAIInsights `json:"ai_insights,omitempty"`
	RelevanceContext []RelevanceContext `json:"relevance_context,omitempty"`
}

// VectorMetadata contains metadata about vector search results
type VectorMetadata struct {
	DistanceScore    float64  `json:"distance_score"`
	SearchMethod     string   `json:"search_method"`
	MatchedCriteria  []string `json:"matched_criteria,omitempty"`
	SemanticClusters []string `json:"semantic_clusters,omitempty"`
	EmbeddingModel   string   `json:"embedding_model,omitempty"`
}

// ProductAIInsights contains AI-generated insights about the product
type ProductAIInsights struct {
	KeyFeatures          []string           `json:"key_features,omitempty"`
	UseCases             []string           `json:"use_cases,omitempty"`
	TargetAudience       []string           `json:"target_audience,omitempty"`
	CompetitiveAdvantage string             `json:"competitive_advantage,omitempty"`
	SentimentAnalysis    *SentimentAnalysis `json:"sentiment_analysis,omitempty"`
	TrendAnalysis        *TrendAnalysis     `json:"trend_analysis,omitempty"`
}

// SentimentAnalysis contains sentiment analysis of product reviews
type SentimentAnalysis struct {
	OverallSentiment string   `json:"overall_sentiment"`
	PositiveScore    float64  `json:"positive_score"`
	NegativeScore    float64  `json:"negative_score"`
	NeutralScore     float64  `json:"neutral_score"`
	KeyTopics        []string `json:"key_topics,omitempty"`
}

// TrendAnalysis contains trend analysis for the product
type TrendAnalysis struct {
	TrendDirection  string   `json:"trend_direction"` // "rising", "declining", "stable"
	TrendStrength   float64  `json:"trend_strength"`
	SeasonalPattern string   `json:"seasonal_pattern,omitempty"`
	PeakPeriods     []string `json:"peak_periods,omitempty"`
	MarketPosition  string   `json:"market_position,omitempty"`
}

// RelevanceContext explains why this product is relevant to the customer
type RelevanceContext struct {
	ContextType string  `json:"context_type"` // "purchase_history", "browsing_behavior", "semantic_match", "collaborative"
	Explanation string  `json:"explanation"`
	Confidence  float64 `json:"confidence"`
	SourceData  string  `json:"source_data,omitempty"`
}

// RecommendationMetadataV2 contains enhanced metadata about the recommendation process
type RecommendationMetadataV2 struct {
	AlgorithmVersion   string              `json:"algorithm_version"`
	ProcessingTimeMs   int64               `json:"processing_time_ms"`
	TotalProducts      int                 `json:"total_products"`
	FilteredProducts   int                 `json:"filtered_products"`
	AIModelUsed        string              `json:"ai_model_used,omitempty"`
	EmbeddingModel     string              `json:"embedding_model,omitempty"`
	SessionID          uuid.UUID           `json:"session_id,omitempty"`
	KnowledgeBaseUsed  bool                `json:"knowledge_base_used"`
	VectorSearchUsed   bool                `json:"vector_search_used"`
	SemanticSearchUsed bool                `json:"semantic_search_used"`
	SearchStrategies   []string            `json:"search_strategies,omitempty"`
	PerformanceMetrics *PerformanceMetrics `json:"performance_metrics,omitempty"`
}

// PerformanceMetrics contains performance analytics
type PerformanceMetrics struct {
	VectorSearchTimeMs    int64   `json:"vector_search_time_ms"`
	KnowledgeBaseTimeMs   int64   `json:"knowledge_base_time_ms"`
	AIProcessingTimeMs    int64   `json:"ai_processing_time_ms"`
	CacheHitRate          float64 `json:"cache_hit_rate"`
	TotalTokensUsed       int     `json:"total_tokens_used"`
	EmbeddingGenerationMs int64   `json:"embedding_generation_ms"`
}

// SemanticInsights contains insights from semantic analysis
type SemanticInsights struct {
	QueryIntent       string            `json:"query_intent"`
	ExtractedEntities []ExtractedEntity `json:"extracted_entities,omitempty"`
	SemanticClusters  []SemanticCluster `json:"semantic_clusters,omitempty"`
	RelatedConcepts   []string          `json:"related_concepts,omitempty"`
	ConfidenceLevel   float64           `json:"confidence_level"`
}

// ExtractedEntity represents an entity extracted from the query
type ExtractedEntity struct {
	Type       string  `json:"type"` // "product", "brand", "category", "feature", "price"
	Value      string  `json:"value"`
	Confidence float64 `json:"confidence"`
	StartPos   int     `json:"start_pos,omitempty"`
	EndPos     int     `json:"end_pos,omitempty"`
}

// SemanticCluster represents a group of semantically related items
type SemanticCluster struct {
	ClusterID  string      `json:"cluster_id"`
	Theme      string      `json:"theme"`
	Keywords   []string    `json:"keywords"`
	Relevance  float64     `json:"relevance"`
	ProductIDs []uuid.UUID `json:"product_ids,omitempty"`
}

// QueryUnderstanding contains analysis of the user's query
type QueryUnderstanding struct {
	OriginalQuery   string            `json:"original_query"`
	ProcessedQuery  string            `json:"processed_query"`
	Intent          string            `json:"intent"`
	Entities        []ExtractedEntity `json:"entities,omitempty"`
	Sentiment       string            `json:"sentiment,omitempty"`
	Complexity      string            `json:"complexity"` // "simple", "medium", "complex"
	RequiredContext []string          `json:"required_context,omitempty"`
}

// SemanticSearchRequest represents a request for semantic search
type SemanticSearchRequest struct {
	Query           string     `json:"query" binding:"required"`
	CustomerID      *uuid.UUID `json:"customer_id,omitempty"`
	CategoryID      *int       `json:"category_id,omitempty"`
	PriceRangeMin   *float64   `json:"price_range_min,omitempty"`
	PriceRangeMax   *float64   `json:"price_range_max,omitempty"`
	Limit           int        `json:"limit,omitempty"`
	IncludeMetadata bool       `json:"include_metadata,omitempty"`
}

// SemanticSearchResponse represents the response from semantic search
type SemanticSearchResponse struct {
	Query              string                    `json:"query"`
	Results            []ProductRecommendationV2 `json:"results"`
	TotalFound         int                       `json:"total_found"`
	ProcessingTimeMs   int64                     `json:"processing_time_ms"`
	QueryUnderstanding *QueryUnderstanding       `json:"query_understanding,omitempty"`
	SearchMetadata     *SearchMetadata           `json:"search_metadata,omitempty"`
}

// SearchMetadata contains metadata about the search operation
type SearchMetadata struct {
	SearchType       string                 `json:"search_type"`
	EmbeddingModel   string                 `json:"embedding_model"`
	SimilarityMetric string                 `json:"similarity_metric"`
	FilterApplied    map[string]interface{} `json:"filters_applied,omitempty"`
	RerankerUsed     bool                   `json:"reranker_used"`
	CacheUsed        bool                   `json:"cache_used"`
}

// VectorSimilarityRequest represents a request for vector similarity search
type VectorSimilarityRequest struct {
	ProductID           uuid.UUID `json:"product_id" binding:"required"`
	Limit               int       `json:"limit,omitempty"`
	IncludeMetadata     bool      `json:"include_metadata,omitempty"`
	SimilarityThreshold float64   `json:"similarity_threshold,omitempty"`
}

// VectorSimilarityResponse represents the response from vector similarity search
type VectorSimilarityResponse struct {
	ProductID        uuid.UUID                 `json:"product_id"`
	SimilarProducts  []ProductRecommendationV2 `json:"similar_products"`
	TotalFound       int                       `json:"total_found"`
	ProcessingTimeMs int64                     `json:"processing_time_ms"`
	VectorMetadata   *VectorMetadata           `json:"vector_metadata,omitempty"`
}

// KnowledgeBasedRecommendationRequest represents a request for knowledge-based recommendations
type KnowledgeBasedRecommendationRequest struct {
	CustomerID         uuid.UUID `json:"customer_id" binding:"required"`
	Intent             string    `json:"intent,omitempty"`
	ContextDescription string    `json:"context_description,omitempty"`
	Limit              int       `json:"limit,omitempty"`
}

// KnowledgeBasedRecommendationResponse represents the response from knowledge-based recommendations
type KnowledgeBasedRecommendationResponse struct {
	CustomerID        uuid.UUID                 `json:"customer_id"`
	Intent            string                    `json:"intent,omitempty"`
	Recommendations   []ProductRecommendationV2 `json:"recommendations"`
	ProcessingTimeMs  int64                     `json:"processing_time_ms"`
	KnowledgeInsights *KnowledgeInsights        `json:"knowledge_insights,omitempty"`
	ReasoningChain    []ReasoningStep           `json:"reasoning_chain,omitempty"`
}

// KnowledgeInsights contains insights derived from the knowledge base
type KnowledgeInsights struct {
	RelevantTopics     []string           `json:"relevant_topics"`
	ConfidenceLevel    float64            `json:"confidence_level"`
	SourceDocuments    []string           `json:"source_documents,omitempty"`
	AlternativeQueries []string           `json:"alternative_queries,omitempty"`
	RelatedConcepts    map[string]float64 `json:"related_concepts,omitempty"`
}

// ReasoningStep represents a step in the AI reasoning process
type ReasoningStep struct {
	Step        int     `json:"step"`
	Description string  `json:"description"`
	Confidence  float64 `json:"confidence"`
	Source      string  `json:"source,omitempty"`
}

// RecommendationExplanationResponse represents detailed explanation for a recommendation
type RecommendationExplanationResponse struct {
	RecommendationID   uuid.UUID                 `json:"recommendation_id"`
	CustomerID         uuid.UUID                 `json:"customer_id"`
	ProductID          uuid.UUID                 `json:"product_id"`
	Explanation        string                    `json:"explanation"`
	FactorsConsidered  []ExplanationFactor       `json:"factors_considered"`
	AlternativeOptions []ProductRecommendationV2 `json:"alternative_options,omitempty"`
	GeneratedAt        time.Time                 `json:"generated_at"`
}

// ExplanationFactor represents a factor considered in the recommendation
type ExplanationFactor struct {
	Factor      string  `json:"factor"`
	Weight      float64 `json:"weight"`
	Impact      string  `json:"impact"` // "positive", "negative", "neutral"
	Description string  `json:"description"`
	Confidence  float64 `json:"confidence"`
}

// TrendingProductsRequestV2 represents a request for trending products with AI insights
type TrendingProductsRequestV2 struct {
	CategoryID      *int   `json:"category_id,omitempty"`
	TimeRange       string `json:"time_range,omitempty"` // "daily", "weekly", "monthly"
	Limit           int    `json:"limit,omitempty"`
	IncludeInsights bool   `json:"include_insights,omitempty"`
}

// TrendingProductsResponseV2 represents the response for trending products with AI insights
type TrendingProductsResponseV2 struct {
	TrendingProducts []TrendingProductV2 `json:"trending_products"`
	TimeRange        string              `json:"time_range"`
	TotalFound       int                 `json:"total_found"`
	ProcessingTimeMs int64               `json:"processing_time_ms"`
	TrendInsights    *TrendInsights      `json:"trend_insights,omitempty"`
	MarketAnalysis   *MarketAnalysis     `json:"market_analysis,omitempty"`
}

// TrendingProductV2 represents a trending product with enhanced data
type TrendingProductV2 struct {
	ProductRecommendationV2
	TrendScore         float64                    `json:"trend_score"`
	TrendVelocity      float64                    `json:"trend_velocity"` // Rate of trend change
	TrendDuration      string                     `json:"trend_duration"` // How long it's been trending
	TrendPrediction    *TrendPrediction           `json:"trend_prediction,omitempty"`
	PerformanceMetrics *ProductPerformanceMetrics `json:"performance_metrics,omitempty"`
}

// TrendPrediction contains AI-powered trend predictions
type TrendPrediction struct {
	NextWeekScore     float64 `json:"next_week_score"`
	NextMonthScore    float64 `json:"next_month_score"`
	PeakPredicted     bool    `json:"peak_predicted"`
	ConfidenceLevel   float64 `json:"confidence_level"`
	RecommendedAction string  `json:"recommended_action,omitempty"`
}

// ProductPerformanceMetrics contains product performance analytics
type ProductPerformanceMetrics struct {
	ViewCount            int     `json:"view_count"`
	ClickThroughRate     float64 `json:"click_through_rate"`
	ConversionRate       float64 `json:"conversion_rate"`
	RevenueGrowth        float64 `json:"revenue_growth"`
	CustomerSatisfaction float64 `json:"customer_satisfaction"`
}

// TrendInsights contains overall trend insights
type TrendInsights struct {
	EmergingCategories  []string         `json:"emerging_categories,omitempty"`
	DecliningCategories []string         `json:"declining_categories,omitempty"`
	SeasonalFactors     []string         `json:"seasonal_factors,omitempty"`
	MarketDrivers       []string         `json:"market_drivers,omitempty"`
	PredictedTrends     []PredictedTrend `json:"predicted_trends,omitempty"`
}

// PredictedTrend represents a predicted future trend
type PredictedTrend struct {
	TrendName   string   `json:"trend_name"`
	Probability float64  `json:"probability"`
	TimeFrame   string   `json:"time_frame"`
	Impact      string   `json:"impact"` // "high", "medium", "low"
	Categories  []string `json:"categories,omitempty"`
	Description string   `json:"description"`
}

// MarketAnalysis contains broader market analysis
type MarketAnalysis struct {
	MarketSentiment   string              `json:"market_sentiment"` // "bullish", "bearish", "neutral"
	CompetitorInsight []CompetitorInsight `json:"competitor_insights,omitempty"`
	OpportunityAreas  []string            `json:"opportunity_areas,omitempty"`
	RiskFactors       []string            `json:"risk_factors,omitempty"`
}

// CompetitorInsight contains insights about competitors
type CompetitorInsight struct {
	CompetitorID   string   `json:"competitor_id,omitempty"`
	CompetitorName string   `json:"competitor_name"`
	MarketShare    float64  `json:"market_share,omitempty"`
	TrendDirection string   `json:"trend_direction"` // "gaining", "losing", "stable"
	KeyAdvantages  []string `json:"key_advantages,omitempty"`
}
