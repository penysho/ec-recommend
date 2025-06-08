package router

import (
	"ec-recommend/internal/handler"
	"ec-recommend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the HTTP router
func SetupRouter(chatHandler *handler.ChatHandler, healthHandler *handler.HealthHandler, recommendationHandler *handler.RecommendationHandler) *gin.Engine {
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

		// Recommendation endpoints
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
			products.GET("/:product_id/similar", recommendationHandler.GetSimilarProducts)
		}
	}

	return router
}
