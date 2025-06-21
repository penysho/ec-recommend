package service

import (
	"context"
	"ec-recommend/internal/dto"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RecommendationServiceV2 implements the RecommendationServiceV2Interface with RAG and vector search capabilities
type RecommendationServiceV2 struct {
	repo             RecommendationRepositoryV2Interface
	rag              RAGInterface
	chatService      ChatServiceInterface
	modelID          string
	knowledgeBaseID  string
	embeddingModelID string
}

// NewRecommendationServiceV2 creates a new enhanced recommendation service instance
func NewRecommendationServiceV2(
	repo RecommendationRepositoryV2Interface,
	rag RAGInterface,
	chatService ChatServiceInterface,
	modelID, knowledgeBaseID, embeddingModelID string,
) *RecommendationServiceV2 {
	return &RecommendationServiceV2{
		repo:             repo,
		rag:              rag,
		chatService:      chatService,
		modelID:          modelID,
		knowledgeBaseID:  knowledgeBaseID,
		embeddingModelID: embeddingModelID,
	}
}

// GetRecommendationsV2 generates advanced product recommendations using RAG and vector search
func (rs *RecommendationServiceV2) GetRecommendationsV2(ctx context.Context, req *dto.RecommendationRequestV2) (*dto.RecommendationResponseV2, error) {
	startTime := time.Now()

	// Set default values
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.RecommendationType == "" {
		req.RecommendationType = "hybrid"
	}
	if req.ContextType == "" {
		req.ContextType = "homepage"
	}

	// Get customer profile
	profile, err := rs.GetCustomerProfile(ctx, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer profile: %w", err)
	}

	var recommendations []dto.ProductRecommendationV2
	var semanticInsights *dto.SemanticInsights
	var queryUnderstanding *dto.QueryUnderstanding
	var searchStrategies []string
	var performanceMetrics = &dto.PerformanceMetrics{}

	// Generate recommendations based on type
	switch req.RecommendationType {
	case "semantic", "vector_search":
		// Both semantic and vector_search now use the same underlying semantic search with automatic vectorization
		// semantic: uses query_text directly
		// vector_search: generates query from product_id and uses semantic search for similarity
		recommendations, semanticInsights, queryUnderstanding, err = rs.generateSemanticRecommendations(ctx, req, profile, performanceMetrics)
		if req.RecommendationType == "semantic" {
			searchStrategies = append(searchStrategies, "semantic_search")
		} else {
			searchStrategies = append(searchStrategies, "vector_similarity")
		}
	case "knowledge_based":
		recommendations, err = rs.generateKnowledgeBasedRecommendations(ctx, req, profile, performanceMetrics)
		searchStrategies = append(searchStrategies, "knowledge_base_rag")
	case "collaborative":
		recommendations, err = rs.generateCollaborativeRecommendations(ctx, req, profile)
		searchStrategies = append(searchStrategies, "collaborative_filtering")
	case "hybrid":
		recommendations, semanticInsights, queryUnderstanding, err = rs.generateHybridRecommendations(ctx, req, profile, performanceMetrics)
		searchStrategies = append(searchStrategies, "hybrid_rag", "semantic_search", "collaborative_filtering")
	default:
		return nil, fmt.Errorf("unsupported recommendation type: %s", req.RecommendationType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to generate recommendations: %w", err)
	}

	// Filter out owned products if requested
	if req.ExcludeOwned {
		recommendations = rs.filterOwnedProductsV2(recommendations, profile.PurchaseHistory)
	}

	// Apply price range filters
	if req.PriceRangeMin != nil || req.PriceRangeMax != nil {
		recommendations = rs.filterByPriceRange(recommendations, req.PriceRangeMin, req.PriceRangeMax)
	}

	// Limit results
	if len(recommendations) > req.Limit {
		recommendations = recommendations[:req.Limit]
	}

	// Generate AI-powered explanations if requested
	if req.EnableExplanation {
		recommendations, err = rs.enhanceWithAIExplanations(ctx, recommendations, profile, req.ContextType, performanceMetrics)
		if err != nil {
			log.Printf("Warning: failed to enhance recommendations with AI explanations: %v", err)
		}
	}

	// Log recommendation for analytics
	sessionID := uuid.New()
	productIDs := make([]uuid.UUID, len(recommendations))
	for i, rec := range recommendations {
		productIDs[i] = rec.ProductID
	}

	err = rs.repo.LogRecommendation(ctx, req.CustomerID, req.RecommendationType, req.ContextType, productIDs, sessionID)
	if err != nil {
		log.Printf("Warning: failed to log recommendation: %v", err)
	}

	processingTime := time.Since(startTime).Milliseconds()
	performanceMetrics.AIProcessingTimeMs = processingTime

	return &dto.RecommendationResponseV2{
		CustomerID:         req.CustomerID,
		Recommendations:    recommendations,
		RecommendationType: req.RecommendationType,
		ContextType:        req.ContextType,
		GeneratedAt:        time.Now(),
		SemanticInsights:   semanticInsights,
		QueryUnderstanding: queryUnderstanding,
		Metadata: dto.RecommendationMetadataV2{
			AlgorithmVersion:   "rag_hybrid_v2.0",
			ProcessingTimeMs:   processingTime,
			TotalProducts:      len(recommendations),
			FilteredProducts:   len(recommendations),
			AIModelUsed:        rs.modelID,
			EmbeddingModel:     rs.embeddingModelID,
			SessionID:          sessionID,
			KnowledgeBaseUsed:  contains(searchStrategies, "knowledge_base_rag"),
			VectorSearchUsed:   contains(searchStrategies, "vector_similarity"),
			SemanticSearchUsed: contains(searchStrategies, "semantic_search"),
			SearchStrategies:   searchStrategies,
			PerformanceMetrics: performanceMetrics,
		},
	}, nil
}

// SemanticSearch performs semantic search using natural language queries
func (rs *RecommendationServiceV2) SemanticSearch(ctx context.Context, req *dto.SemanticSearchRequest) (*dto.SemanticSearchResponse, error) {
	startTime := time.Now()

	// Generate query understanding
	queryUnderstanding, err := rs.analyzeQuery(ctx, req.Query)
	if err != nil {
		log.Printf("Warning: failed to analyze query: %v", err)
	}

	// Build filters
	filters := make(map[string]interface{})
	if req.CategoryID != nil {
		filters["category_id"] = *req.CategoryID
	}
	if req.PriceRangeMin != nil {
		filters["price_min"] = *req.PriceRangeMin
	}
	if req.PriceRangeMax != nil {
		filters["price_max"] = *req.PriceRangeMax
	}

	// Perform semantic search using RAG Knowledge Base
	ragResponse, err := rs.rag.GetProductsWithSemanticSearch(ctx, req.Query, req.Limit, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to perform semantic search: %w", err)
	}

	// Convert RAG results to ProductRecommendationV2
	results, err := rs.convertRAGResultsToProducts(ctx, ragResponse.Results, "semantic_search")
	if err != nil {
		return nil, fmt.Errorf("failed to convert semantic search results: %w", err)
	}

	// Log semantic search
	resultIDs := make([]uuid.UUID, len(results))
	for i, result := range results {
		resultIDs[i] = result.ProductID
	}

	processingTime := time.Since(startTime).Milliseconds()
	err = rs.repo.LogSemanticSearch(ctx, req.CustomerID, req.Query, resultIDs, processingTime)
	if err != nil {
		log.Printf("Warning: failed to log semantic search: %v", err)
	}

	return &dto.SemanticSearchResponse{
		Query:              req.Query,
		Results:            results,
		TotalFound:         len(results),
		ProcessingTimeMs:   processingTime,
		QueryUnderstanding: queryUnderstanding,
		SearchMetadata: &dto.SearchMetadata{
			SearchType:       "semantic",
			EmbeddingModel:   rs.embeddingModelID,
			SimilarityMetric: "cosine",
			FilterApplied:    filters,
			RerankerUsed:     false,
			CacheUsed:        false,
		},
	}, nil
}

// GetVectorSimilarProducts finds products similar to a given product using semantic similarity
// This method now leverages AWS Bedrock Knowledge Base's automatic vectorization capabilities
// instead of manually generating and passing vector embeddings
func (rs *RecommendationServiceV2) GetVectorSimilarProducts(ctx context.Context, req *dto.VectorSimilarityRequest) (*dto.VectorSimilarityResponse, error) {
	startTime := time.Now()

	// Get target product details
	products, err := rs.repo.GetProductsByIDs(ctx, []uuid.UUID{req.ProductID})
	if err != nil || len(products) == 0 {
		return nil, fmt.Errorf("failed to get target product: %w", err)
	}

	targetProduct := products[0]

	// Build comprehensive product description for semantic search
	// The Knowledge Base will automatically convert this to vectors internally
	comprehensiveProductText := rs.buildComprehensiveProductText(targetProduct)

	// Create similarity query using product characteristics
	similarityQuery := fmt.Sprintf("Find products similar to: %s", comprehensiveProductText)

	log.Printf("Using semantic search for product similarity - Target: %s, Query: %s",
		targetProduct.ProductID, similarityQuery)

	// Set up filters to exclude the target product
	filters := make(map[string]interface{})
	filters["exclude_id"] = req.ProductID.String()

	// Use enhanced semantic search which automatically handles vectorization
	ragResponse, err := rs.rag.GetProductsWithSemanticSearch(ctx, similarityQuery, req.Limit, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to perform semantic similarity search: %w", err)
	}

	// Convert RAG results to ProductRecommendationV2
	similarProducts, err := rs.convertRAGResultsToProducts(ctx, ragResponse.Results, "semantic_similarity")
	if err != nil {
		return nil, fmt.Errorf("failed to convert semantic search results: %w", err)
	}

	// Enhance with semantic similarity metadata
	for i := range similarProducts {
		if similarProducts[i].VectorMetadata == nil {
			similarProducts[i].VectorMetadata = &dto.VectorMetadata{}
		}
		similarProducts[i].VectorMetadata.SearchMethod = "semantic_similarity"
		similarProducts[i].VectorMetadata.EmbeddingModel = rs.embeddingModelID
		if similarProducts[i].Reason == "" {
			similarProducts[i].Reason = fmt.Sprintf("Semantically similar to: %s", targetProduct.Name)
		}

		// Add additional similarity context
		if len(ragResponse.Results) > i {
			similarProducts[i].SimilarityScore = ragResponse.Results[i].SimilarityScore
			similarProducts[i].ConfidenceScore = ragResponse.Results[i].ConfidenceScore
		}
	}

	processingTime := time.Since(startTime).Milliseconds()

	return &dto.VectorSimilarityResponse{
		ProductID:        req.ProductID,
		SimilarProducts:  similarProducts,
		TotalFound:       len(similarProducts),
		ProcessingTimeMs: processingTime,
		VectorMetadata: &dto.VectorMetadata{
			SearchMethod:   "semantic_similarity_enhanced",
			EmbeddingModel: rs.embeddingModelID,
		},
	}, nil
}

// GetKnowledgeBasedRecommendations generates recommendations based on comprehensive knowledge base analysis
func (rs *RecommendationServiceV2) GetKnowledgeBasedRecommendations(ctx context.Context, req *dto.KnowledgeBasedRecommendationRequest) (*dto.KnowledgeBasedRecommendationResponse, error) {
	startTime := time.Now()

	// Get customer profile
	profile, err := rs.GetCustomerProfile(ctx, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer profile: %w", err)
	}

	// Create knowledge-based query
	query := rs.createKnowledgeBasedQuery(profile, req.Intent, req.ContextDescription)

	// Query knowledge base
	kbResponse, err := rs.rag.QueryKnowledgeBase(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query knowledge base: %w", err)
	}

	// Extract recommendations from knowledge base results
	recommendations, reasoningChain, err := rs.extractRecommendationsFromKB(ctx, kbResponse, profile, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to extract recommendations from knowledge base: %w", err)
	}

	// Create knowledge insights
	knowledgeInsights := rs.createKnowledgeInsights(kbResponse)

	processingTime := time.Since(startTime).Milliseconds()

	return &dto.KnowledgeBasedRecommendationResponse{
		CustomerID:        req.CustomerID,
		Intent:            req.Intent,
		Recommendations:   recommendations,
		ProcessingTimeMs:  processingTime,
		KnowledgeInsights: knowledgeInsights,
		ReasoningChain:    reasoningChain,
	}, nil
}

// GetRecommendationExplanation provides detailed explanation for a specific recommendation
func (rs *RecommendationServiceV2) GetRecommendationExplanation(ctx context.Context, recommendationID, customerID uuid.UUID) (*dto.RecommendationExplanationResponse, error) {
	// Get customer profile
	profile, err := rs.GetCustomerProfile(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer profile: %w", err)
	}

	// In a real implementation, you would retrieve the recommendation details from a database
	// For now, we'll create a comprehensive explanation based on the customer profile

	// Create AI prompt for explanation
	prompt := rs.createExplanationPrompt(recommendationID, profile)

	// Get AI-generated explanation
	chatResponse, err := rs.chatService.GenerateResponse(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate explanation: %w", err)
	}

	// Parse explanation factors
	factors := rs.parseExplanationFactors(chatResponse.Content, profile)

	// Get alternative options (simplified for this example)
	alternatives, err := rs.repo.GetProductsByCategory(ctx, profile.PreferredCategories[0], 3)
	if err != nil {
		log.Printf("Warning: failed to get alternative products: %v", err)
		alternatives = []dto.ProductRecommendationV2{}
	}

	return &dto.RecommendationExplanationResponse{
		RecommendationID:   recommendationID,
		CustomerID:         customerID,
		ProductID:          recommendationID, // Simplified assumption
		Explanation:        chatResponse.Content,
		FactorsConsidered:  factors,
		AlternativeOptions: alternatives,
		GeneratedAt:        time.Now(),
	}, nil
}

// GetTrendingProductsV2 returns trending products with AI-powered insights
func (rs *RecommendationServiceV2) GetTrendingProductsV2(ctx context.Context, req *dto.TrendingProductsRequestV2) (*dto.TrendingProductsResponseV2, error) {
	startTime := time.Now()

	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.TimeRange == "" {
		req.TimeRange = "weekly"
	}

	// Get trending products
	trendingProducts, err := rs.repo.GetTrendingProductsV2(ctx, req.CategoryID, req.TimeRange, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending products: %w", err)
	}

	var trendInsights *dto.TrendInsights
	var marketAnalysis *dto.MarketAnalysis

	// Generate AI insights if requested
	if req.IncludeInsights {
		trendInsights, err = rs.generateTrendInsights(ctx, trendingProducts, req.TimeRange)
		if err != nil {
			log.Printf("Warning: failed to generate trend insights: %v", err)
		}

		marketAnalysis, err = rs.repo.GetMarketAnalysis(ctx, req.CategoryID, req.TimeRange)
		if err != nil {
			log.Printf("Warning: failed to get market analysis: %v", err)
		}
	}

	processingTime := time.Since(startTime).Milliseconds()

	return &dto.TrendingProductsResponseV2{
		TrendingProducts: trendingProducts,
		TimeRange:        req.TimeRange,
		TotalFound:       len(trendingProducts),
		ProcessingTimeMs: processingTime,
		TrendInsights:    trendInsights,
		MarketAnalysis:   marketAnalysis,
	}, nil
}

// GetCustomerProfile retrieves customer profile data for personalized recommendations
func (rs *RecommendationServiceV2) GetCustomerProfile(ctx context.Context, customerID uuid.UUID) (*dto.CustomerProfile, error) {
	profile, err := rs.repo.GetCustomerByID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	// Get purchase history
	purchaseHistory, err := rs.repo.GetCustomerPurchaseHistory(ctx, customerID, 50)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase history: %w", err)
	}
	profile.PurchaseHistory = purchaseHistory

	// Get recent activities
	activities, err := rs.repo.GetCustomerActivities(ctx, customerID, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get activities: %w", err)
	}
	profile.RecentActivities = activities

	return profile, nil
}

// LogRecommendationInteraction logs customer interactions with recommendations for analytics
func (rs *RecommendationServiceV2) LogRecommendationInteraction(ctx context.Context, analytics *dto.RecommendationAnalytics) error {
	return rs.repo.LogRecommendationInteraction(ctx, analytics)
}

// Private helper methods

func (rs *RecommendationServiceV2) generateSemanticRecommendations(ctx context.Context, req *dto.RecommendationRequestV2, profile *dto.CustomerProfile, metrics *dto.PerformanceMetrics) ([]dto.ProductRecommendationV2, *dto.SemanticInsights, *dto.QueryUnderstanding, error) {
	startTime := time.Now()

	var queryText string
	var searchMethod string
	var queryUnderstanding *dto.QueryUnderstanding

	// Determine query text and search method based on request type
	if req.RecommendationType == "vector_search" {
		// For vector_search, validate product_id and generate query from product
		if req.ProductID == nil {
			return nil, nil, nil, fmt.Errorf("product_id is required for vector search recommendations")
		}

		// Get target product and generate similarity query
		products, err := rs.repo.GetProductsByIDs(ctx, []uuid.UUID{*req.ProductID})
		if err != nil || len(products) == 0 {
			return nil, nil, nil, fmt.Errorf("failed to get target product: %w", err)
		}

		targetProduct := products[0]
		comprehensiveProductText := rs.buildComprehensiveProductText(targetProduct)
		queryText = fmt.Sprintf("Find products similar to: %s", comprehensiveProductText)
		searchMethod = "semantic_similarity"

		log.Printf("Using semantic search for vector recommendations - Target: %s, Query: %s",
			targetProduct.ProductID, queryText)
	} else {
		// For semantic search, validate query_text
		if req.QueryText == "" {
			return nil, nil, nil, fmt.Errorf("query_text is required for semantic recommendations")
		}

		queryText = req.QueryText
		searchMethod = "semantic_search"

		// Analyze query for better understanding
		var err error
		queryUnderstanding, err = rs.analyzeQuery(ctx, req.QueryText)
		if err != nil {
			log.Printf("Warning: failed to analyze query: %v", err)
		}
	}

	// Build personalized filters based on customer profile
	filters := rs.buildPersonalizedFilters(profile, req)

	// For vector_search, exclude the target product
	if req.RecommendationType == "vector_search" && req.ProductID != nil {
		filters["exclude_id"] = req.ProductID.String()
	}

	// Perform semantic search using RAG Knowledge Base (handles both semantic and vector similarity)
	ragResponse, err := rs.rag.GetProductsWithSemanticSearch(ctx, queryText, req.Limit*2, filters)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to perform semantic search: %w", err)
	}

	// Convert RAG results to ProductRecommendationV2
	results, err := rs.convertRAGResultsToProducts(ctx, ragResponse.Results, searchMethod)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to convert semantic search results: %w", err)
	}

	// Generate semantic insights (only for semantic search with actual query)
	var semanticInsights *dto.SemanticInsights
	if req.RecommendationType == "semantic" && queryUnderstanding != nil {
		semanticInsights = rs.generateSemanticInsights(results, queryUnderstanding)
	}

	metrics.VectorSearchTimeMs = time.Since(startTime).Milliseconds()

	return results, semanticInsights, queryUnderstanding, nil
}

