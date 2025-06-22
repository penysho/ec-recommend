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
	inventoryHandler *handler.InventoryHandler,
) *gin.Engine {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Add middlewares
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// Health check endpoint (no authentication required)
	router.GET("/health", healthHandler.Check)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Chat endpoints (require authentication)
		chat := v1.Group("/chat")
		chat.Use(middleware.RequireAuth())
		{
			chat.POST("/ask", chatHandler.AskQuestion)
			chat.POST("/messages", chatHandler.Chat)
		}

		// Recommendation endpoints (V1) - require authentication
		recommendations := v1.Group("/recommendations")
		recommendations.Use(middleware.RequireAuth())
		{
			// Customer access (all authenticated users can access)
			recommendations.GET("", recommendationHandler.GetRecommendations)
			recommendations.POST("", recommendationHandler.PostRecommendations)
			recommendations.POST("/interactions", recommendationHandler.LogRecommendationInteraction)
		}

		// Customer endpoints - require authentication and proper access control
		customers := v1.Group("/customers")
		customers.Use(middleware.RequireAuth())
		{
			// Customer can only access their own profile, admin/employee can access any
			customers.GET("/:customer_id/profile",
				middleware.RequireRole(middleware.RoleAdmin, middleware.RoleEmployee, middleware.RoleCustomer),
				recommendationHandler.GetCustomerProfile)
		}

		// Product endpoints - require authentication
		products := v1.Group("/products")
		products.Use(middleware.RequireAuth())
		{
			// Read-only operations - all authenticated users
			products.GET("/trending", recommendationHandler.GetTrendingProducts)
			products.GET("/similar/:product_id", recommendationHandler.GetSimilarProducts)
		}

		// Inventory endpoints - NEW
		inventory := v1.Group("/inventory")
		inventory.Use(middleware.RequireAuth())
		{
			// Read operations - all authenticated users can view inventory
			inventory.GET("", inventoryHandler.GetInventoryList)
			inventory.GET("/products/:product_id", inventoryHandler.GetInventoryByProductID)
			inventory.GET("/products/:product_id/history", inventoryHandler.GetInventoryHistory)

			// Management operations - admin and employee only
			inventory.PUT("/products/:product_id",
				middleware.RequireRole(middleware.RoleAdmin, middleware.RoleEmployee),
				inventoryHandler.UpdateInventory)
			inventory.POST("/transactions",
				middleware.RequireRole(middleware.RoleAdmin, middleware.RoleEmployee),
				inventoryHandler.PerformInventoryTransaction)
			inventory.GET("/alerts",
				middleware.RequireRole(middleware.RoleAdmin, middleware.RoleEmployee),
				inventoryHandler.GetInventoryAlerts)
			inventory.GET("/stats",
				middleware.RequireRole(middleware.RoleAdmin, middleware.RoleEmployee),
				inventoryHandler.GetInventoryStats)

			// Batch operations - admin only
			inventory.POST("/batch-update",
				middleware.RequireRole(middleware.RoleAdmin),
				inventoryHandler.BatchUpdateInventory)
		}
	}

	// API V2 routes - Enhanced RAG-based recommendations
	v2 := router.Group("/api/v2")
	v2.Use(middleware.RequireAuth()) // All V2 endpoints require authentication
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
