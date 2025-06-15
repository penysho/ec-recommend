package router

import (
	"ec-recommend/internal/handler"
	"ec-recommend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the HTTP router
func SetupRouter(
	chatHandler *handler.ChatHandler,
	healthHandler *handler.HealthHandler,
	recommendationHandler *handler.RecommendationHandler,
	recommendationHandlerV2 *handler.RecommendationHandlerV2,
) *gin.Engine {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Add middlewares
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", healthHandler.Check)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Chat endpoints
		chat := v1.Group("/chat")
		{
			chat.POST("/ask", chatHandler.AskQuestion)
			chat.POST("/messages", chatHandler.Chat)
		}

		// Recommendation endpoints (V1)
		recommendations := v1.Group("/recommendations")
		{
			recommendations.GET("", recommendationHandler.GetRecommendations)
			recommendations.POST("", recommendationHandler.PostRecommendations)
			recommendations.POST("/interactions", recommendationHandler.LogRecommendationInteraction)
		}

		// Customer endpoints
		customers := v1.Group("/customers")
		{
			customers.GET("/:customer_id/profile", recommendationHandler.GetCustomerProfile)
		}

		// Product endpoints
		products := v1.Group("/products")
		{
			products.GET("/trending", recommendationHandler.GetTrendingProducts)
			products.GET("/similar/:product_id", recommendationHandler.GetSimilarProducts)
		}
	}

	// API V2 routes - Enhanced RAG-based recommendations
	v2 := router.Group("/api/v2")
	{
		// Enhanced recommendation endpoints with RAG capabilities
		recommendations := v2.Group("/recommendations")
		{
			recommendations.GET("", recommendationHandlerV2.GetRecommendationsV2)
			recommendations.POST("", recommendationHandlerV2.PostRecommendationsV2)
			recommendations.GET("/semantic-search", recommendationHandlerV2.GetSemanticSearch)
			recommendations.GET("/vector-similar/:product_id", recommendationHandlerV2.GetVectorSimilarProducts)
			recommendations.GET("/knowledge-based", recommendationHandlerV2.GetKnowledgeBasedRecommendations)
			recommendations.GET("/:recommendation_id/explanation", recommendationHandlerV2.GetRecommendationExplanation)
		}

		// Enhanced product endpoints
		products := v2.Group("/products")
		{
			products.GET("/trending", recommendationHandlerV2.GetTrendingProductsV2)
		}
	}

	return router
}