func (rs *RecommendationServiceV2) generateKnowledgeBasedRecommendations(ctx context.Context, req *dto.RecommendationRequestV2, profile *dto.CustomerProfile, metrics *dto.PerformanceMetrics) ([]dto.ProductRecommendationV2, error) {
	startTime := time.Now()

	// Create comprehensive query based on customer profile
	query := rs.createComprehensiveQuery(profile, req.ContextType)

	// Query knowledge base
	kbResponse, err := rs.rag.QueryKnowledgeBase(ctx, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query knowledge base: %w", err)
	}

	// Extract product recommendations from knowledge base response
	recommendations, _, err := rs.extractRecommendationsFromKB(ctx, kbResponse, profile, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to extract recommendations: %w", err)
	}

	metrics.KnowledgeBaseTimeMs = time.Since(startTime).Milliseconds()

	return recommendations, nil
}

func (rs *RecommendationServiceV2) generateCollaborativeRecommendations(ctx context.Context, req *dto.RecommendationRequestV2, profile *dto.CustomerProfile) ([]dto.ProductRecommendationV2, error) {
	// This would integrate with the existing collaborative filtering logic
	// For now, we'll use a simplified approach

	// Get products from preferred categories
	var allRecommendations []dto.ProductRecommendationV2

	for _, categoryID := range profile.PreferredCategories {
		categoryProducts, err := rs.repo.GetProductsByCategory(ctx, categoryID, req.Limit/len(profile.PreferredCategories)+1)
		if err != nil {
			continue
		}
		allRecommendations = append(allRecommendations, categoryProducts...)
	}

	// Remove duplicates and limit results
	uniqueProducts := rs.removeDuplicateProductsV2(allRecommendations)
	if len(uniqueProducts) > req.Limit {
		uniqueProducts = uniqueProducts[:req.Limit]
	}

	return uniqueProducts, nil
}

