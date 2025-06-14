package handler

import (
	"ec-recommend/internal/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RecommendationHandler handles recommendation-related HTTP requests
type RecommendationHandler struct {
	recommendationService RecommendationServiceInterface
}

// NewRecommendationHandler creates a new recommendation handler instance
func NewRecommendationHandler(recommendationService RecommendationServiceInterface) *RecommendationHandler {
	return &RecommendationHandler{
		recommendationService: recommendationService,
	}
}

// GetRecommendations handles GET /api/v1/recommendations
// @Summary Get product recommendations for a customer
// @Description Generate personalized product recommendations based on customer profile and preferences
// @Tags recommendations
// @Accept json
// @Produce json
// @Param customer_id query string true "Customer UUID"
// @Param recommendation_type query string false "Type of recommendation (similar, collaborative, content_based, hybrid)" default(hybrid)
// @Param context_type query string false "Context where recommendations are shown (homepage, product_page, cart, checkout)" default(homepage)
// @Param product_id query string false "Product UUID for similar product recommendations"
// @Param category_id query int false "Category ID for category-based recommendations"
// @Param limit query int false "Number of recommendations to return" default(10)
// @Param exclude_owned query bool false "Exclude already purchased products" default(false)
// @Success 200 {object} dto.RecommendationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/recommendations [get]
func (h *RecommendationHandler) GetRecommendations(c *gin.Context) {
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
	req := &dto.RecommendationRequest{
		CustomerID:         customerID,
		RecommendationType: c.DefaultQuery("recommendation_type", "hybrid"),
		ContextType:        c.DefaultQuery("context_type", "homepage"),
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

	// Get recommendations
	response, err := h.recommendationService.GetRecommendations(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "failed to get recommendations: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// PostRecommendations handles POST /api/v1/recommendations
// @Summary Get product recommendations with detailed request body
// @Description Generate personalized product recommendations using a detailed request body
// @Tags recommendations
// @Accept json
// @Produce json
// @Param request body dto.RecommendationRequest true "Recommendation request"
// @Success 200 {object} dto.RecommendationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/recommendations [post]
func (h *RecommendationHandler) PostRecommendations(c *gin.Context) {
	var req dto.RecommendationRequest
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

	// Get recommendations
	response, err := h.recommendationService.GetRecommendations(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "failed to get recommendations: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetCustomerProfile handles GET /api/v1/customers/{customer_id}/profile
// @Summary Get customer profile for recommendations
// @Description Retrieve comprehensive customer profile data used for generating recommendations
// @Tags customers
// @Produce json
// @Param customer_id path string true "Customer UUID"
// @Success 200 {object} dto.CustomerProfile
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/customers/{customer_id}/profile [get]
func (h *RecommendationHandler) GetCustomerProfile(c *gin.Context) {
	customerIDStr := c.Param("customer_id")
	customerID, err := uuid.Parse(customerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid customer_id format",
		})
		return
	}

	profile, err := h.recommendationService.GetCustomerProfile(c.Request.Context(), customerID)
	if err != nil {
		if err.Error() == "customer not found: "+customerID.String() {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Not Found",
				Message: "customer not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "failed to get customer profile: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// GetSimilarProducts handles GET /api/v1/products/{product_id}/similar
// @Summary Get products similar to a specific product
// @Description Find products similar to the given product using content-based filtering
// @Tags products
// @Produce json
// @Param product_id path string true "Product UUID"
// @Param limit query int false "Number of similar products to return" default(10)
// @Success 200 {object} []dto.ProductRecommendation
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/products/{product_id}/similar [get]
func (h *RecommendationHandler) GetSimilarProducts(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid product_id format",
		})
		return
	}

	// Parse limit
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 || parsedLimit > 100 {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "limit must be a positive integer between 1 and 100",
			})
			return
		}
		limit = parsedLimit
	}

	products, err := h.recommendationService.GetSimilarProducts(c.Request.Context(), productID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "failed to get similar products: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetTrendingProducts handles GET /api/v1/products/trending
// @Summary Get trending products
// @Description Retrieve currently trending products, optionally filtered by category
// @Tags products
// @Produce json
// @Param category_id query int false "Category ID to filter trending products"
// @Param limit query int false "Number of trending products to return" default(10)
// @Success 200 {object} []dto.ProductRecommendation
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/products/trending [get]
func (h *RecommendationHandler) GetTrendingProducts(c *gin.Context) {
	// Parse optional category ID
	var categoryID *int
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		parsedCategoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "invalid category_id format",
			})
			return
		}
		categoryID = &parsedCategoryID
	}

	// Parse limit
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 || parsedLimit > 100 {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Bad Request",
				Message: "limit must be a positive integer between 1 and 100",
			})
			return
		}
		limit = parsedLimit
	}

	products, err := h.recommendationService.GetTrendingProducts(c.Request.Context(), categoryID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "failed to get trending products: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

// LogRecommendationInteraction handles POST /api/v1/recommendations/interactions
// @Summary Log customer interaction with recommendations
// @Description Log customer clicks, purchases, and other interactions with recommended products for analytics
// @Tags recommendations
// @Accept json
// @Produce json
// @Param interaction body dto.RecommendationAnalytics true "Recommendation interaction data"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/recommendations/interactions [post]
func (h *RecommendationHandler) LogRecommendationInteraction(c *gin.Context) {
	var analytics dto.RecommendationAnalytics
	if err := c.ShouldBindJSON(&analytics); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if analytics.CustomerID == uuid.Nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "customer_id is required",
		})
		return
	}

	if analytics.RecommendationID == uuid.Nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "recommendation_id is required",
		})
		return
	}

	err := h.recommendationService.LogRecommendationInteraction(c.Request.Context(), &analytics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "failed to log interaction: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "interaction logged successfully",
	})
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message"`
}
