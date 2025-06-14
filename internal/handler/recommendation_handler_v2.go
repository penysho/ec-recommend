package handler

import (
	"ec-recommend/internal/dto"
	"ec-recommend/internal/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RecommendationHandlerV2 handles advanced RAG-based recommendation requests using Amazon Bedrock Knowledge Bases
type RecommendationHandlerV2 struct {
	recommendationServiceV2 interfaces.RecommendationServiceV2Interface
}

// NewRecommendationHandlerV2 creates a new enhanced recommendation handler instance
func NewRecommendationHandlerV2(recommendationServiceV2 interfaces.RecommendationServiceV2Interface) *RecommendationHandlerV2 {
	return &RecommendationHandlerV2{
		recommendationServiceV2: recommendationServiceV2,
	}
}

// GetRecommendationsV2 handles GET /api/v2/recommendations
// @Summary Get advanced AI-powered product recommendations for a customer using RAG
// @Description Generate personalized product recommendations using Amazon Bedrock Knowledge Bases with vector search and semantic understanding
// @Tags recommendations-v2
// @Accept json
// @Produce json
// @Param customer_id query string true "Customer UUID"
// @Param recommendation_type query string false "Type of recommendation (hybrid, semantic, collaborative, vector_search, knowledge_based)" default(hybrid)
// @Param context_type query string false "Context where recommendations are shown (homepage, product_page, cart, checkout, search_results)" default(homepage)
// @Param query_text query string false "Natural language query for semantic search (e.g., 'Find products similar to wireless headphones for running')"
// @Param product_id query string false "Product UUID for similar product recommendations"
// @Param category_id query int false "Category ID for category-based recommendations"
// @Param price_range_min query float64 false "Minimum price range for filtering"
// @Param price_range_max query float64 false "Maximum price range for filtering"
// @Param limit query int false "Number of recommendations to return" default(10)
// @Param exclude_owned query bool false "Exclude already purchased products" default(false)
// @Param enable_explanation query bool false "Include AI-generated explanations for recommendations" default(true)
// @Success 200 {object} dto.RecommendationResponseV2
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v2/recommendations [get]
func (h *RecommendationHandlerV2) GetRecommendationsV2(c *gin.Context) {
	// Parse customer ID
	customerIDStr := c.Query("customer_id")
	if customerIDStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "customer_id is required",
		})
		return
	}

	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid customer_id format",
		})
		return
	}

	// Build request object
	req := &dto.RecommendationRequestV2{
		CustomerID:         customerID,
		RecommendationType: c.DefaultQuery("recommendation_type", "hybrid"),
		ContextType:        c.DefaultQuery("context_type", "homepage"),
		QueryText:          c.Query("query_text"),
		EnableExplanation:  true, // Default to true for V2
	}

	// Parse optional product ID
	if productIDStr := c.Query("product_id"); productIDStr != "" {
		productID, err := uuid.Parse(productIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "invalid product_id format",
			})
			return
		}
		req.ProductID = &productID
	}

	// Parse optional category ID
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "invalid category_id format",
			})
			return
		}
		req.CategoryID = &categoryID
	}

	// Parse price range
	if priceMinStr := c.Query("price_range_min"); priceMinStr != "" {
		priceMin, err := strconv.ParseFloat(priceMinStr, 64)
		if err != nil || priceMin < 0 {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "price_range_min must be a valid positive number",
			})
			return
		}
		req.PriceRangeMin = &priceMin
	}

	if priceMaxStr := c.Query("price_range_max"); priceMaxStr != "" {
		priceMax, err := strconv.ParseFloat(priceMaxStr, 64)
		if err != nil || priceMax < 0 {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "price_range_max must be a valid positive number",
			})
			return
		}
		req.PriceRangeMax = &priceMax
	}

	// Validate price range
	if req.PriceRangeMin != nil && req.PriceRangeMax != nil && *req.PriceRangeMin > *req.PriceRangeMax {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "price_range_min cannot be greater than price_range_max",
		})
		return
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 100 {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "limit must be a positive integer between 1 and 100",
			})
			return
		}
		req.Limit = limit
	} else {
		req.Limit = 10
	}

	// Parse exclude_owned
	if excludeOwnedStr := c.Query("exclude_owned"); excludeOwnedStr != "" {
		excludeOwned, err := strconv.ParseBool(excludeOwnedStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "exclude_owned must be a boolean value",
			})
			return
		}
		req.ExcludeOwned = excludeOwned
	}

	// Parse enable_explanation
	if enableExplanationStr := c.Query("enable_explanation"); enableExplanationStr != "" {
		enableExplanation, err := strconv.ParseBool(enableExplanationStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "enable_explanation must be a boolean value",
			})
			return
		}
		req.EnableExplanation = enableExplanation
	}

	// Get advanced recommendations
	response, err := h.recommendationServiceV2.GetRecommendationsV2(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "failed to get advanced recommendations: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// PostRecommendationsV2 handles POST /api/v2/recommendations
// @Summary Get advanced product recommendations with detailed request body
// @Description Generate personalized product recommendations using a detailed request body with RAG capabilities
// @Tags recommendations-v2
// @Accept json
// @Produce json
// @Param request body dto.RecommendationRequestV2 true "Advanced recommendation request"
// @Success 200 {object} dto.RecommendationResponseV2
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v2/recommendations [post]
func (h *RecommendationHandlerV2) PostRecommendationsV2(c *gin.Context) {
	var req dto.RecommendationRequestV2
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if req.CustomerID == uuid.Nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "customer_id is required",
		})
		return
	}

	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Validate price range
	if req.PriceRangeMin != nil && req.PriceRangeMax != nil && *req.PriceRangeMin > *req.PriceRangeMax {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "price_range_min cannot be greater than price_range_max",
		})
		return
	}

	// Get advanced recommendations
	response, err := h.recommendationServiceV2.GetRecommendationsV2(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "failed to get advanced recommendations: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetSemanticSearch handles GET /api/v2/search/semantic
// @Summary Perform semantic search using knowledge base
// @Description Search products using natural language queries with semantic understanding
// @Tags semantic-search
// @Accept json
// @Produce json
// @Param query query string true "Natural language search query"
// @Param customer_id query string false "Customer UUID for personalization"
// @Param category_id query int false "Category ID to filter results"
// @Param price_range_min query float64 false "Minimum price range for filtering"
// @Param price_range_max query float64 false "Maximum price range for filtering"
// @Param limit query int false "Number of results to return" default(10)
// @Success 200 {object} dto.SemanticSearchResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v2/search/semantic [get]
func (h *RecommendationHandlerV2) GetSemanticSearch(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "query parameter is required",
		})
		return
	}

	// Build search request
	req := &dto.SemanticSearchRequest{
		Query: query,
		Limit: 10,
	}

	// Parse optional customer ID for personalization
	if customerIDStr := c.Query("customer_id"); customerIDStr != "" {
		customerID, err := uuid.Parse(customerIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "invalid customer_id format",
			})
			return
		}
		req.CustomerID = &customerID
	}

	// Parse optional category ID
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "invalid category_id format",
			})
			return
		}
		req.CategoryID = &categoryID
	}

	// Parse price range filters
	if priceMinStr := c.Query("price_range_min"); priceMinStr != "" {
		priceMin, err := strconv.ParseFloat(priceMinStr, 64)
		if err != nil || priceMin < 0 {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "price_range_min must be a valid positive number",
			})
			return
		}
		req.PriceRangeMin = &priceMin
	}

	if priceMaxStr := c.Query("price_range_max"); priceMaxStr != "" {
		priceMax, err := strconv.ParseFloat(priceMaxStr, 64)
		if err != nil || priceMax < 0 {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "price_range_max must be a valid positive number",
			})
			return
		}
		req.PriceRangeMax = &priceMax
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 100 {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "limit must be a positive integer between 1 and 100",
			})
			return
		}
		req.Limit = limit
	}

	// Perform semantic search
	response, err := h.recommendationServiceV2.SemanticSearch(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "failed to perform semantic search: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetVectorSimilarProducts handles GET /api/v2/products/{product_id}/vector-similar
// @Summary Get products similar to a specific product using vector search
// @Description Find products similar to the given product using advanced vector similarity in knowledge base
// @Tags vector-search
// @Produce json
// @Param product_id path string true "Product UUID"
// @Param limit query int false "Number of similar products to return" default(10)
// @Param include_metadata query bool false "Include additional metadata in results" default(true)
// @Success 200 {object} dto.VectorSimilarityResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v2/products/{product_id}/vector-similar [get]
func (h *RecommendationHandlerV2) GetVectorSimilarProducts(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid product_id format",
		})
		return
	}

	// Build request
	req := &dto.VectorSimilarityRequest{
		ProductID:       productID,
		Limit:           10,
		IncludeMetadata: true,
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 100 {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "limit must be a positive integer between 1 and 100",
			})
			return
		}
		req.Limit = limit
	}

	// Parse include_metadata
	if includeMetadataStr := c.Query("include_metadata"); includeMetadataStr != "" {
		includeMetadata, err := strconv.ParseBool(includeMetadataStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "include_metadata must be a boolean value",
			})
			return
		}
		req.IncludeMetadata = includeMetadata
	}

	// Get vector similar products
	response, err := h.recommendationServiceV2.GetVectorSimilarProducts(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "failed to get vector similar products: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetKnowledgeBasedRecommendations handles GET /api/v2/recommendations/knowledge-based
// @Summary Get knowledge-based recommendations using RAG
// @Description Generate recommendations based on comprehensive knowledge base analysis
// @Tags knowledge-based
// @Accept json
// @Produce json
// @Param customer_id query string true "Customer UUID"
// @Param intent query string false "User intent or goal (e.g., 'workout', 'cooking', 'study')"
// @Param context_description query string false "Additional context description"
// @Param limit query int false "Number of recommendations to return" default(10)
// @Success 200 {object} dto.KnowledgeBasedRecommendationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v2/recommendations/knowledge-based [get]
func (h *RecommendationHandlerV2) GetKnowledgeBasedRecommendations(c *gin.Context) {
	// Parse customer ID
	customerIDStr := c.Query("customer_id")
	if customerIDStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "customer_id is required",
		})
		return
	}

	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid customer_id format",
		})
		return
	}

	// Build request
	req := &dto.KnowledgeBasedRecommendationRequest{
		CustomerID:         customerID,
		Intent:             c.Query("intent"),
		ContextDescription: c.Query("context_description"),
		Limit:              10,
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 100 {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "limit must be a positive integer between 1 and 100",
			})
			return
		}
		req.Limit = limit
	}

	// Get knowledge-based recommendations
	response, err := h.recommendationServiceV2.GetKnowledgeBasedRecommendations(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "failed to get knowledge-based recommendations: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetRecommendationExplanation handles GET /api/v2/recommendations/{recommendation_id}/explanation
// @Summary Get detailed explanation for a specific recommendation
// @Description Retrieve AI-generated explanation for why a product was recommended
// @Tags explanations
// @Produce json
// @Param recommendation_id path string true "Recommendation UUID"
// @Param customer_id query string true "Customer UUID"
// @Success 200 {object} dto.RecommendationExplanationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v2/recommendations/{recommendation_id}/explanation [get]
func (h *RecommendationHandlerV2) GetRecommendationExplanation(c *gin.Context) {
	recommendationIDStr := c.Param("recommendation_id")
	recommendationID, err := uuid.Parse(recommendationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid recommendation_id format",
		})
		return
	}

	customerIDStr := c.Query("customer_id")
	if customerIDStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "customer_id is required",
		})
		return
	}

	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid customer_id format",
		})
		return
	}

	// Get explanation
	response, err := h.recommendationServiceV2.GetRecommendationExplanation(c.Request.Context(), recommendationID, customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "failed to get recommendation explanation: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetTrendingProductsV2 handles GET /api/v2/products/trending
// @Summary Get trending products with enhanced AI insights
// @Description Retrieve currently trending products with AI-powered trend analysis
// @Tags trending-v2
// @Produce json
// @Param category_id query int false "Category ID to filter trending products"
// @Param time_range query string false "Time range for trend analysis (daily, weekly, monthly)" default(weekly)
// @Param limit query int false "Number of trending products to return" default(10)
// @Param include_insights query bool false "Include AI-generated trend insights" default(true)
// @Success 200 {object} dto.TrendingProductsResponseV2
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v2/products/trending [get]
func (h *RecommendationHandlerV2) GetTrendingProductsV2(c *gin.Context) {
	// Build request
	req := &dto.TrendingProductsRequestV2{
		TimeRange:       c.DefaultQuery("time_range", "weekly"),
		Limit:           10,
		IncludeInsights: true,
	}

	// Parse optional category ID
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "invalid category_id format",
			})
			return
		}
		req.CategoryID = &categoryID
	}

	// Validate time_range
	validTimeRanges := map[string]bool{
		"daily":   true,
		"weekly":  true,
		"monthly": true,
	}
	if !validTimeRanges[req.TimeRange] {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "time_range must be one of: daily, weekly, monthly",
		})
		return
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 || limit > 100 {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "limit must be a positive integer between 1 and 100",
			})
			return
		}
		req.Limit = limit
	}

	// Parse include_insights
	if includeInsightsStr := c.Query("include_insights"); includeInsightsStr != "" {
		includeInsights, err := strconv.ParseBool(includeInsightsStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "include_insights must be a boolean value",
			})
			return
		}
		req.IncludeInsights = includeInsights
	}

	// Get trending products with AI insights
	response, err := h.recommendationServiceV2.GetTrendingProductsV2(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "failed to get trending products: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