func (rs *RecommendationServiceV2) generateHybridRecommendations(ctx context.Context, req *dto.RecommendationRequestV2, profile *dto.CustomerProfile, metrics *dto.PerformanceMetrics) ([]dto.ProductRecommendationV2, *dto.SemanticInsights, *dto.QueryUnderstanding, error) {
	var allRecommendations []dto.ProductRecommendationV2
	var semanticInsights *dto.SemanticInsights
	var queryUnderstanding *dto.QueryUnderstanding
	var err error

	// 1. Semantic search (40% weight) if query provided
	if req.QueryText != "" {
		// Create a copy of the request for semantic search
		semanticReq := *req
		semanticReq.RecommendationType = "semantic"
		semanticRecs, insights, understanding, err := rs.generateSemanticRecommendations(ctx, &semanticReq, profile, metrics)
		if err == nil {
			for i := range semanticRecs {
				semanticRecs[i].ConfidenceScore = semanticRecs[i].ConfidenceScore * 0.4
			}
			allRecommendations = append(allRecommendations, semanticRecs...)
			semanticInsights = insights
			queryUnderstanding = understanding
		}
	}

	// 2. Vector search (30% weight) if product ID provided
	if req.ProductID != nil {
		// Create a copy of the request for vector search
		vectorReq := *req
		vectorReq.RecommendationType = "vector_search"
		vectorRecs, _, _, err := rs.generateSemanticRecommendations(ctx, &vectorReq, profile, metrics)
		if err == nil {
			for i := range vectorRecs {
				vectorRecs[i].ConfidenceScore = vectorRecs[i].ConfidenceScore * 0.3
			}
			allRecommendations = append(allRecommendations, vectorRecs...)
		}
	}

	// 3. Knowledge-based recommendations (20% weight)
	kbRecs, err := rs.generateKnowledgeBasedRecommendations(ctx, req, profile, metrics)
	if err == nil {
		for i := range kbRecs {
			kbRecs[i].ConfidenceScore = kbRecs[i].ConfidenceScore * 0.2
		}
		allRecommendations = append(allRecommendations, kbRecs...)
	}

	// 4. Collaborative recommendations (10% weight)
	collabRecs, err := rs.generateCollaborativeRecommendations(ctx, req, profile)
	if err == nil {
		for i := range collabRecs {
			collabRecs[i].ConfidenceScore = collabRecs[i].ConfidenceScore * 0.1
		}
		allRecommendations = append(allRecommendations, collabRecs...)
	}

	// Remove duplicates and sort by confidence score
	uniqueProducts := rs.removeDuplicateProductsV2(allRecommendations)
	sort.Slice(uniqueProducts, func(i, j int) bool {
		return uniqueProducts[i].ConfidenceScore > uniqueProducts[j].ConfidenceScore
	})

	if len(uniqueProducts) > req.Limit {
		uniqueProducts = uniqueProducts[:req.Limit]
	}

	return uniqueProducts, semanticInsights, queryUnderstanding, nil
}

func (rs *RecommendationServiceV2) analyzeQuery(ctx context.Context, query string) (*dto.QueryUnderstanding, error) {
	// Create AI prompt for query analysis
	prompt := fmt.Sprintf(`
Analyze the following user query for e-commerce product search:

Query: "%s"

Please provide a JSON response with the following structure:
{
  "intent": "string (e.g., 'product_search', 'comparison', 'gift_suggestion')",
  "entities": [
    {
      "type": "string (e.g., 'product', 'brand', 'category', 'feature', 'price')",
      "value": "string",
      "confidence": float
    }
  ],
  "sentiment": "string (positive, negative, neutral)",
  "complexity": "string (simple, medium, complex)",
  "required_context": ["string array of required context"]
}
`, query)

	chatResponse, err := rs.chatService.GenerateResponse(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var understanding dto.QueryUnderstanding
	err = json.Unmarshal([]byte(chatResponse.Content), &understanding)
	if err != nil {
		// Fallback to basic analysis
		understanding = dto.QueryUnderstanding{
			OriginalQuery:  query,
			ProcessedQuery: strings.ToLower(query),
			Intent:         "product_search",
			Sentiment:      "neutral",
			Complexity:     "simple",
		}
	}

	understanding.OriginalQuery = query
	understanding.ProcessedQuery = strings.ToLower(query)

	return &understanding, nil
}

func (rs *RecommendationServiceV2) buildPersonalizedFilters(profile *dto.CustomerProfile, req *dto.RecommendationRequestV2) map[string]interface{} {
	filters := make(map[string]interface{})

	// Apply category filters - enhanced with multiple categories support
	if req.CategoryID != nil {
		filters["category_id"] = *req.CategoryID
	} else if len(profile.PreferredCategories) > 0 {
		filters["category_ids"] = profile.PreferredCategories
	}

	// Apply price range filters with enhanced logic
	if req.PriceRangeMin != nil {
		filters["price_min"] = *req.PriceRangeMin
	} else if profile.PriceRangeMin != nil {
		filters["price_min"] = *profile.PriceRangeMin
	}

	if req.PriceRangeMax != nil {
		filters["price_max"] = *req.PriceRangeMax
	} else if profile.PriceRangeMax != nil {
		filters["price_max"] = *profile.PriceRangeMax
	}

	// Apply preferred brands filter
	if len(profile.PreferredBrands) > 0 {
		filters["preferred_brands"] = profile.PreferredBrands
	}

	// Enhanced metadata filters based on customer profile and CSV batch metadata
	rs.applyEnhancedMetadataFilters(filters, profile, req)

	return filters
}

// applyEnhancedMetadataFilters applies advanced metadata filtering based on CSV batch metadata structure
func (rs *RecommendationServiceV2) applyEnhancedMetadataFilters(filters map[string]interface{}, profile *dto.CustomerProfile, req *dto.RecommendationRequestV2) {
	// Rating-based filtering for quality-conscious customers
	if profile.IsPremium || profile.OrderCount > 5 {
		filters["rating_min"] = 4.0      // Premium customers prefer high-rated products
		filters["rating_count_min"] = 10 // Ensure sufficient review volume
	} else {
		filters["rating_min"] = 3.0 // Standard minimum rating
	}

	// Stock status filtering - only show available products
	filters["stock_status"] = []string{"in_stock", "limited_stock"}
	filters["is_active"] = true

	// Popularity-based filtering based on customer engagement level
	if profile.OrderCount > 10 {
		// High-engagement customers might prefer trending or niche products
		filters["popularity_score_min"] = 0 // Allow all popularity levels
	} else {
		// New customers prefer popular, well-established products
		filters["popularity_score_min"] = 50
	}

	// Price-based filtering with smart defaults based on customer behavior
	if profile.TotalSpent > 1000 && profile.OrderCount > 3 {
		// High-value customers - adjust price expectations
		avgOrderValue := profile.TotalSpent / float64(profile.OrderCount)
		if _, exists := filters["price_min"]; !exists {
			filters["price_min"] = avgOrderValue * 0.5 // 50% of avg order value as minimum
		}
		if _, exists := filters["price_max"]; !exists {
			filters["price_max"] = avgOrderValue * 2.0 // 200% of avg order value as maximum
		}
	}

	// Tag-based filtering using customer lifestyle and preferences
	if len(profile.LifestyleTags) > 0 {
		filters["lifestyle_tags"] = profile.LifestyleTags
	}

	// Seasonal filtering
	seasonalTags := rs.getCurrentSeasonalTags()
	if len(seasonalTags) > 0 {
		filters["seasonal_boost"] = seasonalTags
	}

	// Exclude recently purchased categories if customer seeks diversity
	if req.ExcludeOwned && len(profile.PurchaseHistory) > 0 {
		recentCategoryIDs := rs.extractRecentCategoryIDs(profile.PurchaseHistory, 30) // Last 30 days
		if len(recentCategoryIDs) > 0 {
			filters["exclude_category_ids"] = recentCategoryIDs
		}
	}

	// Add context-specific filters
	rs.applyContextSpecificFilters(filters, req.ContextType, profile)
}

// applyContextSpecificFilters applies filters based on the recommendation context
func (rs *RecommendationServiceV2) applyContextSpecificFilters(filters map[string]interface{}, contextType string, profile *dto.CustomerProfile) {
	switch contextType {
	case "homepage":
		// Homepage: Focus on popular, high-rated products
		filters["homepage_boost"] = true
		filters["rating_min"] = 4.2
		filters["popularity_score_min"] = 100
	case "category_page":
		// Category page: Allow more variety within the category
		filters["rating_min"] = 3.5
		filters["include_new_arrivals"] = true
	case "product_detail":
		// Product detail page: Focus on similar and complementary products
		filters["similarity_boost"] = true
		filters["complementary_products"] = true
	case "cart":
		// Cart page: Suggest complementary and frequently bought together items
		filters["frequently_bought_together"] = true
		filters["complementary_products"] = true
		filters["price_max"] = 200 // Limit add-on suggestions
	case "checkout":
		// Checkout: Last-minute add-ons and accessories
		filters["quick_add_suitable"] = true
		filters["price_max"] = 100
		filters["category_types"] = []string{"accessories", "add-ons", "consumables"}
	case "post_purchase":
		// Post-purchase: Complementary items and replenishment
		filters["replenishment_suitable"] = true
		filters["complementary_products"] = true
	case "email_campaign":
		// Email campaigns: Personalized picks based on engagement
		if profile.IsPremium {
			filters["exclusive_items"] = true
		}
		filters["personalization_boost"] = true
	}
}

// getCurrentSeasonalTags returns relevant seasonal tags for filtering
func (rs *RecommendationServiceV2) getCurrentSeasonalTags() []string {
	now := time.Now()
	month := now.Month()

	switch {
	case month >= 12 || month <= 2:
		return []string{"winter", "holiday", "gift", "warm", "indoor"}
	case month >= 3 && month <= 5:
		return []string{"spring", "outdoor", "fresh", "renewal", "gardening"}
	case month >= 6 && month <= 8:
		return []string{"summer", "vacation", "outdoor", "travel", "cooling"}
	case month >= 9 && month <= 11:
		return []string{"autumn", "back-to-school", "cozy", "preparation", "harvest"}
	default:
		return []string{}
	}
}

// extractRecentCategoryIDs extracts category IDs from recent purchase history
func (rs *RecommendationServiceV2) extractRecentCategoryIDs(purchaseHistory []dto.PurchaseItem, days int) []int {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	categoryMap := make(map[int]bool)

	for _, purchase := range purchaseHistory {
		if purchase.PurchasedAt.After(cutoffDate) {
			categoryMap[purchase.CategoryID] = true
		}
	}

	categoryIDs := make([]int, 0, len(categoryMap))
	for categoryID := range categoryMap {
		categoryIDs = append(categoryIDs, categoryID)
	}

	return categoryIDs
}

func (rs *RecommendationServiceV2) generateSemanticInsights(results []dto.ProductRecommendationV2, understanding *dto.QueryUnderstanding) *dto.SemanticInsights {
	insights := &dto.SemanticInsights{
		ConfidenceLevel: 0.8,
	}

	// Extract entities from understanding
	if understanding != nil {
		insights.QueryIntent = understanding.Intent
		insights.ExtractedEntities = understanding.Entities
	}

	// Generate semantic clusters based on results
	categoryGroups := make(map[string][]dto.ProductRecommendationV2)
	for _, product := range results {
		categoryGroups[product.CategoryName] = append(categoryGroups[product.CategoryName], product)
	}

	for categoryName, products := range categoryGroups {
		if len(products) >= 2 {
			cluster := dto.SemanticCluster{
				ClusterID: fmt.Sprintf("cluster_%s", strings.ToLower(strings.ReplaceAll(categoryName, " ", "_"))),
				Theme:     categoryName,
				Relevance: float64(len(products)) / float64(len(results)),
			}

			for _, product := range products {
				cluster.ProductIDs = append(cluster.ProductIDs, product.ProductID)
				cluster.Keywords = append(cluster.Keywords, product.Tags...)
			}

			insights.SemanticClusters = append(insights.SemanticClusters, cluster)
		}
	}

	return insights
}

func (rs *RecommendationServiceV2) createKnowledgeBasedQuery(profile *dto.CustomerProfile, intent, contextDescription string) string {
	query := fmt.Sprintf(`
Customer Profile Analysis for Product Recommendations:

Customer Information:
- Total Spending: %.2f
- Order Count: %d
- Premium Customer: %t
- Preferred Categories: %v
- Preferred Brands: %v
- Lifestyle Tags: %v

Intent: %s
Context: %s

Please recommend products that would be most suitable for this customer based on their profile, preferences, and the given intent/context. Focus on products that align with their spending patterns, preferred categories, and lifestyle.
`,
		profile.TotalSpent,
		profile.OrderCount,
		profile.IsPremium,
		profile.PreferredCategories,
		profile.PreferredBrands,
		profile.LifestyleTags,
		intent,
		contextDescription,
	)

	return query
}

func (rs *RecommendationServiceV2) createComprehensiveQuery(profile *dto.CustomerProfile, contextType string) string {
	return fmt.Sprintf(`
Generate product recommendations for a customer with the following profile:

Customer Profile:
- Total Spent: %.2f
- Order Count: %d
- Is Premium: %t
- Preferred Categories: %v
- Preferred Brands: %v
- Context: %s

Please provide detailed product recommendations that match this customer's preferences and purchasing behavior.
`,
		profile.TotalSpent,
		profile.OrderCount,
		profile.IsPremium,
		profile.PreferredCategories,
		profile.PreferredBrands,
		contextType,
	)
}

func (rs *RecommendationServiceV2) extractRecommendationsFromKB(ctx context.Context, kbResponse *RAGResponse, profile *dto.CustomerProfile, limit int) ([]dto.ProductRecommendationV2, []dto.ReasoningStep, error) {
	recommendations := make([]dto.ProductRecommendationV2, 0)
	reasoningChain := make([]dto.ReasoningStep, 0)
	extractedProductIDs := make([]uuid.UUID, 0)

	// Check if kbResponse and Results are valid
	if kbResponse == nil || kbResponse.Results == nil {
		return recommendations, reasoningChain, nil
	}

	// Parse knowledge base results and extract product recommendations
	for i, result := range kbResponse.Results {
		if i >= limit {
			break
		}

		// Create reasoning step
		step := dto.ReasoningStep{
			Step:        i + 1,
			Description: fmt.Sprintf("Knowledge base result: %s", result.Content[:min(100, len(result.Content))]),
			Confidence:  result.Score,
			Source:      result.Source,
		}
		reasoningChain = append(reasoningChain, step)

		// Extract product IDs from KB content using AI analysis
		extractedIDs := rs.extractProductIDsFromKBContent(ctx, result.Content)
		if len(extractedIDs) > 0 {
			extractedProductIDs = append(extractedProductIDs, extractedIDs...)
		}
	}

	// If we extracted product IDs, enhance with full product details from repository
	if len(extractedProductIDs) > 0 {
		fullProducts, err := rs.repo.GetProductsByIDs(ctx, extractedProductIDs)
		if err != nil {
			log.Printf("Warning: failed to get full product details for KB recommendations: %v", err)
		} else {
			// Create recommendations from full product data with KB insights
			for _, product := range fullProducts {
				// Enhance product with KB-specific metadata
				product.Reason = "Based on knowledge base analysis and recommendations"
				product.ConfidenceScore = 0.8  // Default confidence from KB
				product.SimilarityScore = 0.75 // Default similarity from KB

				// Add KB source information if available
				if len(kbResponse.Results) > 0 {
					product.Reason = fmt.Sprintf("Based on knowledge base insights: %s",
						kbResponse.Results[0].Content[:min(100, len(kbResponse.Results[0].Content))])
				}

				recommendations = append(recommendations, product)
			}
		}
	}

	// Fallback to semantic search using KB content if no specific products extracted
	if len(recommendations) == 0 && len(kbResponse.Results) > 0 {
		// Use the first KB result as a search query (extract key terms)
		searchQuery := rs.extractKeyTermsFromText(kbResponse.Results[0].Content)
		if searchQuery != "" {
			filters := rs.buildPersonalizedFilters(profile, &dto.RecommendationRequestV2{})
			semanticResults, err := rs.rag.GetProductsWithSemanticSearch(ctx, searchQuery, limit, filters)
			if err == nil {
				// Extract product IDs from semantic search results
				productIDs := make([]uuid.UUID, len(semanticResults.Results))
				for i, result := range semanticResults.Results {
					productIDs[i] = result.ProductID
				}

				// Get full product details from repository
				fullProducts, err := rs.repo.GetProductsByIDs(ctx, productIDs)
				if err == nil {
					// Create a map for quick lookup of full product details
					productMap := make(map[uuid.UUID]dto.ProductRecommendationV2)
					for _, product := range fullProducts {
						productMap[product.ProductID] = product
					}

					// Enhance semantic results with full product information and KB reasoning
					for _, semanticResult := range semanticResults.Results {
						if fullProduct, exists := productMap[semanticResult.ProductID]; exists {
							// Copy semantic metadata and scores to the full product
							fullProduct.VectorMetadata = &dto.VectorMetadata{
								DistanceScore:   semanticResult.DistanceScore,
								SearchMethod:    semanticResult.SearchMethod,
								MatchedCriteria: semanticResult.MatchedCriteria,
								EmbeddingModel:  semanticResult.EmbeddingModel,
							}
							fullProduct.SimilarityScore = semanticResult.SimilarityScore
							fullProduct.ConfidenceScore = semanticResult.ConfidenceScore
							fullProduct.Reason = fmt.Sprintf("Based on knowledge base insights: %s", searchQuery)
							recommendations = append(recommendations, fullProduct)
						}
					}
				}
			}
		}
	}

	// Final fallback to category-based recommendations
	if len(recommendations) == 0 && len(profile.PreferredCategories) > 0 {
		categoryProducts, err := rs.repo.GetProductsByCategory(ctx, profile.PreferredCategories[0], limit)
		if err == nil {
			for i := range categoryProducts {
				categoryProducts[i].Reason = "Based on your preferred categories and knowledge base analysis"
			}
			recommendations = categoryProducts
		}
	}

	return recommendations, reasoningChain, nil
}

func (rs *RecommendationServiceV2) createKnowledgeInsights(kbResponse *RAGResponse) *dto.KnowledgeInsights {
	insights := &dto.KnowledgeInsights{
		ConfidenceLevel: 0.8,
		RelatedConcepts: make(map[string]float64),
	}

	// Extract insights from knowledge base response
	for _, result := range kbResponse.Results {
		insights.SourceDocuments = append(insights.SourceDocuments, result.Source)

		// Extract topics from content (simplified)
		words := strings.Fields(strings.ToLower(result.Content))
		for _, word := range words {
			if len(word) > 4 { // Only consider meaningful words
				insights.RelatedConcepts[word] = result.Score
			}
		}
	}

	return insights
}

func (rs *RecommendationServiceV2) enhanceWithAIExplanations(ctx context.Context, recommendations []dto.ProductRecommendationV2, profile *dto.CustomerProfile, contextType string, metrics *dto.PerformanceMetrics) ([]dto.ProductRecommendationV2, error) {
	if len(recommendations) == 0 {
		return recommendations, nil
	}

	startTime := time.Now()

	// Create prompt for AI enhancement
	prompt := rs.createEnhancementPromptV2(recommendations, profile, contextType)

	// Get chat response
	chatResponse, err := rs.chatService.GenerateResponse(ctx, prompt)
	if err != nil {
		return recommendations, err
	}

	// Parse and enhance recommendations
	enhanced, err := rs.parseAIEnhancementsV2(chatResponse.Content, recommendations)
	if err != nil {
		return recommendations, err
	}

	metrics.AIProcessingTimeMs += time.Since(startTime).Milliseconds()

	return enhanced, nil
}

func (rs *RecommendationServiceV2) generateTrendInsights(ctx context.Context, trendingProducts []dto.TrendingProductV2, timeRange string) (*dto.TrendInsights, error) {
	// Create AI prompt for trend analysis
	prompt := fmt.Sprintf(`
Analyze the following trending products for the %s period and provide insights:

Products:
%s

Please provide a JSON response with trend insights including emerging categories, declining categories, seasonal factors, and market drivers.
`, timeRange, rs.formatTrendingProductsForAI(trendingProducts))

	chatResponse, err := rs.chatService.GenerateResponse(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var insights dto.TrendInsights
	err = json.Unmarshal([]byte(chatResponse.Content), &insights)
	if err != nil {
		// Fallback to basic insights
		insights = dto.TrendInsights{
			EmergingCategories: []string{"Electronics", "Fashion"},
			SeasonalFactors:    []string{"Holiday shopping", "Back to school"},
			MarketDrivers:      []string{"Consumer trends", "Price sensitivity"},
		}
	}

	return &insights, nil
}

// Helper methods

func (rs *RecommendationServiceV2) filterOwnedProductsV2(recommendations []dto.ProductRecommendationV2, purchaseHistory []dto.PurchaseItem) []dto.ProductRecommendationV2 {
	ownedProducts := make(map[uuid.UUID]bool)
	for _, purchase := range purchaseHistory {
		ownedProducts[purchase.ProductID] = true
	}

	var filtered []dto.ProductRecommendationV2
	for _, rec := range recommendations {
		if !ownedProducts[rec.ProductID] {
			filtered = append(filtered, rec)
		}
	}

	return filtered
}

func (rs *RecommendationServiceV2) filterByPriceRange(recommendations []dto.ProductRecommendationV2, minPrice, maxPrice *float64) []dto.ProductRecommendationV2 {
	var filtered []dto.ProductRecommendationV2

	for _, rec := range recommendations {
		if minPrice != nil && rec.Price < *minPrice {
			continue
		}
		if maxPrice != nil && rec.Price > *maxPrice {
			continue
		}
		filtered = append(filtered, rec)
	}

	return filtered
}

func (rs *RecommendationServiceV2) removeDuplicateProductsV2(recommendations []dto.ProductRecommendationV2) []dto.ProductRecommendationV2 {
	seen := make(map[uuid.UUID]bool)
	var unique []dto.ProductRecommendationV2

	for _, rec := range recommendations {
		if !seen[rec.ProductID] {
			seen[rec.ProductID] = true
			unique = append(unique, rec)
		}
	}

	return unique
}

func (rs *RecommendationServiceV2) createEnhancementPromptV2(recommendations []dto.ProductRecommendationV2, profile *dto.CustomerProfile, contextType string) string {
	return fmt.Sprintf(`
Enhance the following product recommendations with AI-generated explanations and insights:

Customer Profile:
- Total Spent: %.2f
- Order Count: %d
- Preferred Categories: %v
- Preferred Brands: %v
- Context: %s

Products to enhance:
%s

For each product, provide:
1. A personalized reason why this product is recommended for this customer
2. Key features that match the customer's profile
3. Confidence score (0.0 to 1.0)

Format as JSON array with enhanced product information.
`,
		profile.TotalSpent,
		profile.OrderCount,
		profile.PreferredCategories,
		profile.PreferredBrands,
		contextType,
		rs.formatRecommendationsForAIV2(recommendations),
	)
}

func (rs *RecommendationServiceV2) parseAIEnhancementsV2(content string, recommendations []dto.ProductRecommendationV2) ([]dto.ProductRecommendationV2, error) {
	// This is a simplified parsing implementation
	// In a real system, you would have more sophisticated JSON parsing with proper error handling

	enhanced := make([]dto.ProductRecommendationV2, len(recommendations))
	copy(enhanced, recommendations)

	// For now, add basic AI insights to each recommendation
	for i := range enhanced {
		if enhanced[i].AIInsights == nil {
			enhanced[i].AIInsights = &dto.ProductAIInsights{
				KeyFeatures:    []string{"High quality", "Popular choice", "Great value"},
				UseCases:       []string{"Daily use", "Professional", "Personal"},
				TargetAudience: []string{"General consumers", "Professionals"},
			}
		}

		enhanced[i].Reason = fmt.Sprintf("Recommended based on your preferences and purchasing history. This product aligns with your interests in %s category.", enhanced[i].CategoryName)

		// Enhance confidence score slightly
		if enhanced[i].ConfidenceScore < 0.9 {
			enhanced[i].ConfidenceScore += 0.1
		}
	}

	return enhanced, nil
}

func (rs *RecommendationServiceV2) createExplanationPrompt(recommendationID uuid.UUID, profile *dto.CustomerProfile) string {
	return fmt.Sprintf(`
Provide a detailed explanation for why product ID %s was recommended to this customer:

Customer Profile:
- Total Spent: %.2f
- Order Count: %d
- Preferred Categories: %v
- Preferred Brands: %v
- Lifestyle Tags: %v

Please explain:
1. Why this product matches their preferences
2. How it relates to their purchase history
3. What specific factors influenced this recommendation
4. The confidence level of this recommendation

Provide a comprehensive, personalized explanation.
`,
		recommendationID.String(),
		profile.TotalSpent,
		profile.OrderCount,
		profile.PreferredCategories,
		profile.PreferredBrands,
		profile.LifestyleTags,
	)
}

func (rs *RecommendationServiceV2) parseExplanationFactors(content string, profile *dto.CustomerProfile) []dto.ExplanationFactor {
	// Simplified implementation - in reality, you would parse the AI response more sophisticatedly
	factors := []dto.ExplanationFactor{
		{
			Factor:      "Purchase History",
			Weight:      0.3,
			Impact:      "positive",
			Description: "Based on previous purchases in similar categories",
			Confidence:  0.8,
		},
		{
			Factor:      "Category Preference",
			Weight:      0.25,
			Impact:      "positive",
			Description: "Matches preferred product categories",
			Confidence:  0.9,
		},
		{
			Factor:      "Price Range",
			Weight:      0.2,
			Impact:      "positive",
			Description: "Within typical spending range",
			Confidence:  0.7,
		},
		{
			Factor:      "Brand Affinity",
			Weight:      0.15,
			Impact:      "positive",
			Description: "From preferred or similar brands",
			Confidence:  0.6,
		},
		{
			Factor:      "Seasonal Trends",
			Weight:      0.1,
			Impact:      "neutral",
			Description: "Trending product for current season",
			Confidence:  0.5,
		},
	}

	return factors
}

func (rs *RecommendationServiceV2) formatRecommendationsForAIV2(recommendations []dto.ProductRecommendationV2) string {
	result := ""
	for _, rec := range recommendations {
		result += fmt.Sprintf("- ID: %s, Name: %s, Category: %s, Price: %.2f, Rating: %.1f, Description: %s\n",
			rec.ProductID.String(), rec.Name, rec.CategoryName, rec.Price, rec.RatingAverage, rec.Description)
	}
	return result
}

func (rs *RecommendationServiceV2) formatTrendingProductsForAI(products []dto.TrendingProductV2) string {
	result := ""
	for _, product := range products {
		result += fmt.Sprintf("- Name: %s, Category: %s, Trend Score: %.2f, Velocity: %.2f\n",
			product.Name, product.CategoryName, product.TrendScore, product.TrendVelocity)
	}
	return result
}

// Utility functions

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// extractProductIDsFromKBContent extracts product IDs from knowledge base content using AI analysis
func (rs *RecommendationServiceV2) extractProductIDsFromKBContent(ctx context.Context, content string) []uuid.UUID {
	// Create AI prompt to extract product IDs from KB content
	prompt := fmt.Sprintf(`
Extract product IDs from the following knowledge base content.
Look for any product identifiers, UUIDs, or product references.

Content: %s

IMPORTANT: Return ONLY a valid JSON array of strings without any markdown formatting, code blocks, or explanatory text.
Example: ["uuid1", "uuid2", "uuid3"]
If no valid product IDs are found, return: []
`, content)

	// Use chat service to analyze content
	chatResponse, err := rs.chatService.GenerateResponse(ctx, prompt)
	if err != nil {
		log.Printf("Warning: failed to extract product IDs from KB content: %v", err)
		return []uuid.UUID{}
	}

	// Clean the response by removing markdown code blocks if present
	responseContent := rs.extractJSONFromMarkdown(chatResponse.Content)

	// Parse the response to extract UUIDs
	var productIDStrings []string
	err = json.Unmarshal([]byte(responseContent), &productIDStrings)
	if err != nil {
		log.Printf("Warning: failed to parse extracted product IDs: %v", err)
		log.Printf("Raw response content: %s", chatResponse.Content)

		// Fallback: try to extract UUIDs using regex pattern matching
		productIDs := rs.extractUUIDsWithRegex(chatResponse.Content)
		if len(productIDs) > 0 {
			log.Printf("Successfully extracted %d UUIDs using regex fallback", len(productIDs))
			return productIDs
		}

		return []uuid.UUID{}
	}

	// Convert strings to UUIDs
	var productIDs []uuid.UUID
	for _, idStr := range productIDStrings {
		if id, err := uuid.Parse(idStr); err == nil {
			productIDs = append(productIDs, id)
		}
	}

	return productIDs
}

// extractKeyTermsFromText extracts key terms from text for semantic search
func (rs *RecommendationServiceV2) extractKeyTermsFromText(text string) string {
	// Simple implementation: extract meaningful words
	words := strings.Fields(strings.ToLower(text))
	var keyTerms []string

	// Filter out common stop words and keep meaningful terms
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "is": true, "are": true, "was": true, "were": true,
		"be": true, "been": true, "being": true, "have": true, "has": true, "had": true,
		"do": true, "does": true, "did": true, "will": true, "would": true, "could": true,
		"should": true, "may": true, "might": true, "must": true, "can": true,
		"this": true, "that": true, "these": true, "those": true, "i": true, "you": true,
		"he": true, "she": true, "it": true, "we": true, "they": true, "me": true,
		"him": true, "her": true, "us": true, "them": true,
	}

	for _, word := range words {
		// Remove punctuation and keep only alphanumeric characters
		cleanWord := strings.Trim(word, ".,!?;:\"'()[]{}+-=_*&^%$#@~`")
		if len(cleanWord) > 3 && !stopWords[cleanWord] {
			keyTerms = append(keyTerms, cleanWord)
		}
	}

	// Return top 5 key terms joined with spaces
	if len(keyTerms) > 5 {
		keyTerms = keyTerms[:5]
	}

	return strings.Join(keyTerms, " ")
}

// extractJSONFromMarkdown extracts JSON content from markdown code blocks
func (rs *RecommendationServiceV2) extractJSONFromMarkdown(content string) string {
	// Remove leading/trailing whitespace
	content = strings.TrimSpace(content)

	// Check if content is wrapped in markdown code blocks
	if strings.HasPrefix(content, "```json") {
		// Extract content between ```json and ```
		lines := strings.Split(content, "\n")
		var jsonLines []string
		inCodeBlock := false

		for _, line := range lines {
			if strings.HasPrefix(line, "```json") {
				inCodeBlock = true
				continue
			}
			if strings.HasPrefix(line, "```") && inCodeBlock {
				break
			}
			if inCodeBlock {
				jsonLines = append(jsonLines, line)
			}
		}

		return strings.Join(jsonLines, "\n")
	}

	// Check if content is wrapped in simple code blocks
	if strings.HasPrefix(content, "```") {
		// Extract content between ``` and ```
		lines := strings.Split(content, "\n")
		var jsonLines []string
		inCodeBlock := false

		for _, line := range lines {
			if strings.HasPrefix(line, "```") && !inCodeBlock {
				inCodeBlock = true
				continue
			}
			if strings.HasPrefix(line, "```") && inCodeBlock {
				break
			}
			if inCodeBlock {
				jsonLines = append(jsonLines, line)
			}
		}

		return strings.Join(jsonLines, "\n")
	}

	// Try to find JSON array pattern in the content
	startIndex := strings.Index(content, "[")
	endIndex := strings.LastIndex(content, "]")

	if startIndex != -1 && endIndex != -1 && endIndex > startIndex {
		return content[startIndex : endIndex+1]
	}

	// Return original content if no markdown formatting found
	return content
}

// extractUUIDsWithRegex extracts UUIDs from text using regex pattern matching
func (rs *RecommendationServiceV2) extractUUIDsWithRegex(content string) []uuid.UUID {
	// Find all UUID matches in the content
	matches := make(map[string]bool) // Use map to avoid duplicates

	// Split content into words and check each word
	words := strings.Fields(content)
	for _, word := range words {
		// Clean the word by removing common punctuation
		cleanWord := strings.Trim(word, ".,!?;:\"'()[]{}+-=_*&^%$#@~`")

		// Try to parse as UUID
		if _, err := uuid.Parse(cleanWord); err == nil {
			matches[cleanWord] = true
		}
	}

	// Also check for UUIDs in quoted strings
	if strings.Contains(content, `"`) {
		// Extract quoted UUIDs
		parts := strings.Split(content, `"`)
		for i := 1; i < len(parts); i += 2 { // Check every odd index (quoted content)
			if _, err := uuid.Parse(parts[i]); err == nil {
				matches[parts[i]] = true
			}
		}
	}

	// Convert map keys to UUID slice
	var productIDs []uuid.UUID
	for uuidStr := range matches {
		if id, err := uuid.Parse(uuidStr); err == nil {
			productIDs = append(productIDs, id)
		}
	}

	return productIDs
}

// convertRAGResultsToProducts converts RAG search results to ProductRecommendationV2
func (rs *RecommendationServiceV2) convertRAGResultsToProducts(ctx context.Context, ragResults []RAGSearchResult, searchMethod string) ([]dto.ProductRecommendationV2, error) {
	if len(ragResults) == 0 {
		return []dto.ProductRecommendationV2{}, nil
	}

	// Extract product IDs from RAG search results
	productIDs := make([]uuid.UUID, len(ragResults))
	for i, result := range ragResults {
		productIDs[i] = result.ProductID
	}

	// Get full product details from repository
	fullProducts, err := rs.repo.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get full product details: %w", err)
	}

	// Create a map for quick lookup of full product details
	productMap := make(map[uuid.UUID]dto.ProductRecommendationV2)
	for _, product := range fullProducts {
		productMap[product.ProductID] = product
	}

	// Enhance RAG results with full product information
	results := make([]dto.ProductRecommendationV2, 0, len(ragResults))
	for _, ragResult := range ragResults {
		if fullProduct, exists := productMap[ragResult.ProductID]; exists {
			// Copy RAG metadata and scores to the full product
			fullProduct.VectorMetadata = &dto.VectorMetadata{
				DistanceScore:    ragResult.DistanceScore,
				SearchMethod:     ragResult.SearchMethod,
				MatchedCriteria:  ragResult.MatchedCriteria,
				EmbeddingModel:   ragResult.EmbeddingModel,
				SemanticClusters: ragResult.SemanticClusters,
			}
			fullProduct.SimilarityScore = ragResult.SimilarityScore
			fullProduct.ConfidenceScore = ragResult.ConfidenceScore

			// Set appropriate reason based on search method
			switch searchMethod {
			case "vector_search":
				fullProduct.Reason = "Similar to your selected product based on vector analysis"
			case "semantic_search":
				fullProduct.Reason = "Matches your search query based on semantic understanding"
			case "hybrid_search":
				fullProduct.Reason = "Recommended based on both similarity and semantic relevance"
			default:
				fullProduct.Reason = fmt.Sprintf("Recommended using %s", searchMethod)
			}

			results = append(results, fullProduct)
		}
	}

	return results, nil
}

// buildComprehensiveProductText creates a comprehensive text representation of a product for vector embedding
func (rs *RecommendationServiceV2) buildComprehensiveProductText(product dto.ProductRecommendationV2) string {
	var textComponents []string

	// Primary product information (highest weight)
	textComponents = append(textComponents, fmt.Sprintf("Product: %s", product.Name))
	if product.Description != "" {
		textComponents = append(textComponents, fmt.Sprintf("Description: %s", product.Description))
	}

	// Category and brand information (high weight)
	if product.CategoryName != "" {
		textComponents = append(textComponents, fmt.Sprintf("Category: %s", product.CategoryName))
	}
	if product.Brand != "" {
		textComponents = append(textComponents, fmt.Sprintf("Brand: %s", product.Brand))
	}

	// Price information with semantic context
	priceCategory := rs.categorizePriceRange(product.Price)
	textComponents = append(textComponents, fmt.Sprintf("Price Range: %s", priceCategory))

	// Add price discount information if available
	if product.OriginalPrice != nil && *product.OriginalPrice > product.Price {
		discountPercent := ((*product.OriginalPrice - product.Price) / *product.OriginalPrice) * 100
		discountCategory := rs.categorizeDiscount(discountPercent)
		textComponents = append(textComponents, fmt.Sprintf("Discount: %s", discountCategory))
	}

	// Tags and keywords (medium weight)
	if len(product.Tags) > 0 {
		textComponents = append(textComponents, fmt.Sprintf("Tags: %s", strings.Join(product.Tags, ", ")))
	}

	// Quality and popularity indicators (medium weight)
	if product.RatingAverage > 0 {
		ratingCategory := rs.categorizeRating(product.RatingAverage)
		textComponents = append(textComponents, fmt.Sprintf("Quality: %s", ratingCategory))

		// Add review volume context for quality assessment
		if product.RatingCount > 0 {
			reviewVolumeCategory := rs.categorizeReviewVolume(product.RatingCount)
			textComponents = append(textComponents, fmt.Sprintf("Review Volume: %s", reviewVolumeCategory))
		}
	}

	if product.PopularityScore > 0 {
		popularityCategory := rs.categorizePopularity(product.PopularityScore)
		textComponents = append(textComponents, fmt.Sprintf("Popularity: %s", popularityCategory))
	}

	// AI insights for enhanced context (if available)
	if product.AIInsights != nil {
		if len(product.AIInsights.KeyFeatures) > 0 {
			textComponents = append(textComponents, fmt.Sprintf("Key Features: %s", strings.Join(product.AIInsights.KeyFeatures, ", ")))
		}
		if len(product.AIInsights.UseCases) > 0 {
			textComponents = append(textComponents, fmt.Sprintf("Use Cases: %s", strings.Join(product.AIInsights.UseCases, ", ")))
		}
		if len(product.AIInsights.TargetAudience) > 0 {
			textComponents = append(textComponents, fmt.Sprintf("Target Audience: %s", strings.Join(product.AIInsights.TargetAudience, ", ")))
		}

		// Add sentiment analysis if available
		if product.AIInsights.SentimentAnalysis != nil {
			sentiment := product.AIInsights.SentimentAnalysis
			textComponents = append(textComponents, fmt.Sprintf("Customer Sentiment: %s", sentiment.OverallSentiment))
			if len(sentiment.KeyTopics) > 0 {
				textComponents = append(textComponents, fmt.Sprintf("Review Topics: %s", strings.Join(sentiment.KeyTopics, ", ")))
			}
		}

		// Add trend information if available
		if product.AIInsights.TrendAnalysis != nil {
			trend := product.AIInsights.TrendAnalysis
			textComponents = append(textComponents, fmt.Sprintf("Market Trend: %s %s", trend.TrendDirection, trend.MarketPosition))
			if trend.SeasonalPattern != "" {
				textComponents = append(textComponents, fmt.Sprintf("Seasonal Pattern: %s", trend.SeasonalPattern))
			}
		}
	}

	// Add seasonal context based on current time
	seasonalContext := rs.getCurrentSeasonalContext()
	if seasonalContext != "" {
		textComponents = append(textComponents, fmt.Sprintf("Seasonal Context: %s", seasonalContext))
	}

	// Join all components with structured separators for better embedding understanding
	return strings.Join(textComponents, ". ")
}

// categorizePriceRange converts numerical price to categorical range for better semantic understanding
func (rs *RecommendationServiceV2) categorizePriceRange(price float64) string {
	switch {
	case price < 10:
		return "budget low-cost affordable"
	case price < 50:
		return "budget moderate affordable"
	case price < 100:
		return "mid-range affordable quality"
	case price < 300:
		return "mid-range quality premium"
	case price < 1000:
		return "premium high-quality expensive"
	default:
		return "luxury premium high-end expensive"
	}
}

// categorizeRating converts numerical rating to categorical quality descriptor
func (rs *RecommendationServiceV2) categorizeRating(rating float64) string {
	switch {
	case rating >= 4.5:
		return "excellent highly-rated top-quality"
	case rating >= 4.0:
		return "very-good highly-rated quality"
	case rating >= 3.5:
		return "good rated decent-quality"
	case rating >= 3.0:
		return "average moderate-quality"
	case rating >= 2.0:
		return "below-average poor-quality"
	default:
		return "poor low-quality"
	}
}

// categorizePopularity converts numerical popularity score to categorical descriptor
func (rs *RecommendationServiceV2) categorizePopularity(popularity int) string {
	switch {
	case popularity >= 1000:
		return "extremely-popular trending bestseller"
	case popularity >= 500:
		return "very-popular trending well-known"
	case popularity >= 100:
		return "popular well-known"
	case popularity >= 50:
		return "moderately-popular known"
	default:
		return "niche specialized emerging"
	}
}

// categorizeDiscount converts discount percentage to semantic descriptor
func (rs *RecommendationServiceV2) categorizeDiscount(discountPercent float64) string {
	switch {
	case discountPercent >= 50:
		return "massive-discount clearance heavily-discounted"
	case discountPercent >= 30:
		return "major-discount significant-savings heavily-discounted"
	case discountPercent >= 20:
		return "good-discount discounted on-sale"
	case discountPercent >= 10:
		return "moderate-discount slightly-discounted on-sale"
	default:
		return "minor-discount small-savings"
	}
}

// categorizeReviewVolume converts review count to reliability descriptor
func (rs *RecommendationServiceV2) categorizeReviewVolume(reviewCount int) string {
	switch {
	case reviewCount >= 1000:
		return "extensively-reviewed well-established trusted"
	case reviewCount >= 500:
		return "well-reviewed established reliable"
	case reviewCount >= 100:
		return "moderately-reviewed verified"
	case reviewCount >= 10:
		return "some-reviews emerging"
	default:
		return "few-reviews new-product"
	}
}

// getCurrentSeasonalContext returns seasonal keywords based on current time
func (rs *RecommendationServiceV2) getCurrentSeasonalContext() string {
	now := time.Now()
	month := now.Month()

	switch {
	case month >= 12 || month <= 2:
		return "winter holiday-season cold-weather gift-giving"
	case month >= 3 && month <= 5:
		return "spring renewal fresh-start outdoor-activities"
	case month >= 6 && month <= 8:
		return "summer vacation hot-weather outdoor-recreation"
	case month >= 9 && month <= 11:
		return "autumn back-to-school harvest-season preparation"
	default:
		return ""
	}
}
